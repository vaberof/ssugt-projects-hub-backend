package http

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/protocols/apiv1"
	"net/http"
)

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	AccessToken string `json:"access_token"`
}

func (handler *Handler) Login(ctx *gin.Context) {
	var loginReqBody loginRequestBody
	if err := ctx.Bind(&loginReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "invalid request body"))
		return
	}

	accessToken, err := handler.authService.Login(domain.Email(loginReqBody.Email), domain.Password(loginReqBody.Password))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))
		return
	}

	accessTokenString := string(*accessToken)

	payload, _ := json.Marshal(loginResponseBody{AccessToken: accessTokenString})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
