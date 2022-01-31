package infra

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/remoting"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/usecase"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/di"
	"github.com/sirupsen/logrus"
)

// WireDependencies provides entry point to wire dependencies together for further dependency injection
func WireDependencies(cfg *app.Config) *di.DI {
	filmAdapter, err := remoting.NewFilmAdapter(cfg)
	if err != nil {
		logrus.Fatalf("error wiring Film adapter dependency: %v", err)
	}
	userAdapter, err := remoting.NewUserAdapter(cfg)
	if err != nil {
		logrus.Fatalf("error wiring User adapter dependency: %v", err)
	}

	filmUseCase := usecase.NewFilmUseCase(cfg, filmAdapter)
	userUseCase := usecase.NewUserUseCase(cfg, userAdapter)

	cleanupCallback := func() error {
		_ = filmAdapter.Close()
		_ = userAdapter.Close()
		return nil
	}

	return &di.DI{
		AppConfig:   cfg,
		FilmUseCase: filmUseCase,
		UserUseCase: userUseCase,
		Close:       cleanupCallback,
	}
}
