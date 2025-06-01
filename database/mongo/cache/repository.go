package cache

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	emailConfirmationsCollection = "email_confirmations"
)

type Cache interface {
	Insert(ctx context.Context, emailConfirmation EmailConfirmation) error
	UpdateOrInsert(ctx context.Context, emailConfirmation EmailConfirmation) error
	Get(ctx context.Context, email string) (EmailConfirmation, error)
	DeleteByEmail(ctx context.Context, email string) error
}

type mongoCache struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Cache {
	return &mongoCache{
		collection: db.Collection(emailConfirmationsCollection),
	}
}

func (r *mongoCache) Insert(ctx context.Context, emailConfirmation EmailConfirmation) error {
	_, err := r.collection.InsertOne(ctx, emailConfirmation)
	return err
}

func (r *mongoCache) UpdateOrInsert(ctx context.Context, emailConfirmation EmailConfirmation) error {
	emailConfirmation.CreatedAt = time.Now().UTC()

	filter := bson.M{
		"Email": emailConfirmation.Email,
	}

	update := bson.M{
		"$set": emailConfirmation,
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (r *mongoCache) Get(ctx context.Context, email string) (EmailConfirmation, error) {
	var result EmailConfirmation
	err := r.collection.FindOne(ctx, bson.M{"Email": email}).Decode(&result)
	return result, err
}

func (r *mongoCache) DeleteByEmail(ctx context.Context, email string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"Email": email})
	return err
}
