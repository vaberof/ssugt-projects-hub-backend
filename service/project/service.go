package project

import (
	"context"
	"log/slog"
	"ssugt-projects-hub/database/mongo/files"
	projectrepo "ssugt-projects-hub/database/postgres/project"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/logging/logs"
	"time"
)

type Service interface {
	Create(ctx context.Context, project models.Project) (models.Project, error)
	GetById(ctx context.Context, id int) (models.Project, error)
	GetByUserId(ctx context.Context, userId int) ([]models.Project, error)
	Update(ctx context.Context, project models.Project) (models.Project, error)
	Search(ctx context.Context, filters models.ProjectSearchFilters) ([]models.Project, error)
}

type projectServiceImpl struct {
	log               *slog.Logger
	projectRepository projectrepo.Repository
	fileRepository    files.Repository
}

func NewProjectService(logs *logs.Logs, projectRepository projectrepo.Repository, fileRepository files.Repository) Service {
	return &projectServiceImpl{
		log:               logs.WithName("project-service"),
		projectRepository: projectRepository,
		fileRepository:    fileRepository,
	}
}

func (p projectServiceImpl) Create(ctx context.Context, project models.Project) (models.Project, error) {
	project.Status = models.InProcessProjectStatus
	project.CreatedAt = time.Now().UTC()

	project, err := p.projectRepository.Insert(ctx, project)
	if err != nil {
		p.log.Error("Failed to create project", "error", err)
		return models.Project{}, err
	}

	files, err := p.fileRepository.GetByProjectId(ctx, project.Id)
	if err != nil {
		p.log.Error("Failed to get project files", "error", err)
		return models.Project{}, err
	}

	project.Files = files

	return project, nil
}

func (p projectServiceImpl) GetById(ctx context.Context, id int) (models.Project, error) {
	project, err := p.projectRepository.GetById(ctx, id)
	if err != nil {
		p.log.Error("Failed to get project", "error", err)
		return models.Project{}, err
	}

	files, err := p.fileRepository.GetByProjectId(ctx, id)
	if err != nil {
		p.log.Error("Failed to get project files", "error", err)
		return models.Project{}, err
	}

	project.Files = files

	return project, nil
}

func (p projectServiceImpl) GetByUserId(ctx context.Context, userId int) ([]models.Project, error) {
	projects, err := p.projectRepository.GetByUserId(ctx, userId)
	if err != nil {
		p.log.Error("Failed to get projects", "error", err)
		return []models.Project{}, err
	}

	for i := range projects {
		files, err := p.fileRepository.GetByProjectId(ctx, projects[i].Id)
		if err != nil {
			p.log.Error("Failed to get project files", "error", err)
			return []models.Project{}, err
		}
		projects[i].Files = files
	}

	return projects, nil
}

func (p projectServiceImpl) Update(ctx context.Context, project models.Project) (models.Project, error) {
	project, err := p.projectRepository.Update(ctx, project)
	if err != nil {
		p.log.Error("Failed to update", "error", err)
		return models.Project{}, err
	}

	files, err := p.fileRepository.GetByProjectId(ctx, project.Id)
	if err != nil {
		p.log.Error("Failed to get project files", "error", err)
		return models.Project{}, err
	}

	project.Files = files

	return project, nil
}

func (p projectServiceImpl) Search(ctx context.Context, filters models.ProjectSearchFilters) ([]models.Project, error) {
	projects, err := p.projectRepository.Search(ctx, filters)
	if err != nil {
		p.log.Error("Failed to search projects", "error", err)
		return []models.Project{}, err
	}

	for i := range projects {
		files, err := p.fileRepository.GetByProjectId(ctx, projects[i].Id)
		if err != nil {
			p.log.Error("Failed to get project files", "error", err)
			return []models.Project{}, err
		}
		projects[i].Files = files
	}

	return projects, nil
}
