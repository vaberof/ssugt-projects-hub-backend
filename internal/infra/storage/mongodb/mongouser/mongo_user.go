package mongouser

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Role         string             `bson:"role"`
	FullName     string             `bson:"fullName"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	RegisteredAt time.Time          `bson:"registered_at"`
}
