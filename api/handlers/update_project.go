package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	projectservice "ssugt-projects-hub/service/project"
	"strconv"
)

func UpdateProjectHandler(logs *logs.Logs, projectService projectservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("update-project-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Error(err.Error())
			xhttp.Forbidden(w)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		convId, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
		}

		var project models.Project
		if err := xhttp.ReadRequestJson(r, &project); err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}
		project.Id = convId

		if auth.UserIdFromContext(ctx) != project.UserId {
			// TODO: сейчас только создатель может обновить проект
			log.Debug("Обновление не своего проекта")
			xhttp.Forbidden(w)
			return
		}

		response, err := projectService.Update(ctx, project)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
