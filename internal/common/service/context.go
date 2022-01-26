package service

import (
	"context"
	"errors"
)

var reqContextKey = struct{}{}

func ContextWithAuthorization(ctx context.Context, authz string) context.Context {
	return context.WithValue(ctx, reqContextKey, authz)
}

func GetAuthorization(ctx context.Context) (string, error) {
	if auth := ctx.Value(reqContextKey).(string); auth != "" {
		return auth, nil
	}
	return "", errors.New("authorization header is required")
}
