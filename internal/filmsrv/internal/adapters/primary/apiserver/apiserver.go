package apiserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/di"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

// Serve launches gRPC service to handle Film functionality
func Serve(listenAddr string, di *di.DI) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(defaultInterceptor))
	pb.RegisterFilmServer(s, di.FilmUseCase)
	logrus.Info("Listening gRPC server: ", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("failed to serve:", err)
	}
}

// defaultInterceptor allows us to add request-specific values into request context
func defaultInterceptor(
	ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	ctxLogger := logrus.StandardLogger().WithFields(logrus.Fields{
		"reqId": getRequestID(ctx),
	})
	ctx = context.WithValue(ctx, service.LoggerContextKey, ctxLogger)
	m, err := handler(ctx, req)
	if err != nil {
		logrus.Errorf("gRPC failed with error %v", err)
	}
	return m, err
}

// getRequestID returns value of x-request-id header (we use this value to inject request id into context)
func getRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v, ok := md["x-request-id"]; ok {
			return v[0]
		}
	}
	return uuid.New().String()
}
