package usecase

import (
	"context"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/outport"
)

type User struct {
	appcfg      *app.Config
	userPersist outport.UserPersist
}

func NewUserUseCase(appcfg *app.Config, userPersist outport.UserPersist) *User {
	return &User{
		appcfg:      appcfg,
		userPersist: userPersist,
	}
}

func (impl *User) FetchUsers(ctx context.Context, req *pb.FetchUsersRequest) (*pb.FetchUsersResponse, error) {
	return impl.userPersist.FetchUsers(ctx, req)
}
