package httpserver

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/http/httpserver/middleware/logging"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"log/slog"
	"net/http"
)

const MaxMultipartMemory = 100 << 20 // 100 MB

const defaultMultipartMemory = 100 << 20 // 100 MB

type AppServer struct {
	Server *http.Server
	config *ServerConfig
	logger *slog.Logger
}

func New(config *ServerConfig, logs *logs.Logs) *AppServer {
	loggingMw := logging.New(logs)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.MaxMultipartMemory = defaultMultipartMemory
	router.Use(loggingMw.Handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}

	return &AppServer{
		Server: server,
		config: config,
		logger: loggingMw.Logger,
	}
}

func (server *AppServer) StartAsync() <-chan error {
	server.logger.Info("Starting http server")

	exitChannel := make(chan error)

	go func() {
		err := server.Server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("Failed to start HTTP server", slog.Group("error", err))
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Info("Started HTTP server", slog.Group("http-server", "address", server.Server.Addr))

	return exitChannel
}

func (server *AppServer) GetLogger() *slog.Logger {
	return server.logger
}

func (server *AppServer) LoadGinEngineFromHTTPHandler() *gin.Engine {
	v, ok := server.Server.Handler.(*gin.Engine)
	if !ok {
		return nil
	}
	return v
}
