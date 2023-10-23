package auth

import (
	"github.com/vaberof/ssugt-projects/pkg/auth/accesstoken"
	"github.com/vaberof/ssugt-projects/pkg/domain"
	"time"
)

const defaultTokenTtl = 1 * time.Hour

type AuthService interface {
	Login(email domain.Email, password domain.Password) (*AccessToken, error)
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

	/*err = xpassword.Check(password.String(), user.Password.String())
	if err != nil {
		// TODO: hash password
		//return nil, nil, fmt.Errorf("incorrect password: %w", err)
	}*/

	return service.getAccessToken(user.Id, defaultTokenTtl, service.config.AccessTokenSecretKey)
}

func (service *authServiceImpl) VerifyAccessToken(token string) (*domain.UserId, error) {
	payload, err := accesstoken.Verify(token, accesstoken.SecretKey(service.config.AccessTokenSecretKey))
	if err != nil {
		return nil, err
	}

	user, err := service.userService.Get(payload.UserId)
	if err != nil {
		return nil, err
	}

	return &user.Id, nil
}

func (service *authServiceImpl) getAccessToken(userId domain.UserId, ttl time.Duration, jwtSecretKey string) (*AccessToken, error) {
	accessTokenString, err := accesstoken.Create(userId, ttl, accesstoken.SecretKey(jwtSecretKey))
	if err != nil {
		return nil, err
	}

	accessToken := AccessToken(accessTokenString)

	return &accessToken, nil
}
