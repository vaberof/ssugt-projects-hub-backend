package auth

import (
	"context"
	"fmt"
	"log/slog"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth/accesstoken"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xpassword"
	"ssugt-projects-hub/service/user"
	"time"
)

type AccessToken string

const (
	_defaultTokenTtl = 24 * time.Hour
)

type Service interface {
	Register(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, loginUserRequestParams models.LoginUserRequestParams) (*AccessToken, error)
}

type authServiceImpl struct {
	log         *slog.Logger
	userService user.Service
}

func NewService(log *logs.Logs, userService user.Service) Service {
	return &authServiceImpl{
		log:         log.WithName("auth-service"),
		userService: userService,
	}
}

func (as authServiceImpl) Register(ctx context.Context, user models.User) (models.User, error) {
	hashedPassword, err := xpassword.Hash(user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	user, err = as.userService.Create(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create a user: %w", err)
	}

	return user, nil
}

func (as authServiceImpl) Login(ctx context.Context, loginUserRequestParams models.LoginUserRequestParams) (*AccessToken, error) {
	user, err := as.userService.GetByEmail(ctx, loginUserRequestParams.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if err = xpassword.Check(loginUserRequestParams.Password, user.Password); err != nil {
		return nil, fmt.Errorf("incorrect password: %w", err)
	}

	return as.getAccessToken(user.Id, _defaultTokenTtl)
}

func (as authServiceImpl) getAccessToken(userId int, ttl time.Duration) (*AccessToken, error) {
	accessTokenString, err := accesstoken.Create(userId, ttl)
	if err != nil {
		return nil, err
	}

	accessToken := AccessToken(accessTokenString)

	return &accessToken, nil
}
