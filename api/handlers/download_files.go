package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	fileservice "ssugt-projects-hub/service/files"
	"strconv"
)

type downloadFilesResponse struct {
	Files []models.ProjectFile `json:"files"`
}

func DownloadFilesHandler(logs *logs.Logs, fileService fileservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("download-files-handler")

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

		files, err := fileService.GetByProjectId(ctx, projectId)
		if err != nil {
			log.Error("Ошибка при получении файлов:", err)
			xhttp.BadRequest(w)
			return
		}

		response := downloadFilesResponse{Files: files}
		if err := xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
