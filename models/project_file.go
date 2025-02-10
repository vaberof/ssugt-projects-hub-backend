package models

type ProjectFileType string

const (
	PatentProjectFileType      ProjectFileType = "Патент"
	PublicationProjectFileType ProjectFileType = "Публикация"
	ImageProjectFileType       ProjectFileType = "Изображение"

	OtherProjectFileType ProjectFileType = "Другое"
)

type ProjectFileContentBase64 string

type ProjectFile struct {
	Id      string                   `json:"id"`
	Type    ProjectFileType          `json:"type"`
	Name    string                   `json:"name"`
	Content ProjectFileContentBase64 `json:"content"` // Только для изображений
}
