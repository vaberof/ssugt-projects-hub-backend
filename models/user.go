package models

import "time"

type UserRole int

const (
	DefaultRole UserRole = iota + 1
	RoleAdmin
)

type User struct {
	Id           int          `json:"id"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	FullName     string       `json:"fullName"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Role         UserRole     `json:"role"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type PersonalInfo struct {
	HasOrganisation bool         `json:"hasOrganisation"`
	Organisation    Organisation `json:"organisation"`
}

type Organisation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
