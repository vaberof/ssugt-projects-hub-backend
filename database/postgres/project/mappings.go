package project

import (
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"github.com/samber/lo"
	"ssugt-projects-hub/models"
)

func mapToDbProject(project models.Project) DbProject {
	return DbProject{
		Id:            project.Id,
		UserId:        project.UserId,
		TypeId:        int(project.Type),
		Status:        string(project.Status),
		Attributes:    types.JSONText(project.Attributes),
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
		Collaborators: mapToDbCollaborators(project.Collaborators),
	}
}

func mapToDbCollaborators(collaborators []models.Collaborator) []DbCollaborator {
	dbCollaborators := make([]DbCollaborator, 0, len(collaborators))
	for _, collaborator := range collaborators {
		dbCollaborators = append(dbCollaborators, mapToDbCollaborator(collaborator))
	}
	return dbCollaborators
}

func mapToDbCollaborator(collaborator models.Collaborator) DbCollaborator {
	return DbCollaborator{
		Id:        collaborator.Id,
		UserId:    collaborator.UserId,
		ProjectId: collaborator.ProjectId,
		Role:      string(collaborator.Role),
	}
}

func mapDbProjectsWithDbCollaborators(dbProjects []DbProject, dbCollaborators []DbCollaborator) []DbProject {
	collaboratorsByProjectId := lo.GroupBy(dbCollaborators, func(collaborator DbCollaborator) int {
		return collaborator.ProjectId
	})

	return lo.Map(dbProjects, func(project DbProject, _ int) DbProject {
		project.Collaborators = collaboratorsByProjectId[project.Id]
		return project
	})
}

func mapToDbProjectReview(projectReview models.ProjectReview) DbProjectReview {
	return DbProjectReview{
		Id:         projectReview.Id,
		ProjectId:  projectReview.ProjectData.Id,
		ReviewedBy: projectReview.ReviewedBy,
		Status:     string(projectReview.Status),
		Comment:    projectReview.Comment,
		CreatedAt:  projectReview.CreatedAt,
	}
}

func mapToProjects(dbProjects []DbProject) []models.Project {
	projects := make([]models.Project, 0, len(dbProjects))
	for _, dbProject := range dbProjects {
		projects = append(projects, mapToProject(dbProject))
	}
	return projects
}

func mapToProject(dbProject DbProject) models.Project {
	return models.Project{
		Id:            dbProject.Id,
		UserId:        dbProject.UserId,
		Type:          models.ProjectType(dbProject.TypeId),
		Status:        models.ProjectStatus(dbProject.Status),
		Attributes:    json.RawMessage(dbProject.Attributes),
		CreatedAt:     dbProject.CreatedAt,
		UpdatedAt:     dbProject.UpdatedAt,
		Collaborators: mapToCollaborators(dbProject.Collaborators),
	}
}

func mapToCollaborators(dbCollaborators []DbCollaborator) []models.Collaborator {
	collaborators := make([]models.Collaborator, 0, len(dbCollaborators))
	for _, collaborator := range dbCollaborators {
		collaborators = append(collaborators, mapToCollaborator(collaborator))
	}
	return collaborators
}

func mapToCollaborator(dbCollaborator DbCollaborator) models.Collaborator {
	return models.Collaborator{
		Id:        dbCollaborator.Id,
		UserId:    dbCollaborator.UserId,
		ProjectId: dbCollaborator.ProjectId,
		Role:      models.CollaboratorRole(dbCollaborator.Role),
	}
}
