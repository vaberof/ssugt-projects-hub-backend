package models

import "time"

type GetUsersRequest struct {
	Ids []int `json:"ids"`
}

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type GetUserResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	Id           int          `json:"id"`
	Email        string       `json:"email"`
	FullName     string       `json:"fullName"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Role         UserRole     `json:"role"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

func MapToUsersResponse(users []User) []UserResponse {
	usersResponse := make([]UserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = MapToUserResponse(user)
	}
	return usersResponse
}

func MapToUserResponse(user User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		FullName: user.FullName,
		PersonalInfo: PersonalInfo{
			HasOrganisation: user.PersonalInfo.HasOrganisation,
			Organisation: Organisation{
				Name:    user.PersonalInfo.Organisation.Name,
				Address: user.PersonalInfo.Organisation.Address,
			},
		},
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
