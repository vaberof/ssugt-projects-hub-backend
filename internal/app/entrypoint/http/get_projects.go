package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
	"strings"
)

type getProjectsResponseBody struct {
	Projects []*views.ProjectResponsePayload `json:"projects"`
}

func (handler *Handler) GetProjects(ctx *gin.Context) {
	userIdParam := ctx.Query("userId")
	projectTypeParam := ctx.Query("projectType")
	organizationNameParam := ctx.Query("organizationName")
	tagsParam := ctx.Query("tags")

	var tagsValues []string

	if tagsParam != "" {
		tagsValues = strings.Split(tagsParam, ",")
	}

	projects, err := handler.projectService.ListByFilters(domain.UserId(userIdParam), domain.ProjectType(projectTypeParam), organizationNameParam, tagsValues)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	projectResponsePayloads, err := views.FromDomainProjects(projects)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	payload, _ := json.Marshal(getProjectsResponseBody{
		Projects: projectResponsePayloads,
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
