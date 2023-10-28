package mongoproject

import (
	"encoding/json"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

func (storage *MongoProjectStorage) toDomainProjects(mongoProjects []*Project) ([]*project.Project, error) {
	domainProjects := make([]*project.Project, len(mongoProjects))
	for i := 0; i < len(domainProjects); i++ {
		domainProject, err := storage.toDomainProject(mongoProjects[i])
		if err != nil {
			return nil, err
		}
		domainProjects[i] = domainProject
	}
	return domainProjects, nil
}

func (storage *MongoProjectStorage) toDomainProject(mongoProject *Project) (*project.Project, error) {
	var domainProject project.Project

	if mongoProject.ProjectType == project.ProjectTypeStudentScientificConference {
		projectTemplateBytes, err := json.Marshal(mongoProject.SscProjectTemplate)
		if err != nil {
			return nil, err
		}
		domainProject.Template = projectTemplateBytes
	} else if mongoProject.ProjectType == project.ProjectTypeLaboratory {
		projectTemplateBytes, err := json.Marshal(mongoProject.LaboratoryProjectTemplate)
		if err != nil {
			return nil, err
		}
		domainProject.Template = projectTemplateBytes
	}

	domainProject.Id = domain.ProjectId(mongoProject.Id.Hex())
	domainProject.UserId = domain.UserId(mongoProject.UserId.Hex())
	domainProject.ProjectType = domain.ProjectType(mongoProject.ProjectType)
	domainProject.Authors = storage.toDomainAuthors(mongoProject.Authors)
	domainProject.Organization = storage.toDomainOrganization(mongoProject.Organization)
	domainProject.Director = storage.toDomainDirector(mongoProject.Director)
	domainProject.Tags = mongoProject.Tags
	domainProject.CreatedAt = mongoProject.CreatedAt
	domainProject.ModifiedAt = mongoProject.ModifiedAt

	return &domainProject, nil
}

func (storage *MongoProjectStorage) toDomainAuthors(mongoAuthors []*Author) []*project.Author {
	authors := make([]*project.Author, len(mongoAuthors))
	for i := 0; i < len(authors); i++ {
		authors[i] = storage.toDomainAuthor(mongoAuthors[i])
	}
	return authors
}

func (storage *MongoProjectStorage) toDomainAuthor(mongoAuthor *Author) *project.Author {
	return &project.Author{
		FullName: mongoAuthor.FullName,
		Degree:   mongoAuthor.Degree,
		Course:   mongoAuthor.Course,
		Group:    mongoAuthor.Group,
	}
}

func (storage *MongoProjectStorage) toDomainDirector(mongoDirector *Director) *project.Director {
	return &project.Director{
		FullName: mongoDirector.FullName,
		Email:    mongoDirector.Email,
		Phone:    mongoDirector.Phone,
	}
}

func (storage *MongoProjectStorage) toDomainOrganization(mongoOrganization *Organization) *project.Organization {
	return &project.Organization{
		Name:    mongoOrganization.Name,
		Address: mongoOrganization.Address,
	}
}

func (storage *MongoProjectStorage) toDomainSscProjectTemplate(mongoSscProjectTemplate *SSCProjectTemplate) *project.StudentScientificConferenceProjectTemplate {
	return &project.StudentScientificConferenceProjectTemplate{
		Title:            mongoSscProjectTemplate.Title,
		Object:           mongoSscProjectTemplate.Object,
		Summary:          mongoSscProjectTemplate.Summary,
		Cost:             mongoSscProjectTemplate.Cost,
		DevelopingStage:  mongoSscProjectTemplate.DevelopingStage,
		RealizationTerm:  mongoSscProjectTemplate.RealizationTerm,
		ApplicationScope: mongoSscProjectTemplate.ApplicationScope,
	}
}
