package project

import "github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"

const (
	FileTypeImage = "image"
	FileTypeOther = "other"
)

type FileType string

type ProjectFileContent []byte

type ProjectFile struct {
	ProjectId   domain.ProjectId   `json:"project_id,omitempty"`
	Type        FileType           `json:"type,omitempty"`
	ContentType string             `json:"content_type,omitempty"`
	Name        string             `json:"name,omitempty"`
	Size        int64              `json:"size,omitempty"`
	Content     ProjectFileContent `json:"content,omitempty"`
}
