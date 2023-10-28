package views

import (
	"encoding/json"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"time"
)

type ProjectResponsePayload struct {
	Id                 string                                      `json:"id"`
	UserId             string                                      `json:"user_id"`
	ProjectType        string                                      `json:"project_type"`
	Authors            []*AuthorResponse                           `json:"authors"`
	Organization       *OrganizationResponse                       `json:"organization"`
	Director           *DirectorResponse                           `json:"director"`
	SscTemplate        *StudentScientificConferenceProjectTemplate `json:"ssc_template,omitempty"`
	LaboratoryTemplate *LaboratoryProjectTemplate                  `json:"laboratory_template,omitempty"`
	Tags               []string                                    `json:"tags,omitempty"`
	Files              []*project.ProjectFile                      `json:"files,omitempty"` // TODO: смапить в респонс
	CreatedAt          time.Time                                   `json:"created_at"`
	ModifiedAt         time.Time                                   `json:"modified_at"`
}

type AuthorResponse struct {
	FullName string `json:"full_name"`
	Degree   string `json:"degree,omitempty"`
	Course   int    `json:"course,omitempty"`
	Group    string `json:"group"`
}

type OrganizationResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type DirectorResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func FromDomainProjects(domainProjects []*project.Project) ([]*ProjectResponsePayload, error) {
	projectResponsePayloads := make([]*ProjectResponsePayload, len(domainProjects))
	for i := 0; i < len(domainProjects); i++ {
		projectResponsePayload, err := FromDomainProject(domainProjects[i])
		if err != nil {
			return nil, err
		}
		projectResponsePayloads[i] = projectResponsePayload
	}
	return projectResponsePayloads, nil
}

func FromDomainProject(domainProject *project.Project) (*ProjectResponsePayload, error) {
	var getProjectRespPayload ProjectResponsePayload

	switch domainProject.ProjectType {
	case project.ProjectTypeStudentScientificConference:
		var projectTemplate StudentScientificConferenceProjectTemplate
		err := json.Unmarshal(domainProject.Template, &projectTemplate)
		if err != nil {
			return nil, err
		}
		getProjectRespPayload.SscTemplate = &projectTemplate

	case project.ProjectTypeLaboratory:
		var projectTemplate LaboratoryProjectTemplate
		err := json.Unmarshal(domainProject.Template, &projectTemplate)
		if err != nil {
			return nil, err
		}
		getProjectRespPayload.LaboratoryTemplate = &projectTemplate
	}

	getProjectRespPayload.Id = string(domainProject.Id)
	getProjectRespPayload.UserId = string(domainProject.UserId)
	getProjectRespPayload.ProjectType = string(domainProject.ProjectType)
	getProjectRespPayload.Authors = toAuthorsResponse(domainProject.Authors)
	getProjectRespPayload.Organization = toOrganizationResponse(domainProject.Organization)
	getProjectRespPayload.Director = toDirectorResponse(domainProject.Director)
	getProjectRespPayload.Tags = domainProject.Tags
	getProjectRespPayload.CreatedAt = domainProject.CreatedAt
	getProjectRespPayload.ModifiedAt = domainProject.ModifiedAt
	getProjectRespPayload.Files = domainProject.Files // TODO сменить тип

	return &getProjectRespPayload, nil
}

func toAuthorsResponse(domainAuthors []*project.Author) []*AuthorResponse {
	authors := make([]*AuthorResponse, len(domainAuthors))
	for i := 0; i < len(authors); i++ {
		authors[i] = toAuthorResponse(domainAuthors[i])
	}
	return authors
}

func toAuthorResponse(domainAuthor *project.Author) *AuthorResponse {
	return &AuthorResponse{
		FullName: domainAuthor.FullName,
		Degree:   domainAuthor.Degree,
		Course:   domainAuthor.Course,
		Group:    domainAuthor.Group,
	}
}

func toDirectorResponse(domainDirector *project.Director) *DirectorResponse {
	return &DirectorResponse{
		FullName: domainDirector.FullName,
		Email:    domainDirector.Email,
		Phone:    domainDirector.Phone,
	}
}

func toOrganizationResponse(domainOrganization *project.Organization) *OrganizationResponse {
	return &OrganizationResponse{
		Name:    domainOrganization.Name,
		Address: domainOrganization.Address,
	}
}
