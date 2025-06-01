package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"ssugt-projects-hub/database/mongo/cache"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth/accesstoken"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xpassword"
	"ssugt-projects-hub/pkg/xrand"
	"ssugt-projects-hub/service/sender/email"
	"ssugt-projects-hub/service/user"
	"time"
)

type AccessToken string

const (
	_defaultTokenTtl = 24 * time.Hour
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrNeedEmailConfirmation = errors.New("need to confirm email")
)

type Service interface {
	Register(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, loginUserRequestParams models.LoginUserRequestParams) (*AccessToken, error)
	VerifyEmail(ctx context.Context, email, code string) error
	IsAdmin(ctx context.Context, userId int) (bool, error)
}

type authServiceImpl struct {
	log          *slog.Logger
	userService  user.Service
	emailService email.Service
	cache        cache.Cache
}

func NewService(log *logs.Logs, userService user.Service, emailService email.Service, cache cache.Cache) Service {
	return &authServiceImpl{
		log:          log.WithName("auth-service"),
		userService:  userService,
		emailService: emailService,
		cache:        cache,
	}
}

func (as authServiceImpl) Register(ctx context.Context, user models.User) (models.User, error) {
	_, err := as.userService.GetByEmail(ctx, user.Email)
	if err == nil {
		return models.User{}, ErrEmailAlreadyExists // errors.New("пользователь с таким email уже существует")
	}

	_, err = as.cache.Get(ctx, user.Email)
	if err == nil {
		return models.User{}, ErrNeedEmailConfirmation
	}

	hashedPassword, err := xpassword.Hash(user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.Role = models.DefaultRole

	code, err := xrand.GenerateRandomCode(6)
	if err != nil {
		as.log.Error("failed to generate random code: ", err.Error())
	}

	err = as.cache.Insert(ctx, cache.MapToEmailConfirmation(user, code))
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user into cache: %w", err)
	}

	go func() {
		err = as.emailService.SendConfirmationEmail(user.Email, code)
		if err != nil {
			as.log.Error("failed to send confirmation email: ", err.Error())
		}
	}()

	return user, nil
}

func (as authServiceImpl) VerifyEmail(ctx context.Context, email, code string) error {
	emailConfirmation, err := as.cache.Get(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to verify an email: %w", err)
	}

	if !isValidCode(code, emailConfirmation.Code) {
		return fmt.Errorf("invalid code")
	}

	_, err = as.userService.Create(ctx, mapCachedUserDataToUser(emailConfirmation.UserData))
	if err != nil {
		return fmt.Errorf("failed to create a user: %w", err)
	}

	go func() {
		err = as.cache.DeleteByEmail(ctx, email)
		if err != nil {
			log.Printf("failed to delete a user from cache: %v", err)
		}
	}()

	return nil
}

func (as authServiceImpl) Login(ctx context.Context, loginUserRequestParams models.LoginUserRequestParams) (*AccessToken, error) {
	user, err := as.userService.GetByEmail(ctx, loginUserRequestParams.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if err = xpassword.Check(loginUserRequestParams.Password, user.Password); err != nil {
		return nil, fmt.Errorf("incorrect password: %w", err)
	}

	return getAccessToken(user.Id, _defaultTokenTtl)
}

func (as authServiceImpl) IsAdmin(ctx context.Context, userId int) (bool, error) {
	users, err := as.userService.GetByIds(ctx, []int{userId})
	if err != nil {
		return false, fmt.Errorf("failed to get user by id: %w", err)
	}

	if len(users) == 0 {
		return false, fmt.Errorf("user not found by id")
	}

	isAdmin := users[0].Role == models.RoleAdmin

	return isAdmin, nil
}

func getAccessToken(userId int, ttl time.Duration) (*AccessToken, error) {
	accessTokenString, err := accesstoken.Create(userId, ttl)
	if err != nil {
		return nil, err
	}

	accessToken := AccessToken(accessTokenString)

	return &accessToken, nil
}

func isValidCode(code, initialCode string) bool {
	return code == initialCode
}
