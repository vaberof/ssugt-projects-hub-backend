package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/app/entrypoint/http/middleware"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/auth"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"log/slog"
)

type Handler struct {
	authService    auth.AuthService
	projectService project.ProjectService
	logger         *slog.Logger
}

func NewHandler(authService auth.AuthService, projectService project.ProjectService, logs *logs.Logs) *Handler {
	logger := logs.WithName("http-handler")
	return &Handler{authService: authService, projectService: projectService, logger: logger}
}

func (handler *Handler) InitRoutes(router *gin.Engine, logs *logs.Logs) *gin.Engine {
	apiV1 := router.Group("/api/v1")

	auth := apiV1.Group("/auth")
	auth.POST("/login", handler.Login)

	projects := apiV1.Group("/projects")
	projects.GET("/:id", handler.GetProject)
	projects.GET("/", handler.GetProjects)
	projects.POST("/", middleware.AuthMiddleware(handler.authService, logs), handler.CreateProject)

	return router
}
