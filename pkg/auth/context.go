package auth

import (
	"context"
	"github.com/vaberof/ssugt-projects-hub-backend/pkg/domain"
)

type contextKey struct {
	name string
}

var authClientCtxKey = &contextKey{"AuthClient"}

func UserIdFromContext(ctx context.Context) *domain.UserId {
	v := ctx.Value(authClientCtxKey)
	if v == nil {
		return nil
	}

	userId, ok := v.(*domain.UserId)
	if !ok {
		return nil
	}

	return userId
}

func UserIdToContext(ctx context.Context, userId *domain.UserId) context.Context {
	return context.WithValue(ctx, authClientCtxKey, userId)
}
