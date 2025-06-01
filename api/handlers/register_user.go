package handlers

import (
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	authservice "ssugt-projects-hub/service/auth"
)

type registerUserResponse struct {
	Message string `json:"message"`
}

func RegisterUserHandler(logs *logs.Logs, authService authservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("register-user-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
		}

		var user models.User
		if err := xhttp.ReadRequestJson(r, &user); err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		_, err = authService.Register(ctx, user)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		response := registerUserResponse{
			Message: getResponseMessage(err),
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}

func getResponseMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return "ok"
}
