package infra

import (
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/persist"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/usecase"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/di"
)

func WireDependencies(appcfg *app.Config) *di.DI {
	db := persist.NewDBConnector(appcfg)
	userAdapter := persist.NewUserAdapter(appcfg, db)

	userUseCase := usecase.NewUserUseCase(appcfg, userAdapter)

	return &di.DI{
		AppConfig:   appcfg,
		UserUseCase: userUseCase,
	}
}
