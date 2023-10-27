package mongoproject

import (
	"context"
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collectionProjects = "ssugt_projects"

type MongoProjectStorage struct {
	collection *mongo.Collection
}

func NewMongoProjectStorage(db *mongo.Database) *MongoProjectStorage {
	return &MongoProjectStorage{collection: db.Collection(collectionProjects)}
}

func (storage *MongoProjectStorage) CreateSSCProject(userId domain.UserId, projectType domain.ProjectType, authors []*project.Author, organization *project.Organization, director *project.Director, template *project.StudentScientificConferenceProjectTemplate, tags []string) (domain.ProjectId, error) {
	var project Project

	primitiveObjectUserId, err := primitive.ObjectIDFromHex(string(userId))
	if err != nil {
		return "", err
	}

	project.ProjectType = string(projectType)
	project.UserId = primitiveObjectUserId
	project.Authors = storage.toMongoAuthors(authors)
	project.Organization = storage.toMongoOrganization(organization)
	project.Director = storage.toMongoDirector(director)
	project.SscProjectTemplate = storage.toMongoSscProjectTemplate(template)
	project.Tags = tags
	project.CreatedAt = time.Now()
	project.ModifiedAt = time.Now()

	insertOneResult, err := storage.collection.InsertOne(context.Background(), &project)
	if err != nil {
		return "", err
	}

	primitiveObjectProjectId, ok := insertOneResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get _id of created ssc project")
	}

	projectId := domain.ProjectId(primitiveObjectProjectId.Hex())

	return projectId, nil
}

func (storage *MongoProjectStorage) CreateLaboratoryProject(userId domain.UserId, projectType domain.ProjectType, authors []*project.Author, organization *project.Organization, director *project.Director, template *project.LaboratoryProjectTemplate, tags []string) (domain.ProjectId, error) {
	var project Project

	primitiveObjectUserId, err := primitive.ObjectIDFromHex(string(userId))
	if err != nil {
		return "", err
	}

	project.ProjectType = string(projectType)
	project.UserId = primitiveObjectUserId
	project.Authors = storage.toMongoAuthors(authors)
	project.Organization = storage.toMongoOrganization(organization)
	project.Director = storage.toMongoDirector(director)
	project.LaboratoryProjectTemplate = storage.toMongoLaboratoryProjectTemplate(template)
	project.Tags = tags
	project.CreatedAt = time.Now()
	project.ModifiedAt = time.Now()

	insertOneResult, err := storage.collection.InsertOne(context.Background(), &project)
	if err != nil {
		return "", err
	}

	primitiveObjectProjectId, ok := insertOneResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get _id of created laboratory project")
	}

	projectId := domain.ProjectId(primitiveObjectProjectId.Hex())

	return projectId, nil
}

func (storage *MongoProjectStorage) Get(id domain.ProjectId) (*project.Project, error) {
	var project Project

	objectId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	err = storage.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&project)
	if err != nil {
		return nil, err
	}

	fmt.Printf("project: %+v\n", project.SscProjectTemplate)

	return storage.toDomainProject(&project)
}

func (storage *MongoProjectStorage) ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*project.Project, error) {
	filterQuery := make(bson.M)

	if userId != "" {
		objectUserId, err := primitive.ObjectIDFromHex(userId.String())
		if err != nil {
			return nil, err
		}
		filterQuery["user_id"] = objectUserId
	}

	if projectType != "" {
		filterQuery["project_type"] = projectType
	}

	if organizationName != "" {
		filterQuery["organization.name"] = organizationName
	}

	if len(tags) != 0 {
		filterQuery["tags"] = bson.M{"$in": tags}
	}

	var projects []*Project

	cursor, err := storage.collection.Find(context.Background(), filterQuery)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &projects)
	if err != nil {
		return nil, err
	}

	fmt.Println(projects)

	domainProjects, err := storage.toDomainProjects(projects)
	if err != nil {
		return nil, err
	}

	return domainProjects, nil
}
