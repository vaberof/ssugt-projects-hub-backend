package files

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ssugt-projects-hub/models"
)

const (
	projectFilesBucket = "project_files"
)

type Repository interface {
	Save(ctx context.Context, files []models.ProjectFile) error
	GetByProjectId(ctx context.Context, projectId int) ([]models.ProjectFile, error)
}

type repositoryImpl struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repositoryImpl{db: db}
}

func (r repositoryImpl) Save(ctx context.Context, files []models.ProjectFile) error {
	// Создаём GridFS bucket с префиксом project_files
	bucket, err := gridfs.NewBucket(
		r.db,
		options.GridFSBucket().SetName(projectFilesBucket),
	)
	if err != nil {
		return err
	}

	for _, f := range files {
		// Декодируем содержимое из Base64
		data, err := base64.StdEncoding.DecodeString(string(f.Content))
		if err != nil {
			return err
		}

		// Готовим metadata: связываем файл с projectId, сохраняем тип, имя и время загрузки
		metadata := bson.M{
			"projectId":  f.ProjectId,
			"type":       f.Type,
			"name":       f.Name,
			"uploadedAt": f.UploadedAt,
		}
		uploadOpts := options.GridFSUpload().SetMetadata(metadata)

		// Загружаем потоком (имя файла используется как filename в GridFS)
		_, err = bucket.UploadFromStream(f.Name, bytes.NewReader(data), uploadOpts)
		if err != nil {
			return err
		}
	}

	log.Println("files saved")
	return nil
}

func (r repositoryImpl) GetByProjectId(ctx context.Context, projectId int) ([]models.ProjectFile, error) {
	bucket, err := gridfs.NewBucket(
		r.db,
		options.GridFSBucket().SetName(projectFilesBucket),
	)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"metadata.projectId": projectId}
	cursor, err := bucket.Find(filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []models.ProjectFile
	for cursor.Next(ctx) {
		var fileDoc bson.M
		if err := cursor.Decode(&fileDoc); err != nil {
			return nil, err
		}

		// Извлекаем ID файла (ObjectID), имя, metadata
		fileID, ok := fileDoc["_id"].(primitive.ObjectID)
		if !ok {
			continue
		}
		filename, _ := fileDoc["filename"].(string)
		meta, _ := fileDoc["metadata"].(bson.M)

		// Извлекаем поля из metadata
		var ftype models.ProjectFileType
		if t, ok := meta["type"].(string); ok {
			ftype = models.ProjectFileType(t)
		}
		var uploadedAt time.Time
		if dt, ok := meta["uploadedAt"].(primitive.DateTime); ok {
			uploadedAt = dt.Time()
		}

		// Скачиваем содержимое в буфер
		var buf bytes.Buffer
		if _, err := bucket.DownloadToStream(fileID, &buf); err != nil {
			return nil, err
		}
		// Кодируем содержимое в Base64
		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

		result = append(result, models.ProjectFile{
			Id:         fileID,
			ProjectId:  projectId,
			Type:       ftype,
			Name:       filename,
			Content:    models.ProjectFileContentBase64(encoded),
			UploadedAt: uploadedAt,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
