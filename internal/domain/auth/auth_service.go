package auth

import (
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/auth/accesstoken"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"time"
)

const defaultTokenTtl = 1 * time.Hour

type AuthService interface {
	Login(email domain.Email, password domain.Password) (*AccessToken, error)
	VerifyAccessToken(token string) (*domain.UserId, error)
}

type authServiceImpl struct {
	config      AuthConfig
	userService UserService
}

func NewAuthService(config AuthConfig, userService UserService) AuthService {
	return &authServiceImpl{
		config:      config,
		userService: userService,
	}
}

func (service *authServiceImpl) Login(email domain.Email, password domain.Password) (*AccessToken, error) {
	user, err := service.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("incorrect password")
	}
	/*err = xpassword.Check(password.String(), user.Password.String())
	if err != nil {
		// TODO: hash password
		//return nil, nil, fmt.Errorf("incorrect password: %w", err)
	}*/

	return service.getAccessToken(domain.UserId(user.Id.String()), defaultTokenTtl, service.config.AccessTokenSecretKey)
}

func (service *authServiceImpl) VerifyAccessToken(token string) (*domain.UserId, error) {
	payload, err := accesstoken.Verify(token, accesstoken.SecretKey(service.config.AccessTokenSecretKey))
	if err != nil {
		return nil, err
	}

	fmt.Printf("userId from payload: %v\n", payload.UserId)

	user, err := service.userService.Get(payload.UserId)
	if err != nil {
		return nil, err
	}

	domainUserId := domain.UserId(user.Id.String())

	return &domainUserId, nil
}

func (service *authServiceImpl) getAccessToken(userId domain.UserId, ttl time.Duration, jwtSecretKey string) (*AccessToken, error) {
	accessTokenString, err := accesstoken.Create(userId, ttl, accesstoken.SecretKey(jwtSecretKey))
	if err != nil {
		return nil, err
	}

	accessToken := AccessToken(accessTokenString)

	return &accessToken, nil
}
