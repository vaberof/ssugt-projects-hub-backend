package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	service "github.com/vaberof/ssugt-projects-hub-backend/internal/domain/auth"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/logging/logs"
	"net/http"
)

func AuthMiddleware(authService service.AuthService, logs *logs.Logs) gin.HandlerFunc {
	loggerName := "auth-middleware"
	logger := logs.WithName(loggerName)

	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		fmt.Printf("accessToken: %s\n", accessToken)
		if accessToken == "" {
			logger.Error("Client not authenticated: empty access-token")

			c.JSON(http.StatusUnauthorized, "Need to provide access access-token")
			c.Abort()
			return
		}

		userId, err := authService.VerifyAccessToken(accessToken)
		if err != nil {
			logger.Error(fmt.Sprintf("Invalid access-token: %v", err))

			c.JSON(http.StatusUnauthorized, fmt.Sprintf("invalid access-token: %s", err.Error()))
			c.Abort()
			return
		}

		logger.Info("Client is authenticated")

		c.Set("userId", userId)

		c.Next()
	}
}
