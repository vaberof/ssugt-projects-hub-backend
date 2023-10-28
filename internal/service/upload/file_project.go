package upload

import "github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"

type ProjectFileContent []byte

type ProjectFile struct {
	Id          int64
	ProjectId   domain.ProjectId
	Type        FileType
	ContentType string
	Name        string
	Size        int64
	Content     ProjectFileContent
}
