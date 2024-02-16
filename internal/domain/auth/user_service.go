package auth

import (
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/user"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

type UserService interface {
	Get(id domain.UserId) (*user.User, error)
	GetByEmail(email domain.Email) (*user.User, error)
}
