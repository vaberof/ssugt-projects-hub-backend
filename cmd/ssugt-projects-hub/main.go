package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	httproutes "github.com/vaberof/ssugt-projects-hub-backend/internal/app/entrypoint/http"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/auth"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/project"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/domain/user"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/infra/storage/mongodb/mongoproject"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/infra/storage/mongodb/mongouser"
	"github.com/vaberof/ssugt-projects-hub-backend/internal/service/upload"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/database/mongodb"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/httpserver"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"os"
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
	authService := auth.NewAuthService(appConfig.AuthConfig, userService)
	projectService := project.NewProjectService(projectStorage)

	uploadService := upload.NewUploadService()

	httpHandler := httproutes.NewHandler(authService, projectService, uploadService, logger)

	appServer := httpserver.New(&appConfig.Server, logger)

	appServer.Server = httpHandler.InitRoutes(appServer.Server, logger)

	starter := appServer.StartAsync()

	<-(*starter)
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
