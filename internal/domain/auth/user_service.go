package auth

import (
	domain2 "github.com/vaberof/ssugt-projects-hub-backend/internal/domain/user"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

type UserService interface {
	Get(id domain.UserId) (*domain2.User, error)
	GetByEmail(email domain.Email) (*domain2.User, error)
}
