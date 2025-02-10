package auth

import (
	"context"
	"errors"
	"net/http"
	"ssugt-projects-hub/pkg/auth/accesstoken"
)

type contextKey struct {
	name string
}

var authClientCtxKey = &contextKey{"AuthClient"}

func GetContext(r *http.Request) (context.Context, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return r.Context(), errors.New("empty token")
	}

	payload, err := accesstoken.Verify(token)
	if payload == nil || err != nil {
		return r.Context(), err
	}

	return UserIdToContext(r.Context(), payload.UserId), nil
}

func IsAuthorized(ctx context.Context) bool {
	v := ctx.Value(authClientCtxKey)
	if v == nil {
		return false
	}

	_, ok := v.(int)
	if !ok {
		return false
	}

	return true
}

func UserIdFromContext(ctx context.Context) int {
	v := ctx.Value(authClientCtxKey)
	if v == nil {
		return 0
	}

	userId, ok := v.(int)
	if !ok {
		return 0
	}

	return userId
}

func UserIdToContext(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, authClientCtxKey, userId)
}
