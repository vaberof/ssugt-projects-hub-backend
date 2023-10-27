package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

func (handler *Handler) userIdFromContext(ctx *gin.Context) (*domain.UserId, error) {
	v, ok := ctx.Get("userId")
	if !ok {
		return nil, errors.New("failed to get userId from context")
	}

	userId := v.(*domain.UserId)
	if userId == nil {
		return nil, errors.New("failed to get domain userId from context value")
	}

	return userId, nil
}
