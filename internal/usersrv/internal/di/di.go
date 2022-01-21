package di

import (
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/usecase"
)

type DI struct {
	AppConfig   *app.Config
	UserUseCase *usecase.User
}
