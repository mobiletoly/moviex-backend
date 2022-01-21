package di

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/usecase"
)

type DI struct {
	AppConfig   *app.Config
	FilmUseCase *usecase.Film
}
