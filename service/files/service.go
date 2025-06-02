package files

import (
	"context"
	"ssugt-projects-hub/database/mongo/files"
	"ssugt-projects-hub/models"
)

type Service interface {
	Save(ctx context.Context, files []models.ProjectFile) error
	GetByProjectId(ctx context.Context, projectId int) ([]models.ProjectFile, error)
	Update(ctx context.Context, projectId int, files []models.ProjectFile) error
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

func (s serviceImpl) Update(ctx context.Context, projectId int, files []models.ProjectFile) error {
	err := s.filesRepo.DeleteByProjectId(ctx, projectId)
	if err != nil {
		return err
	}
	return s.Save(ctx, files)
}
