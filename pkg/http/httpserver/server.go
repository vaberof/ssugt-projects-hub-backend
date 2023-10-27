package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/httpserver/middleware/logging"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type AppServer struct {
	Server  *gin.Engine
	config  *ServerConfig
	logger  *slog.Logger
	address string
}

func New(config *ServerConfig, logs *logs.Logs) *AppServer {
	loggingMw := logging.New(logs)

	ginServer := gin.New()
	ginServer.Use(loggingMw.Handler)

	return &AppServer{
		Server:  ginServer,
		config:  config,
		logger:  loggingMw.Logger,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

func (server *AppServer) StartAsync() *chan error {
	exitChannel := make(chan error)

	server.logger.Info("Starting http server")

	go func() {
		err := http.ListenAndServe(server.address, server.Server)
		if err != nil {
			server.logger.Error("Failed to start HTTP server", slog.Group("error", err))
			exitChannel <- err
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Info("Started HTTP server", slog.Group("http-server", "address", server.address))

	return &exitChannel
}

func (server *AppServer) GetLogger() *slog.Logger {
	return server.logger
}
