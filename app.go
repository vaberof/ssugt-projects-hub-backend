package main

import (
	"context"
	"errors"
	"net/http"
	"ssugt-projects-hub/api"
	"ssugt-projects-hub/config"
	projectrepo "ssugt-projects-hub/database/postgres/project"
	userrepo "ssugt-projects-hub/database/postgres/user"
	"ssugt-projects-hub/pkg/database/postgres"
	"ssugt-projects-hub/pkg/logging/logs"
	authservice "ssugt-projects-hub/service/auth"
	projectservice "ssugt-projects-hub/service/project"
	userservice "ssugt-projects-hub/service/user"
	"time"
)

type App struct {
	mainCtx context.Context
	log     *logs.Logs
	server  *http.Server

	userRepository    userrepo.Repository
	projectRepository projectrepo.Repository

	authService    authservice.Service
	userService    userservice.Service
	projectService projectservice.Service
}

func NewApp(mainCtx context.Context, log *logs.Logs) *App {
	return &App{
		mainCtx: mainCtx,
		log:     log,
	}
}

func (a *App) initDatabases() {
	db := postgres.NewPgx(context.Background(), config.PostgresConnection())

	a.userRepository = userrepo.NewRepository(db)
	a.projectRepository = projectrepo.NewRepository(db)

	err := runMigrations(db)
	if err != nil {
		a.log.GetLogger().Error("не смог прогнать миграции: %v", err)
		panic(err)
	}
}

func (a *App) initServices() {
	a.userService = userservice.NewService(a.log, a.userRepository)
	a.authService = authservice.NewService(a.log, a.userService)
	a.projectService = projectservice.NewProjectService(a.log, a.projectRepository)
}

func (a *App) initServer() {
	a.server = api.NewServer(a.mainCtx, a.log, a.authService, a.projectService, a.userService)
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
