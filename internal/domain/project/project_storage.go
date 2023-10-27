package project

import "github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"

type ProjectStorage interface {
	CreateSSCProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, template *StudentScientificConferenceProjectTemplate, tags []string) (domain.ProjectId, error)
	CreateLaboratoryProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, template *LaboratoryProjectTemplate, tags []string) (domain.ProjectId, error)
	Get(id domain.ProjectId) (*Project, error)
	ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error)
}
