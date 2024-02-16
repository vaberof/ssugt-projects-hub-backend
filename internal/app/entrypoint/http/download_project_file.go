package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
	"path/filepath"
)

type DownloadProjectFileResponse struct {
	Message string `json:"message"`
}

func (handler *Handler) DownloadProjectFile(ctx *gin.Context) {
	projectId := ctx.Query("projectId")
	if projectId == "" {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "missing required query parameter 'projectId'"))
		return
	}

	fileName := ctx.Query("fileName")
	if fileName == "" {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "missing required query parameter 'fileName'"))
		return
	}

	joinedFilepath := filepath.Join(handler.projectService.GetUploadsPath(), fileName)

	if err := handler.projectService.IsFileExists(joinedFilepath); err != nil {
		ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, err.Error()))
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filepath="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")

	ctx.File(joinedFilepath)

	payload, _ := json.Marshal(DownloadProjectFileResponse{
		Message: fmt.Sprintf("file '%s' successfully downloaded", fileName),
	})

	ctx.JSON(http.StatusBadRequest, apiv1.Success(payload))
}
