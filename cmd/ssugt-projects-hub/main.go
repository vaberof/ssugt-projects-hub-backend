package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	httproutes "github.com/vaberof/ssugt-projects-hub-backend/internal/app/entrypoint/http"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/auth"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/user"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/infra/storage/mongodb/mongoproject"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/infra/storage/mongodb/mongouser"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/database/mongodb"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/httpserver"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()

	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	logger := logs.New(os.Stdout, nil)

	appConfig := getAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	managedDatabase, err := mongodb.New(&appConfig.Database)
	if err != nil {
		panic(err)
	}

	userStorage := mongouser.NewMongoUserStorage(managedDatabase.Db)
	projectStorage := mongoproject.NewMongoProjectStorage(managedDatabase.Db)

	userService := user.NewUserService(userStorage)
	authService := auth.NewAuthService(userService, appConfig.AuthConfig)
	projectService := project.NewProjectService(projectStorage, appConfig.ProjectService)

	httpHandler := httproutes.NewHandler(authService, projectService, logger)

	appServer := httpserver.New(&appConfig.Server, logger)

	router := appServer.LoadGinEngineFromHTTPHandler()
	if router == nil {
		panic("Failed to load gin router")
	}

	appServer.Server.Handler = httpHandler.InitRoutes(router, logger)

	serverExitChannel := appServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-quitCh:
		logger.GetLogger().Info("stopping application", slog.String("signal", sign.String()))

		gracefulShutdown(appServer, managedDatabase)
	case err := <-serverExitChannel:
		logger.GetLogger().Info("stopping application", slog.String("err", err.Error()))

		gracefulShutdown(appServer, managedDatabase)
	}
}

func gracefulShutdown(server *httpserver.AppServer, mongo *mongodb.ManagedDatabase) {
	// log
	serverShutdownTimeout := time.Second * 5
	databaseDisconnectTimeout := time.Second * 10

	emptyCtx := context.Background()

	ctx, cancel := context.WithTimeout(emptyCtx, serverShutdownTimeout)
	defer cancel()

	if err := server.Server.Shutdown(ctx); err != nil {
		// log
	}

	ctx, cancel = context.WithTimeout(emptyCtx, databaseDisconnectTimeout)
	defer cancel()

	if err := mongo.Disconnect(ctx); err != nil {
		// log
	}

	// log
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
