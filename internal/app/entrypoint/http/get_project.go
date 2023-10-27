package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
)

type getProjectResponseBody struct {
	Project *views.ProjectResponsePayload `json:"project"`
}

func (handler *Handler) GetProject(ctx *gin.Context) {
	projectId := ctx.Param("id")

	domainProject, err := handler.projectService.Get(domain.ProjectId(projectId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	getProjectResponsePayload, err := views.FromDomainProject(domainProject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	payload, _ := json.Marshal(getProjectResponseBody{
		Project: getProjectResponsePayload,
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
