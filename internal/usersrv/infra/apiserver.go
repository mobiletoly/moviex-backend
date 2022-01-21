package infra

import (
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/appconfig"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/primary/apiserver"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var appPackage = "github.com/mobiletoly/moviex-backend/"

func RunAPIServer(port int) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// filename := path.Base(f.File)
			return fmt.Sprintf("%s():%d", strings.Replace(f.Function, appPackage, "", 1), f.Line), ""
		},
	})

	listenAddr := fmt.Sprintf(":%d", port)
	appConfig := appconfig.Load()
	di := WireDependencies(appConfig)
	apiserver.Serve(listenAddr, di)
}
