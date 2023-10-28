package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
)

// TODO: map from/to view model

type createProjectRequestBody struct {
	ProjectType        domain.ProjectType                                  `json:"project_type"`
	Authors            []*project.Author                                   `json:"authors"`
	Organization       *project.Organization                               `json:"organization"`
	Director           *project.Director                                   `json:"director"`
	SscTemplate        *project.StudentScientificConferenceProjectTemplate `json:"ssc_template,omitempty"`
	LaboratoryTemplate *project.LaboratoryProjectTemplate                  `json:"laboratory_template,omitempty"`
	Tags               []string                                            `json:"tags"`
}

type createProjectResponseBody struct {
	ProjectId domain.ProjectId `json:"project_id"`
}

func (handler *Handler) CreateProject(ctx *gin.Context) {
	var createProjectReqBody createProjectRequestBody

	if err := ctx.Bind(&createProjectReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "invalid request body"))
		return
	}

	userId, err := handler.userIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, apiv1.Error(apiv1.CodeUnauthorized, fmt.Sprintf("unauthorized: %s", err.Error())))
		return
	}

	projectTemplateBytes, err := handler.getProjectTemplateBytes(&createProjectReqBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	projectId, err := handler.projectService.Create(*userId, createProjectReqBody.ProjectType, createProjectReqBody.Authors, createProjectReqBody.Organization, createProjectReqBody.Director, projectTemplateBytes, createProjectReqBody.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	payload, _ := json.Marshal(&createProjectResponseBody{
		ProjectId: projectId,
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}

func (handler *Handler) getProjectTemplateBytes(createProjectReqBody *createProjectRequestBody) ([]byte, error) {
	switch createProjectReqBody.ProjectType {
	case project.ProjectTypeStudentScientificConference:
		projectTemplateBytes, err := json.Marshal(&createProjectReqBody.SscTemplate)
		if err != nil {
			return nil, err
		}
		return projectTemplateBytes, nil
	case project.ProjectTypeLaboratory:
		projectTemplateBytes, err := json.Marshal(&createProjectReqBody.LaboratoryTemplate)
		if err != nil {
			return nil, err
		}
		return projectTemplateBytes, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown project type '%s'", createProjectReqBody.ProjectType))
	}
}
