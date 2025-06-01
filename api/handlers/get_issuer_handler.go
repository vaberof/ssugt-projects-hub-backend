package handlers

import (
	"net/http"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
)

type getIssuerResponse struct {
	Issuer int `json:"issuer"`
}

func GetIssuerHandler(logs *logs.Logs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("get-issuer-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Error(err.Error())
			xhttp.Forbidden(w)
			return
		}

		userId := auth.UserIdFromContext(ctx)

		response := getIssuerResponse{
			Issuer: userId,
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
