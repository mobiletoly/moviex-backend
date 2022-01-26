package apiserver

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/requestid"
	"github.com/mobiletoly/moviex-backend/internal/common/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sirupsen/logrus"
)

type LogrusExtension struct {
	Logger *logrus.Entry
}

var _ interface {
	graphql.HandlerExtension
	graphql.FieldInterceptor
} = LogrusExtension{}

func (n LogrusExtension) ExtensionName() string {
	return "LogrusExtension"
}

func (n LogrusExtension) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (n LogrusExtension) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	oc := graphql.GetOperationContext(ctx)
	fc := graphql.GetFieldContext(ctx)
	ctxLogger := n.Logger.WithFields(logrus.Fields{
		"reqId":     requestid.GetFromContext(ctx),
		"operation": oc.OperationName,
		"field":     fc.Field.Name,
	})
	ctx = context.WithValue(ctx, service.LoggerContextKey, ctxLogger)
	return next(ctx)
}
