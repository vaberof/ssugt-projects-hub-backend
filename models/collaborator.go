package models

type CollaboratorRole string

const (
	DirectorCollaboratorRole          CollaboratorRole = "Руководитель"
	ArchitectCollaboratorRole         CollaboratorRole = "Архитектор"
	BackendDeveloperCollaboratorRole  CollaboratorRole = "Бэкенд-разработчик"
	FrontendDeveloperCollaboratorRole CollaboratorRole = "Фронтенд-разработчик"
	FullStackCollaboratorRole         CollaboratorRole = "Фуллстек-Разработчик"
	QAEngineerCollaboratorRole        CollaboratorRole = "Тестировщик"
	ProductManagerCollaboratorRole    CollaboratorRole = "Менеджер продукта"
	BusinessAnalyticCollaboratorRole  CollaboratorRole = "Бизнес-аналитик"

	OtherCollaboratorRole CollaboratorRole = "Другое"
)

type CollaboratorType int

type Collaborator struct {
	Id     int              `json:"id"`
	UserId int              `json:"userId"`
	Role   CollaboratorRole `json:"role"`
}
