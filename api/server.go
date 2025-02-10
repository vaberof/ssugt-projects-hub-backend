package api

import (
	"context"
	"fmt"
	gorillahandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"ssugt-projects-hub/api/handlers"
	"ssugt-projects-hub/config"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/service/auth"
	"ssugt-projects-hub/service/project"
	"ssugt-projects-hub/service/user"
)

func NewServer(
	ctx context.Context,
	log *logs.Logs,
	authService auth.Service,
	projectService project.Service,
	userService user.Service) *http.Server {

	route := mux.NewRouter()
	route.Use(jsonContentTypeMiddleware)
	route.Use(gorillahandler.RecoveryHandler(gorillahandler.PrintRecoveryStack(true)))

	// Статус сервиса
	route.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	addAuthHandlers(route, log, authService)
	addProjectsHandlers(route, log, projectService)

	addUsersHandlers(route, log, userService)

	//addProjectsReviewsHandlers(route, log, userService)

	// CORS
	headersOk := gorillahandler.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := gorillahandler.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := gorillahandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	return &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port()),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},

		Handler: gorillahandler.CORS(headersOk, originsOk, methodsOk)(route),
	}
}

func addAuthHandlers(route *mux.Router, log *logs.Logs, authService auth.Service) {
	route.HandleFunc("/auth/register", handlers.RegisterUserHandler(log, authService)).Methods(http.MethodPost)
	route.HandleFunc("/auth/login", handlers.LoginUserHandler(log, authService)).Methods(http.MethodPost)
}

func addProjectsHandlers(route *mux.Router, log *logs.Logs, projectService project.Service) {
	route.HandleFunc("/projects", handlers.CreateProjectHandler(log, projectService)).Methods(http.MethodPost)
	route.HandleFunc("/projects/{id}", handlers.GetProjectByIdHandler(log, projectService)).Methods(http.MethodGet)
	route.HandleFunc("/projects/{userId}", handlers.GetProjectByUserIdHandler(log, projectService)).Methods(http.MethodGet)
	route.HandleFunc("/projects/search", handlers.SearchProjectHandler(log, projectService)).Methods(http.MethodPost)
	route.HandleFunc("/projects", handlers.UpdateProjectHandler(log, projectService)).Methods(http.MethodPut)
}

func addUsersHandlers(route *mux.Router, log *logs.Logs, userService user.Service) {
	route.HandleFunc("/users/{email}", handlers.GetUserByEmailHandler(log, userService)).Methods(http.MethodGet)
}

func addProjectsReviewsHandlers(route *mux.Router, log *logs.Logs, authService auth.Service) {
	//route.HandleFunc("/reviews/projects", handlers.(log, authService)).Methods(http.MethodPost)
	//route.HandleFunc("/reviews/projects/{id}", handlers.(log, authService)).Methods(http.MethodGet)
	//route.HandleFunc("/reviews/projects/search", handlers.(log, authService)).Methods(http.MethodPost)
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
