package handlers

import (
	"net/http"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	authservice "ssugt-projects-hub/service/auth"
)

type isAdminResponse struct {
	IsAdmin bool `json:"isAdmin"`
}

func IsAdminHandler(logs *logs.Logs, authService authservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("is-admin-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Error(err.Error())
			xhttp.Forbidden(w)
			return
		}

		userId := auth.UserIdFromContext(ctx)

		isAdmin, err := authService.IsAdmin(ctx, userId)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		response := isAdminResponse{
			IsAdmin: isAdmin,
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
