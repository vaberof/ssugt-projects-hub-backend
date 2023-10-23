package user

import (
	"github.com/vaberof/ssugt-projects/pkg/domain"
)

type UserStorage interface {
	Get(id domain.UserId) (*User, error)
	GetByEmail(email domain.Email) (*User, error)
}
