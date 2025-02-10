package handlers

import (
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	authservice "ssugt-projects-hub/service/auth"
)

func LoginUserHandler(logs *logs.Logs, authService authservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("login-user-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
		}

		var loginUserRequestParams models.LoginUserRequestParams
		if err := xhttp.ReadRequestJson(r, &loginUserRequestParams); err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		accessToken, err := authService.Login(ctx, loginUserRequestParams)
		if err != nil || accessToken == nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		response := models.LoginUserResponse{
			AccessToken: string(*accessToken),
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
