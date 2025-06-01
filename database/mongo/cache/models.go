package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EmailConfirmation struct {
	Id        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"Email"`
	UserData  UserData           `bson:"UserData"`
	Code      string             `bson:"Code"`
	CreatedAt time.Time          `bson:"CreatedAt"`
}

type UserData struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	FullName     string       `json:"fullName"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Role         int          `json:"role"`
}

type PersonalInfo struct {
	HasOrganisation bool         `json:"hasOrganisation"`
	Organisation    Organisation `json:"organisation"`
}

type Organisation struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
