package handlers

import (
	"bytes"
	"encoding/base64"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"path/filepath"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	fileservice "ssugt-projects-hub/service/files"
	"strconv"
	"strings"
	"time"
)

type updateFilesResponse struct {
	Message string `json:"message"`
}

func UpdateFilesHandler(logs *logs.Logs, fileService fileservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("update-files-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
			xhttp.Forbidden(w)
		}

		vars := mux.Vars(r)
		id := vars["id"]

		projectId, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Error("Не удалось распарсить форму:", err)
			xhttp.BadRequest(w)
			return
		}

		form := r.MultipartForm
		fileHeaders := form.File["files"]
		if len(fileHeaders) == 0 {
			log.Error("Не переданы файлы для загрузки")
			xhttp.BadRequest(w)
			return
		}

		var toSave []models.ProjectFile
		for _, fh := range fileHeaders {
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				continue
			}

			file, err := fh.Open()
			if err != nil {
				log.Error("Не удалось открыть файл:", err)
				continue
			}
			defer file.Close()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, file); err != nil {
				log.Error("Не удалось прочитать файл:", err)
				continue
			}

			encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

			toSave = append(toSave, models.ProjectFile{
				ProjectId:  projectId,
				Type:       models.ImageProjectFileType,
				Name:       fh.Filename,
				Content:    models.ProjectFileContentBase64(encoded),
				UploadedAt: time.Now(),
			})
		}

		if len(toSave) == 0 {
			log.Error("Нет валидных файлов для обновления")
			xhttp.BadRequest(w)
			return
		}

		if err := fileService.Update(ctx, projectId, toSave); err != nil {
			log.Error("Ошибка при обновлении файлов:", err)
			xhttp.BadRequest(w)
			return
		}

		response := updateFilesResponse{Message: "Successfully uploaded files"}
		if err := xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
