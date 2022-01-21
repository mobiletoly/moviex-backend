package infra

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/appconfig"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/primary/apiserver"
	"github.com/sirupsen/logrus"
)

func RunAPIServer(port int) {
	listenAddr := fmt.Sprintf(":%d", port)

	appConfig := appconfig.Load()
	di := WireDependencies(appConfig)
	srv := apiserver.NewAPIServer(listenAddr, di)
	logrus.Info("Starting API Gateway server...")
	srv.Start()
}
