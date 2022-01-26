package usecase

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/model"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/outport"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
)

type User struct {
	appcfg       *app.Config
	userRemoting outport.UserRemoting
}

func NewUserUseCase(appcfg *app.Config, fr outport.UserRemoting) *User {
	return &User{
		appcfg:       appcfg,
		userRemoting: fr,
	}
}

func (impl *User) FetchUsers(ctx context.Context, page *model.PageReqParams) (*model.UserPage, error) {
	service.LogEntry(ctx).Debug("fetch users")
	return impl.userRemoting.FetchUsers(ctx, page)
}
