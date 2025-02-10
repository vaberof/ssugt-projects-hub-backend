package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	projectservice "ssugt-projects-hub/service/project"
	"strconv"
)

func GetProjectByIdHandler(logs *logs.Logs, projectService projectservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("get-project-by-id-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
		}

		vars := mux.Vars(r)
		id := vars["id"]

		convId, err := strconv.Atoi(id)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
		}

		response, err := projectService.GetById(ctx, convId)
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
