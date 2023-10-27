package user

import (
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

type UserService interface {
	Get(id domain.UserId) (*User, error)
	GetByEmail(email domain.Email) (*User, error)
}

type userServiceImpl struct {
	userStorage UserStorage
}

func NewUserService(userStorage UserStorage) UserService {
	return &userServiceImpl{userStorage: userStorage}
}

func (service *userServiceImpl) Get(id domain.UserId) (*User, error) {
	return service.userStorage.Get(id)
}

func (service *userServiceImpl) GetByEmail(email domain.Email) (*User, error) {
	return service.userStorage.GetByEmail(email)
}
