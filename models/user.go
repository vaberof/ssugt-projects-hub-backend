package models

type UserRole int

const (
	DefaultRole UserRole = iota + 1
	ModeratorRole
	RoleAdmin
)

type User struct {
	Id                     int          `json:"id"`
	Email                  string       `json:"email"`
	Password               string       `json:"password"`
	FullName               string       `json:"fullName"`
	PhoneNumber            string       `json:"phoneNumber"`
	IsEmailConfirmed       bool         `json:"isEmailConfirmed"`
	IsPhoneNumberConfirmed bool         `json:"isPhoneNumberConfirmed"`
	PersonalInfo           PersonalInfo `json:"personalInfo"`
	Settings               Settings     `json:"settings"`
	Roles                  []UserRole   `json:"roles"`
}

type PersonalInfo struct {
	HasOrganisation bool         `json:"hasOrganisation"`
	Organisation    Organisation `json:"organisation"`
	HasEducation    bool         `json:"hasEducation"`
	Education       []Education  `json:"education"`
}

type Organisation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Education struct {
	Degree string `json:"degree"`
	Course string `json:"course"`
	Group  string `json:"group"`
}

type Settings struct {
}
