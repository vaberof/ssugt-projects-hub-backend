package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
	"os"
)

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

	projectRootDirectory, _ := os.Getwd()

	pathToUploadsProjectDirectory := projectRootDirectory + "\\" + projectUploadsRelativePath + "\\" + projectId + "\\" + fileName

	if _, err := os.Stat(pathToUploadsProjectDirectory); err != nil {
		ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, err.Error()))
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")

	ctx.File(pathToUploadsProjectDirectory)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message:": "downloaded file successfully"})
}
