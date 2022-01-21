package infra

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/remoting"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/usecase"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/di"
	"github.com/sirupsen/logrus"
)

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

	return &di.DI{
		AppConfig:   cfg,
		FilmUseCase: filmUseCase,
		UserUseCase: userUseCase,
	}
}
