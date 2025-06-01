package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ProjectFileType string

const (
	PatentProjectFileType      ProjectFileType = "Патент"
	PublicationProjectFileType ProjectFileType = "Публикация"
	ImageProjectFileType       ProjectFileType = "Изображение"

	OtherProjectFileType ProjectFileType = "Другое"
)

type ProjectFileContentBase64 string

type ProjectFile struct {
	Id         primitive.ObjectID       `json:"id" bson:"_id,omitempty"`
	ProjectId  int                      `json:"projectId" bson:"projectId"`
	UserId     int                      `json:"userId" bson:"userId"`
	Type       ProjectFileType          `json:"type" bson:"type"`
	Name       string                   `json:"name" bson:"name"`
	Content    ProjectFileContentBase64 `json:"content,omitempty" bson:"content,omitempty"` // Только для изображений
	UploadedAt time.Time                `json:"uploaded_at" bson:"uploaded_at"`
}
