package remoting

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/requestid"
	"google.golang.org/grpc/metadata"
)

func newCallOutgoingContext(ctx context.Context) context.Context {
	rid := requestid.GetFromContext(ctx)
	md := metadata.Pairs("x-request-id", rid)
	callCtx := metadata.NewOutgoingContext(ctx, md)
	return callCtx
}
