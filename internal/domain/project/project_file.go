package project

import "github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"

const (
	FileTypePatent      = "patent"
	FileTypePublication = "publication"
	FileTypeImage       = "image"
	FileTypeOther       = "other"
)

type FileType string

type ContentType string

type Filename string

type FileContentBase64 string

type ProjectFile struct {
	ProjectId   domain.ProjectId  `json:"project_id,omitempty"`
	Type        FileType          `json:"type,omitempty"`
	ContentType ContentType       `json:"content_type,omitempty"`
	Name        Filename          `json:"name,omitempty"`
	Content     FileContentBase64 `json:"content,omitempty"`
}
