package handlers

import (
	"context"
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	authservice "ssugt-projects-hub/service/auth"
)

func VerifyEmailHandler(logs *logs.Logs, authService authservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("verify-email-handler")

		var verifyEmailParams models.VerifyEmail
		if err := xhttp.ReadRequestJson(r, &verifyEmailParams); err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		err := authService.VerifyEmail(context.Background(), verifyEmailParams.Email, verifyEmailParams.Code)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
