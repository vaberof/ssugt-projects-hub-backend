package models

type RegisterUserRequest struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	FullName     string       `json:"fullName"`
	PhoneNumber  string       `json:"phoneNumber"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Roles        []UserRole   `json:"roles"`
}

type RegisterUserResponse struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	FullName     string       `json:"fullName"`
	PhoneNumber  string       `json:"phoneNumber"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Roles        []UserRole   `json:"roles"`
}
