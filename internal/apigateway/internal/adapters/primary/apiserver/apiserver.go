package apiserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/di"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

var appPackage = "github.com/mobiletoly/moviex-backend/"

const gracefulShutdownTimeout = 10 * time.Second

// Server provides an http.Server
type Server struct {
	*http.Server
	ShutdownCallback func()
}

func NewAPIServer(listenAddr string, di *di.DI) *Server {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	// this formatter prints logrus logs formatted as:
	//    [directory-relative-to-project].[function]:[line-number] message
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s():%d", strings.Replace(f.Function, appPackage, "", 1), f.Line), ""
		},
	})

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"*"}, // by some reasons this library is case-sensitive for headers
		AllowCredentials: true,
		Debug:            true, // TODO Enable Debugging for testing, consider disabling in production
	})

	newAPIRoutes(mux, di)
	hdl := c.Handler(mux)

	srv := http.Server{
		Addr:         listenAddr,
		Handler:      hdl,
		ErrorLog:     log.New(logrus.StandardLogger().Writer(), "", 0),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Server: &srv,
	}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) Start() {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal("Could not listen on address (", srv.Addr, ") / Error: ", err.Error())
		}
	}()
	logrus.Info("Server is ready to handle requests on address ", srv.Addr)
	srv.gracefulShutdown()
}

// Graceful shutdown. Normally if application receives a signal to be stopped, it will do just that - quit.
// We give some HTTP server some time to close open connections (and may be finish serving HTTP calls).
func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logrus.Info("Server is shutting down / Reasons: ", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Could not gracefully shutdown the server / Reason: ", err.Error())
	}
	logrus.Info("Server stopped")

	(srv.ShutdownCallback)()
}
