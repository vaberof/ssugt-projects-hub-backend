package project

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"ssugt-projects-hub/models"
	"strings"
)

type Repository interface {
	Insert(ctx context.Context, project models.Project) (models.Project, error)
	GetById(ctx context.Context, id int) (models.Project, error)
	GetByUserId(ctx context.Context, userId int) ([]models.Project, error)
	Update(ctx context.Context, project models.Project) (models.Project, error)
	Search(ctx context.Context, filters models.ProjectSearchFilters) ([]models.Project, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

//go:embed sql/insert_project.sql
var _insertProjectSql string

func (r repositoryImpl) Insert(ctx context.Context, project models.Project) (models.Project, error) {
	dbProject := mapToDbProject(project)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, tx.Rebind(_insertProjectSql))
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	if err = stmt.GetContext(ctx, &dbProject.Id, &dbProject); err != nil {
		return models.Project{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err = insertCollaborators(ctx, tx, dbProject.Id, dbProject.Collaborators); err != nil {
		return models.Project{}, fmt.Errorf("failed to insert collaborators: %w", err)
	}

	if err = insertProjectReview(ctx, tx, createInitialProjectReview(dbProject)); err != nil {
		return models.Project{}, fmt.Errorf("failed to insert project review: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.Project{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapToProject(dbProject), nil
}

//go:embed sql/get_project_by_id.sql
var _getProjectByIdSql string

func (r repositoryImpl) GetById(ctx context.Context, id int) (models.Project, error) {
	var dbProject DbProject

	err := r.db.GetContext(ctx, &dbProject, _getProjectByIdSql, id)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	dbProject.Collaborators, err = r.getCollaborators(ctx, dbProject.Id)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to get project collaborators: %w", err)
	}

	return mapToProject(dbProject), nil
}

//go:embed sql/get_projects_by_user_id.sql
var _getProjectsByUserIdSql string

func (r repositoryImpl) GetByUserId(ctx context.Context, userId int) ([]models.Project, error) {
	var dbProjects []DbProject

	err := r.db.GetContext(ctx, &dbProjects, _getProjectsByUserIdSql, userId)
	if err != nil {
		return []models.Project{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	for i := range dbProjects {
		dbProjects[i].Collaborators, err = r.getCollaborators(ctx, dbProjects[i].Id)
		if err != nil {
			return []models.Project{}, fmt.Errorf("failed to get project collaborators: %w", err)
		}
	}

	return mapToProjects(dbProjects), nil
}

//go:embed sql/update_project.sql
var _updateProjectSql string

func (r repositoryImpl) Update(ctx context.Context, project models.Project) (models.Project, error) {
	dbProject := mapToDbProject(project)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, tx.Rebind(_updateProjectSql))
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, dbProject)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	err = updateCollaborators(ctx, tx, dbProject.Id, dbProject.Collaborators)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to update project collaborators: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.Project{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapToProject(dbProject), nil
}

//go:embed sql/update_project_collaborator.sql
var _updateCollaboratorSql string

//go:embed sql/delete_removed_collaborators.sql
var _deleteRemovedCollaborators string

func updateCollaborators(ctx context.Context, tx *sqlx.Tx, projectId int, collaborators []DbCollaborator) error {
	collaboratorsPK := lo.Map(collaborators, func(collaborator DbCollaborator, _ int) int { return collaborator.Id })
	if _, err := tx.ExecContext(ctx, _deleteRemovedCollaborators, projectId, pq.Array(collaboratorsPK)); err != nil {
		return err
	}

	var (
		stmt *sqlx.NamedStmt
		err  error
	)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	for i := range collaborators {
		collaborators[i].ProjectId = projectId

		if collaborators[i].Id == 0 {
			stmt, err = tx.PrepareNamedContext(ctx, tx.Rebind(_insertCollaboratorSql))
			if err != nil {
				return fmt.Errorf("failed to prepare statement: %w", err)
			}

			err = stmt.GetContext(ctx, &collaborators[i].Id, collaborators[i])
			if err != nil {
				return fmt.Errorf("failed to execute statement: %w", err)
			}

			continue
		}

		if _, err = tx.NamedExecContext(ctx, _updateCollaboratorSql, collaborators[i]); err != nil {
			return fmt.Errorf("failed to update project collaborators: %w", err)
		}
	}

	return nil
}

//go:embed sql/get_projects_by_filters.sql
var _getProjectsByFiltersSql string

func (r repositoryImpl) Search(ctx context.Context, filters models.ProjectSearchFilters) ([]models.Project, error) {
	query := _getProjectsByFiltersSql

	var conditions []string
	var args []interface{}

	// Добавляем условия на основе фильтров
	if filters.BaseFilters.Type != 0 {
		conditions = append(conditions, "type_id = $1")
		args = append(args, filters.BaseFilters.Type)
	}

	if filters.BaseFilters.Status != "" {
		conditions = append(conditions, "status = $2")
		args = append(args, filters.BaseFilters.Status)
	}

	if filters.AttributeFilters.Title != "" {
		conditions = append(conditions, "LOWER(attributes->>'title') LIKE $3")
		args = append(args, "%"+strings.ToLower(filters.AttributeFilters.Title)+"%")
	}

	if len(filters.AttributeFilters.Tags) > 0 {
		conditions = append(conditions, "attributes->'tags' ?| $4")
		args = append(args, pq.Array(filters.AttributeFilters.Tags))
	}

	// Добавляем условия к базовому запросу
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Выполнение запроса
	var dbProjects []DbProject
	err := r.db.SelectContext(ctx, &dbProjects, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	for i := range dbProjects {
		dbProjects[i].Collaborators, err = r.getCollaborators(ctx, dbProjects[i].Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get project collaborators: %w", err)
		}
	}

	return mapToProjects(dbProjects), nil
}

//go:embed sql/get_collaborators.sql
var _getCollaboratorsSql string

func (r repositoryImpl) getCollaborators(ctx context.Context, projectId int) ([]DbCollaborator, error) {
	var dbCollaborators []DbCollaborator
	return dbCollaborators, r.db.GetContext(ctx, &dbCollaborators, _getCollaboratorsSql, projectId)
}

//go:embed sql/get_collaborators_by_project_ids.sql
var _getCollaboratorsByProjectIdsSql string

func (r repositoryImpl) getProjectsCollaborators(ctx context.Context, projectIds []int) ([]DbCollaborator, error) {
	var dbCollaborators []DbCollaborator

	if err := r.db.GetContext(ctx, &dbCollaborators, _getCollaboratorsByProjectIdsSql, pq.Array(projectIds)); err != nil {
		return []DbCollaborator{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	return dbCollaborators, nil
}

//go:embed sql/insert_project_collaborator.sql
var _insertCollaboratorSql string

func insertCollaborators(ctx context.Context, tx *sqlx.Tx, projectId int, collaborators []DbCollaborator) error {
	var (
		stmt *sqlx.NamedStmt
		err  error
	)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	for i := range collaborators {
		stmt, err = tx.PrepareNamedContext(ctx, tx.Rebind(_insertCollaboratorSql))
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}

		collaborators[i].ProjectId = projectId

		err = stmt.GetContext(ctx, &collaborators[i].Id, collaborators[i])
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}

//go:embed sql/insert_project_review.sql
var _insertProjectReviewSql string

func insertProjectReview(ctx context.Context, tx *sqlx.Tx, projectReview DbProjectReview) error {
	var (
		stmt *sqlx.NamedStmt
		err  error
	)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	stmt, err = tx.PrepareNamedContext(ctx, tx.Rebind(_insertProjectReviewSql))
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	err = stmt.GetContext(ctx, &projectReview.Id, projectReview)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func getProjectsIds(dbProjects []DbProject) []int {
	return lo.Map(dbProjects, func(dbProject DbProject, index int) int {
		return dbProject.Id
	})
}

func createInitialProjectReview(project DbProject) DbProjectReview {
	return DbProjectReview{
		Id:         0,
		ProjectId:  project.Id,
		ReviewedBy: 0,
		Status:     project.Status,
		Comment:    "",
		CreatedAt:  project.CreatedAt,
	}
}
