package mongoproject

import "github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"

func (storage *MongoProjectStorage) toMongoAuthors(domainAuthors []*project.Author) []*Author {
	authors := make([]*Author, len(domainAuthors))
	for i := 0; i < len(authors); i++ {
		authors[i] = storage.toMongoAuthor(domainAuthors[i])
	}
	return authors
}

func (storage *MongoProjectStorage) toMongoAuthor(domainAuthor *project.Author) *Author {
	return &Author{
		FullName: domainAuthor.FullName,
		Degree:   domainAuthor.Degree,
		Course:   domainAuthor.Course,
		Group:    domainAuthor.Group,
	}
}

func (storage *MongoProjectStorage) toMongoDirector(domainDirector *project.Director) *Director {
	return &Director{
		FullName: domainDirector.FullName,
		Email:    domainDirector.Email,
		Phone:    domainDirector.Phone,
	}
}

func (storage *MongoProjectStorage) toMongoOrganization(domainOrganization *project.Organization) *Organization {
	return &Organization{
		Name:    domainOrganization.Name,
		Address: domainOrganization.Address,
	}
}

func (storage *MongoProjectStorage) toMongoSscProjectTemplate(domainSscProjectTemplate *project.StudentScientificConferenceProjectTemplate) *SSCProjectTemplate {
	return &SSCProjectTemplate{
		Title:            domainSscProjectTemplate.Title,
		Object:           domainSscProjectTemplate.Object,
		Summary:          domainSscProjectTemplate.Summary,
		Cost:             domainSscProjectTemplate.Cost,
		DevelopingStage:  domainSscProjectTemplate.DevelopingStage,
		RealizationTerm:  domainSscProjectTemplate.RealizationTerm,
		ApplicationScope: domainSscProjectTemplate.ApplicationScope,
	}
}

func (storage *MongoProjectStorage) toMongoLaboratoryProjectTemplate(domainLaboratoryProjectTemplate *project.LaboratoryProjectTemplate) *LaboratoryProjectTemplate {
	return &LaboratoryProjectTemplate{
		LaboratoryName:   domainLaboratoryProjectTemplate.LaboratoryName,
		Title:            domainLaboratoryProjectTemplate.Title,
		Object:           domainLaboratoryProjectTemplate.Object,
		Summary:          domainLaboratoryProjectTemplate.Summary,
		Problematic:      domainLaboratoryProjectTemplate.Problematic,
		Solution:         domainLaboratoryProjectTemplate.Solution,
		Functionality:    domainLaboratoryProjectTemplate.Functionality,
		TechnologyStack:  domainLaboratoryProjectTemplate.TechnologyStack,
		Advantages:       domainLaboratoryProjectTemplate.Advantages,
		Cost:             domainLaboratoryProjectTemplate.Cost,
		DevelopingStage:  domainLaboratoryProjectTemplate.DevelopingStage,
		RealizationTerm:  domainLaboratoryProjectTemplate.RealizationTerm,
		ApplicationScope: domainLaboratoryProjectTemplate.ApplicationScope,
	}
}
