package service

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LoggerContextKeyType struct{}

var (
	LoggerContextKey = LoggerContextKeyType(struct{}{})
)

func LogEntry(ctx context.Context) *logrus.Entry {
	l, _ := ctx.Value(LoggerContextKey).(*logrus.Entry)
	return l
}
