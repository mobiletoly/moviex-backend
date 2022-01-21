package apiserver

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/di"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/generated"
	"github.com/mobiletoly/moviex-backend/internal/common/requestid"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
)

func newAPIRoutes(mux *http.ServeMux, di *di.DI) {
	graphqlSrv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			// Dependency injection
			generated.Config{Resolvers: &graph.Resolver{DI: di}},
		),
	)
	graphqlSrv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		if err != nil {
			ctxLogger := logrus.StandardLogger().WithFields(logrus.Fields{
				"reqId": requestid.GetFromContext(ctx),
			})
			ctxLogger.Error(err)
		}
		return err
	})
	graphqlSrv.Use(LogrusExtension{Logger: logrus.NewEntry(logrus.StandardLogger())})
	queryMiddleware := service.DefaultHTTPHeadersMiddleware(graphqlSrv)
	queryMiddleware = requestid.HTTPHandler(queryMiddleware)

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", queryMiddleware)
	mux.HandleFunc("/version", versionHandler)
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "version 0.1")
}
