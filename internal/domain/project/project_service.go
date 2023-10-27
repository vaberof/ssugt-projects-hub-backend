package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

type ProjectService interface {
	Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error)
	Get(id domain.ProjectId) (*Project, error)
	ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error)
}

type projectServiceImpl struct {
	projectStorage ProjectStorage
}

func NewProjectService(projectStorage ProjectStorage) ProjectService {
	return &projectServiceImpl{projectStorage: projectStorage}
}

func (service *projectServiceImpl) Create(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	switch projectType {
	case ProjectTypeStudentScientificConference:
		return service.createSSCProject(userId, projectType, authors, organization, director, projectTemplate, tags)
	case ProjectTypeLaboratory:
		return service.createLaboratoryProject(userId, projectType, authors, organization, director, projectTemplate, tags)
	default:
		return "", errors.New(fmt.Sprintf("unknown project type '%s'", projectType))
	}
}

func (service *projectServiceImpl) Get(id domain.ProjectId) (*Project, error) {
	return service.projectStorage.Get(id)
}

func (service *projectServiceImpl) ListByFilters(userId domain.UserId, projectType domain.ProjectType, organizationName string, tags []string) ([]*Project, error) {
	return service.projectStorage.ListByFilters(userId, projectType, organizationName, tags)
}

func (service *projectServiceImpl) createSSCProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	var sscProjectTemplate StudentScientificConferenceProjectTemplate
	err := json.Unmarshal(projectTemplate, &sscProjectTemplate)
	if err != nil {
		return "", err
	}
	return service.projectStorage.CreateSSCProject(userId, projectType, authors, organization, director, &sscProjectTemplate, tags)
}

func (service *projectServiceImpl) createLaboratoryProject(userId domain.UserId, projectType domain.ProjectType, authors []*Author, organization *Organization, director *Director, projectTemplate ProjectTemplate, tags []string) (domain.ProjectId, error) {
	var laboratoryProjectTemplate LaboratoryProjectTemplate
	err := json.Unmarshal(projectTemplate, &laboratoryProjectTemplate)
	if err != nil {
		return "", err
	}
	return service.projectStorage.CreateLaboratoryProject(userId, projectType, authors, organization, director, &laboratoryProjectTemplate, tags)
}
