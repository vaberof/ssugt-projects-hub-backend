package user

import "ssugt-projects-hub/models"

func mapToUser(dbUser DbUser) models.User {
	roles := make([]models.UserRole, 0, len(dbUser.Roles))
	for _, role := range dbUser.Roles {
		roles = append(roles, models.UserRole(role.Id))
	}

	return models.User{
		Id:                     dbUser.Id,
		Email:                  dbUser.Email,
		Password:               dbUser.PasswordHash,
		FullName:               dbUser.FullName,
		PhoneNumber:            dbUser.PhoneNumber,
		IsEmailConfirmed:       dbUser.IsEmailConfirmed,
		IsPhoneNumberConfirmed: dbUser.IsPhoneNumberConfirmed,
		Roles:                  roles,
	}
}

func mapToDbUser(user models.User) DbUser {
	roles := make([]DbRole, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, DbRole{Id: int(role), UserId: user.Id})
	}

	return DbUser{
		Id:                     user.Id,
		Email:                  user.Email,
		PasswordHash:           user.Password,
		FullName:               user.FullName,
		PhoneNumber:            user.PhoneNumber,
		IsEmailConfirmed:       user.IsEmailConfirmed,
		IsPhoneNumberConfirmed: user.IsPhoneNumberConfirmed,
		Roles:                  roles,
		Profile:                mapToDbUserProfile(user),
	}
}

func mapToDbUserProfile(user models.User) DbUserProfile {
	return DbUserProfile{
		UserId:       user.Id,
		PersonalInfo: mapToDbPersonalInfo(user.PersonalInfo),
		Settings:     mapToDbSettings(user.Settings),
	}
}

func mapToDbPersonalInfo(personalInfo models.PersonalInfo) DbPersonalInfo {
	return DbPersonalInfo{
		HasOrganisation: personalInfo.HasOrganisation,
		Organisation:    mapToDbOrganisation(personalInfo.Organisation),
		HasEducation:    personalInfo.HasEducation,
		Education:       mapToDbEducations(personalInfo.Education),
	}
}

func mapToDbOrganisation(organisation models.Organisation) DbOrganisation {
	return DbOrganisation{
		Name:    organisation.Name,
		Address: organisation.Address,
	}
}

func mapToDbEducations(education []models.Education) []DbEducation {
	dbEducations := make([]DbEducation, 0, len(education))
	for _, e := range education {
		dbEducations = append(dbEducations, mapToDbEducation(e))
	}
	return dbEducations
}

func mapToDbEducation(education models.Education) DbEducation {
	return DbEducation{
		Degree: education.Degree,
		Course: education.Course,
		Group:  education.Group,
	}
}

func mapToDbSettings(settings models.Settings) DbSettings {
	return DbSettings{}
}
