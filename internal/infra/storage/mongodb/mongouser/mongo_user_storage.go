package mongouser

import (
	"context"
	"github.com/vaberof/ssugt-projects/internal/domain/user"
	"github.com/vaberof/ssugt-projects/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionUser = "users"

type MongoUserStorage struct {
	collection *mongo.Collection
}

func NewMongoUserStorage(db *mongo.Database) *MongoUserStorage {
	return &MongoUserStorage{collection: db.Collection(collectionUser)}
}

func (storage *MongoUserStorage) Get(id domain.UserId) (*user.User, error) {
	var user User

	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	err = storage.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return storage.fromMongoUser(&user), nil
}

func (storage *MongoUserStorage) GetByEmail(email domain.Email) (*user.User, error) {
	var user User

	err := storage.collection.FindOne(context.Background(), bson.M{"email": string(email)}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return storage.fromMongoUser(&user), nil
}

func (storage *MongoUserStorage) fromMongoUser(mongoUser *User) *user.User {
	return &user.User{
		Id:           domain.UserId(mongoUser.Id.String()),
		Role:         domain.Role(mongoUser.Role),
		FullName:     domain.FullName(mongoUser.FullName),
		Password:     domain.Password(mongoUser.Password),
		RegisteredAt: mongoUser.RegisteredAt,
	}
}
