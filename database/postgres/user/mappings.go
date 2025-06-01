package user

import "ssugt-projects-hub/models"

func mapToUsers(dbUser []DbUser) []models.User {
	users := make([]models.User, len(dbUser))
	for i := range dbUser {
		users[i] = mapToUser(dbUser[i])
	}
	return users
}

func mapToUser(dbUser DbUser) models.User {
	return models.User{
		Id:       dbUser.Id,
		Email:    dbUser.Email,
		Password: dbUser.PasswordHash,
		FullName: dbUser.FullName,
		PersonalInfo: models.PersonalInfo{
			HasOrganisation: dbUser.Profile.PersonalInfo.HasOrganisation,
			Organisation: models.Organisation{
				Name:    dbUser.Profile.PersonalInfo.Organisation.Name,
				Address: dbUser.Profile.PersonalInfo.Organisation.Address,
			},
		},
		Role:      models.UserRole(dbUser.RoleId),
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func mapToDbUser(user models.User) DbUser {
	return DbUser{
		Id:           user.Id,
		RoleId:       int(user.Role),
		Email:        user.Email,
		PasswordHash: user.Password,
		FullName:     user.FullName,
		Profile:      mapToDbUserProfile(user),
		CreatedAt:    user.CreatedAt,
	}
}

func mapToDbUserProfile(user models.User) DbUserProfile {
	return DbUserProfile{
		UserId:       user.Id,
		PersonalInfo: mapToDbPersonalInfo(user.PersonalInfo),
		CreatedAt:    user.CreatedAt,
	}
}

func mapToDbPersonalInfo(personalInfo models.PersonalInfo) DbPersonalInfo {
	return DbPersonalInfo{
		HasOrganisation: personalInfo.HasOrganisation,
		Organisation:    mapToDbOrganisation(personalInfo.Organisation),
	}
}

func mapToDbOrganisation(organisation models.Organisation) DbOrganisation {
	return DbOrganisation{
		Name:    organisation.Name,
		Address: organisation.Address,
	}
}
