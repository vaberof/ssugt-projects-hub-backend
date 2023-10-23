package user

import (
	"github.com/vaberof/ssugt-projects/pkg/domain"
	"time"
)

const (
	RoleUser       = domain.Role("user")
	RoleAdmin      = domain.Role("admin")
	RoleSuperAdmin = domain.Role("superAdmin")
)

type User struct {
	Id           domain.UserId
	Role         domain.Role
	FullName     domain.FullName
	Email        domain.Email
	Password     domain.Password
	RegisteredAt time.Time
}
