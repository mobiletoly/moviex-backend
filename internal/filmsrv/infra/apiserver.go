package infra

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/appconfig"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/primary/apiserver"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var appPackage = "github.com/mobiletoly/moviex-backend/"

func RunAPIServer(deployment string) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// filename := path.Base(f.File)
			return fmt.Sprintf("%s():%d", strings.Replace(f.Function, appPackage, "", 1), f.Line), ""
		},
	})

	cfg := appconfig.Load(deployment)
	di := WireDependencies(cfg)
	listenAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	apiserver.Serve(listenAddr, di)
}
