package auth

import (
	"github.com/vaberof/ssugt-projects/internal/domain/user"
	"github.com/vaberof/ssugt-projects/pkg/domain"
)

type UserService interface {
	Get(id domain.UserId) (*user.User, error)
	GetByEmail(email domain.Email) (*user.User, error)
}
