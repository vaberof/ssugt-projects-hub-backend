package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/httpserver"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"mime/multipart"
	"net/http"
)

type uploadProjectFilesResponseBody struct {
	UploadedFiles []project.Filename `json:"uploaded_files"`
}

func (handler *Handler) UploadProjectFiles(ctx *gin.Context) {
	projectId := ctx.Query("projectId")

	if projectId == "" {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "missing required query parameter 'projectId'"))
		return
	}

	if err := ctx.Request.ParseMultipartForm(httpserver.MaxMultipartMemory); err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, fmt.Sprintf("files size limit exceeded. maxSize=%d MB", httpserver.MaxMultipartMemory)))
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	files := form.File["files"]

	uploadedFiles, err := handler.projectService.SaveFiles(domain.ProjectId(projectId), files)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message:": err.Error()})
		return
	}

	payload, _ := json.Marshal(&uploadProjectFilesResponseBody{
		UploadedFiles: uploadedFiles,
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}

func (handler *Handler) validateFiles(files []*multipart.FileHeader) error {
	for _, file := range files {
		mimeType := file.Header.Get("Content-Type")
		//TODO: add doc, docx, pdf, pptx
		if mimeType != MIMETypeJpeg && mimeType != MIMETypePng {
			return errors.New(fmt.Sprintf("unsupported MIME type: %s\n", mimeType))
		}
	}
	return nil
}
