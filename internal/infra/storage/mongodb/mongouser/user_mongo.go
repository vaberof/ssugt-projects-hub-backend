package mongouser

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID
	Role         string
	FullName     string
	Email        string
	Password     string
	RegisteredAt time.Time
}
