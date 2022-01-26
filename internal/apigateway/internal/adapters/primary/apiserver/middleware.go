package apiserver

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/requestid"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DefaultHTTPHeadersMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		ctx := service.ContextWithAuthorization(r.Context(), authz)
		ctxLogger := logrus.StandardLogger().WithFields(logrus.Fields{
			"reqId": requestid.GetFromContext(ctx),
		})
		ctx = context.WithValue(ctx, service.LoggerContextKey, ctxLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
