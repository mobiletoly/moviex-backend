package service

import (
	"context"
	"errors"
)

type reqContextKey string

var (
	authzKey = reqContextKey("Authorization")
)

func ContextWithAuthorization(ctx context.Context, authz string) context.Context {
	return context.WithValue(ctx, authzKey, authz)
}

func GetAuthorization(ctx context.Context) (string, error) {
	if auth := ctx.Value(authzKey).(string); auth != "" {
		return auth, nil
	}
	return "", errors.New("authorization header is required")
}
