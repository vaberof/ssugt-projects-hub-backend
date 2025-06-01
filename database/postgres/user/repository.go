package user

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"ssugt-projects-hub/models"
	"time"
)

type Repository interface {
	Insert(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByIds(ctx context.Context, userIds []int) ([]models.User, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

//go:embed sql/insert_user.sql
var _insertUserSql string

func (r repositoryImpl) Insert(ctx context.Context, user models.User) (models.User, error) {
	dbUser := mapToDbUser(user)

	createdAt := time.Now().UTC()

	dbUser.CreatedAt = createdAt
	dbUser.Profile.CreatedAt = createdAt

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, tx.Rebind(_insertUserSql))
	if err != nil {
		return models.User{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &dbUser.Id, &dbUser)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to execute statement HERE: %w", err)
	}

	dbUser.Profile, err = insertUserProfile(ctx, tx, dbUser.Id, dbUser.Profile)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user profile: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapToUser(dbUser), nil
}

//go:embed sql/get_user_by_email.sql
var _getUserByEmailSql string

func (r repositoryImpl) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var dbUser DbUser

	err := r.db.GetContext(ctx, &dbUser, _getUserByEmailSql, email)
	if err != nil {
		return models.User{}, fmt.Errorf("GetByEmail: failed to execute statement HERE1: %w", err)
	}

	dbUser.Profile, err = r.getUserProfileByUserId(ctx, dbUser.Id)
	if err != nil {
		return models.User{}, fmt.Errorf("GetByEmail: failed to execute statement HERE2: %w", err)
	}

	return mapToUser(dbUser), nil
}

//go:embed sql/get_user_by_ids.sql
var _getUserByIdsSql string

func (r repositoryImpl) GetByIds(ctx context.Context, ids []int) ([]models.User, error) {
	var dbUsers []DbUser

	err := r.db.SelectContext(ctx, &dbUsers, _getUserByIdsSql, ids)
	if err != nil {
		return []models.User{}, fmt.Errorf("GetByIds: failed to execute statement HERE1: %w", err)
	}

	for i := range dbUsers {
		dbUsers[i].Profile, err = r.getUserProfileByUserId(ctx, dbUsers[i].Id)
		if err != nil {
			return []models.User{}, fmt.Errorf("GetByIds: failed to execute statement HERE2: %w", err)
		}
	}

	return mapToUsers(dbUsers), nil
}

//go:embed sql/get_user_profile_by_user_id.sql
var _getUserProfileByUserIdSql string

func (r repositoryImpl) getUserProfileByUserId(ctx context.Context, userId int) (DbUserProfile, error) {
	var profile DbUserProfile

	err := r.db.GetContext(ctx, &profile, _getUserProfileByUserIdSql, userId)
	if err != nil {
		return DbUserProfile{}, fmt.Errorf("failed to execute statement: %w", err)
	}

	return profile, nil
}

//go:embed sql/insert_user_profile.sql
var _insertUserProfileSql string

func insertUserProfile(ctx context.Context, tx *sqlx.Tx, userId int, profile DbUserProfile) (DbUserProfile, error) {
	var (
		stmt *sqlx.NamedStmt
		err  error
	)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	stmt, err = tx.PrepareNamedContext(ctx, tx.Rebind(_insertUserProfileSql))
	if err != nil {
		return DbUserProfile{}, fmt.Errorf("failed to prepare statement: %w", err)
	}

	profile.UserId = userId

	err = stmt.GetContext(ctx, &profile.Id, profile)
	if err != nil {
		return DbUserProfile{}, fmt.Errorf("failed to execute statement USERPROFILE: %w", err)
	}

	return profile, nil
}
