package auth

import (
	"ssugt-projects-hub/database/mongo/cache"
	"ssugt-projects-hub/models"
)

func mapCachedUserDataToUser(userData cache.UserData) models.User {
	return models.User{
		Email:    userData.Email,
		Password: userData.Password,
		FullName: userData.FullName,
		PersonalInfo: models.PersonalInfo{
			HasOrganisation: userData.PersonalInfo.HasOrganisation,
			Organisation: models.Organisation{
				Name:    userData.PersonalInfo.Organisation.Name,
				Address: userData.PersonalInfo.Organisation.Address,
			},
		},
		Role: models.UserRole(userData.Role),
	}
}
