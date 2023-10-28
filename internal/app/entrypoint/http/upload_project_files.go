package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const projectUploadsRelativePath = "public\\uploads\\projects"

func (handler *Handler) UploadProjectFiles(ctx *gin.Context) {
	projectId := ctx.Query("projectId")

	if projectId == "" {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "missing required query parameter 'projectId'"))
		return
	}

	projectRootDirectory, _ := os.Getwd()

	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	/*for _, file := range files {
		mimeType := file.Header.Get("Content-Type")
		//TODO: add pdf, pptx
		if mimeType != MIMETypeJpeg && mimeType != MIMETypePng {
			ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, fmt.Sprintf("unsupported MIMEType=%s", mimeType)))
			return
		}
	}*/

	pathToUploadsProjectDirectory := projectRootDirectory + "\\" + projectUploadsRelativePath + "\\" + projectId

	if _, err := os.Stat(pathToUploadsProjectDirectory); err != nil {
		err = os.Mkdir(pathToUploadsProjectDirectory, os.ModePerm)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
			return
		}
	}

	for _, file := range files {
		fileExtension := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))

		//newFileName := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v%s", time.Now().Unix(), fileExtension)

		filePath := pathToUploadsProjectDirectory + "\\" + originalFileName + fileExtension
		out, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		readerFile, _ := file.Open()
		_, err = io.Copy(out, readerFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message:": "uploaded successfully"})
}
