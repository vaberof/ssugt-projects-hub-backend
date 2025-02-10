package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	userservice "ssugt-projects-hub/service/user"
)

func GetUserByEmailHandler(logs *logs.Logs, userService userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("get-user-by-email-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
		}

		vars := mux.Vars(r)
		email := vars["email"]

		response, err := userService.GetByEmail(ctx, email)
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
