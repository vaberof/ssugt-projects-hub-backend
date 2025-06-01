package main

import (
	"context"
	"errors"
	"net/http"
	"ssugt-projects-hub/api"
	"ssugt-projects-hub/config"
	"ssugt-projects-hub/database/mongo/cache"
	filerepo "ssugt-projects-hub/database/mongo/files"
	projectrepo "ssugt-projects-hub/database/postgres/project"
	userrepo "ssugt-projects-hub/database/postgres/user"
	"ssugt-projects-hub/pkg/database/mongo"
	"ssugt-projects-hub/pkg/database/postgres"
	"ssugt-projects-hub/pkg/logging/logs"
	authservice "ssugt-projects-hub/service/auth"
	fileservice "ssugt-projects-hub/service/files"
	projectservice "ssugt-projects-hub/service/project"
	emailservice "ssugt-projects-hub/service/sender/email"
	userservice "ssugt-projects-hub/service/user"
	"time"
)

type App struct {
	mainCtx context.Context
	log     *logs.Logs
	server  *http.Server

	userRepository    userrepo.Repository
	projectRepository projectrepo.Repository
	fileRepository    filerepo.Repository

	cache cache.Cache

	authService    authservice.Service
	userService    userservice.Service
	emailService   emailservice.Service
	projectService projectservice.Service
	fileService    fileservice.Service
}

func NewApp(mainCtx context.Context, log *logs.Logs) *App {
	return &App{
		mainCtx: mainCtx,
		log:     log,
	}
}

func (a *App) initDatabases() {
	mongoCfg := mongo.MongoDatabaseConfig{
		AppName:  "",
		URI:      "mongodb://localhost:27017",
		Database: "ssugt-projects-hub",
	}

	mongoDb, err := mongo.New(&mongoCfg)
	if err != nil {
		panic(err)
	}

	postgresDb := postgres.NewPgx(context.Background(), config.PostgresConnection())

	a.userRepository = userrepo.NewRepository(postgresDb)
	a.projectRepository = projectrepo.NewRepository(postgresDb)
	a.fileRepository = filerepo.NewRepository(mongoDb.Db)

	a.cache = cache.NewMongoRepository(mongoDb.Db)

	err = runMigrations(postgresDb)
	if err != nil {
		a.log.GetLogger().Error("не смог прогнать миграции: %v", err)
		panic(err)
	}
}

func (a *App) initServices() {
	a.userService = userservice.NewService(a.log, a.userRepository)
	a.emailService = emailservice.NewService(a.log, emailservice.NewSmtpConfig())
	a.authService = authservice.NewService(a.log, a.userService, a.emailService, a.cache)
	a.projectService = projectservice.NewProjectService(a.log, a.projectRepository, a.fileRepository)
	a.fileService = fileservice.NewService(a.fileRepository)
}

func (a *App) initServer() {
	a.server = api.NewServer(a.mainCtx, a.log, a.authService, a.projectService, a.fileService, a.userService)
}

func (a *App) Start() {
	go func() {
		a.log.GetLogger().Debug("Сервер прослушивает %s порт", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.GetLogger().Error("Сервер неожиданно остановлен: %v", err)
		}
	}()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(a.mainCtx, 15*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.GetLogger().Error("Не удалось остановить сервер: %v", err)
	}

	a.log.GetLogger().Debug("Приложение успешно остановлено")
}
