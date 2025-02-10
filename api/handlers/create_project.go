package handlers

import (
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	projectservice "ssugt-projects-hub/service/project"
)

func CreateProjectHandler(logs *logs.Logs, projectService projectservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("create-project-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Error(err.Error())
			xhttp.Forbidden(w)
			return
		}

		var project models.Project
		if err := xhttp.ReadRequestJson(r, &project); err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		response, err := projectService.Create(ctx, project)
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
