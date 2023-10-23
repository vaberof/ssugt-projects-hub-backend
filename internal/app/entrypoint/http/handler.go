package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects/internal/domain/auth"
)

type Handler struct {
	authService auth.AuthService
}

func NewHandler(authService auth.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) InitRoutes(router *gin.Engine) *gin.Engine {
	apiV1 := router.Group("/api/v1")

	auth := apiV1.Group("/auth")
	auth.POST("/login", h.Login)

	return router
}
