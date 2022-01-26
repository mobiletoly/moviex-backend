package requestid

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

var ridKey = struct{}{}

const xRequestIDHeaderName = "X-Request-ID"

// GetFromContext returns the request id from context
func GetFromContext(ctx context.Context) string {
	rid := ctx.Value(ridKey).(string)
	return rid
}

// HTTPHandler sets unique request id.
// If header `X-Request-ID` is already present in the request, that is considered the
// request id. Otherwise, generates a new unique ID.
func HTTPHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(xRequestIDHeaderName)
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(xRequestIDHeaderName, rid)
		}
		ctx := newRequestIDContext(r.Context(), rid)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// newRequestIDContext creates a context with request id
func newRequestIDContext(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ridKey, rid)
}
