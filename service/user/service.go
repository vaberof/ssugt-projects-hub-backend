package user

import (
	"context"
	"log/slog"
	"ssugt-projects-hub/database/postgres/user"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/logging/logs"
)

type Service interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByIds(ctx context.Context, userIds []int) ([]models.User, error)
}

type userServiceImpl struct {
	log            *slog.Logger
	userRepository user.Repository
}

func NewService(log *logs.Logs, userRepository user.Repository) Service {
	return &userServiceImpl{
		log:            log.WithName("user-service"),
		userRepository: userRepository,
	}
}

func (us userServiceImpl) Create(ctx context.Context, user models.User) (models.User, error) {
	u, err := us.userRepository.Insert(ctx, user)
	if err != nil {
		us.log.Error("failed to create new user", "error", err)
		return models.User{}, err
	}

	return u, nil
}

func (us userServiceImpl) GetByEmail(ctx context.Context, email string) (models.User, error) {
	u, err := us.userRepository.GetByEmail(ctx, email)
	if err != nil {
		us.log.Error("failed to get user by email", "error", err)
		return models.User{}, err
	}
	return u, nil
}

func (us userServiceImpl) GetByIds(ctx context.Context, userIds []int) ([]models.User, error) {
	u, err := us.userRepository.GetByIds(ctx, userIds)
	if err != nil {
		us.log.Error("failed to get user by id", "error", err)
		return []models.User{}, err
	}
	return u, nil
}
