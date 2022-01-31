package infra

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/appconfig"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/primary/apiserver"
	"github.com/sirupsen/logrus"
)

func RunAPIServer(deployment string) {
	cfg := appconfig.Load(deployment)
	di := WireDependencies(cfg)

	listenAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := apiserver.NewAPIServer(listenAddr, di)
	srv.ShutdownCallback = func() {
		_ = di.Close()
	}
	logrus.Info("Starting API Gateway server...")
	srv.Start()
}
