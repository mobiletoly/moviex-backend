package persist

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/sqlhelp"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/persist/internal/mapper"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/persist/internal/repo"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/outport"
)

type userAdapter struct {
	config *app.Config
	repo   repo.DBRepo
}

func NewUserAdapter(config *app.Config, dbc outport.DBConnector) outport.UserPersist {
	return &userAdapter{
		config: config,
		repo:   repo.NewDBRepo(dbc.DB()),
	}
}

func (impl userAdapter) FetchUsers(ctx context.Context, req *pb.FetchUsersRequest) (*pb.FetchUsersResponse, error) {
	page := sqlhelp.NewDBPageReq(req.GetNumRecords(), req.GetFirstRecord())
	entities, total, err := impl.repo.FetchUsers(ctx, page)
	return mapper.UserEntitiesToResponse(entities, total), err
}
