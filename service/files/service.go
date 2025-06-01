package files

import (
	"context"
	"ssugt-projects-hub/database/mongo/files"
	"ssugt-projects-hub/models"
)

type Service interface {
	Save(ctx context.Context, files []models.ProjectFile) error
	GetByProjectId(ctx context.Context, projectId int) ([]models.ProjectFile, error)
}

type serviceImpl struct {
	filesRepo files.Repository
}

func NewService(filesRepo files.Repository) Service {
	return &serviceImpl{
		filesRepo: filesRepo,
	}
}

func (s serviceImpl) Save(ctx context.Context, files []models.ProjectFile) error {
	return s.filesRepo.Save(ctx, files)
}

func (s serviceImpl) GetByProjectId(ctx context.Context, projectId int) ([]models.ProjectFile, error) {
	return s.filesRepo.GetByProjectId(ctx, projectId)
}
