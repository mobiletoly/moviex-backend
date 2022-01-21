package infra

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/usecase"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/di"
)

func WireDependencies(appcfg *app.Config) *di.DI {
	db := persist.NewDBConnector(appcfg)
	filmAdapter := persist.NewFilmAdapter(appcfg, db)

	filmUseCase := usecase.NewFilmUseCase(appcfg, filmAdapter)

	return &di.DI{
		AppConfig:   appcfg,
		FilmUseCase: filmUseCase,
	}
}
