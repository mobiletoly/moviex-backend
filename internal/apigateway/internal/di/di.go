package di

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/usecase"
)

type DI struct {
	AppConfig   *app.Config
	FilmUseCase *usecase.Film
	UserUseCase *usecase.User
}
