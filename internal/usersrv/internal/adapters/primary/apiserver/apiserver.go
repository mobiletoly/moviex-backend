package apiserver

import (
	"context"
	"github.com/google/uuid"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/di"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

type server struct {
	pb.UnsafeUserServer
	di *di.DI
}

func Serve(listenAddr string, di *di.DI) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(defaultInterceptor))
	pb.RegisterUserServer(s, &server{di: di})
	logrus.Info("Listening gRPC server: ", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("failed to serve:", err)
	}
}

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

func getRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v, ok := md["x-request-id"]; ok {
			return v[0]
		}
	}
	return uuid.New().String()
}
