package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ssugt-projects-hub/models"
	"time"
)

func MapToEmailConfirmation(user models.User, code string) EmailConfirmation {
	return EmailConfirmation{
		Id:    primitive.NewObjectID(),
		Email: user.Email,
		UserData: UserData{
			Email:    user.Email,
			Password: user.Password,
			FullName: user.FullName,
			PersonalInfo: PersonalInfo{
				HasOrganisation: user.PersonalInfo.HasOrganisation,
				Organisation: Organisation{
					Name:    user.PersonalInfo.Organisation.Name,
					Address: user.PersonalInfo.Organisation.Address,
				},
			},
			Role: int(user.Role),
		},
		Code:      code,
		CreatedAt: time.Now().UTC(),
	}
}
