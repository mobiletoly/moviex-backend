package outport

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

type DBConnector interface {
	DB() *sqlx.DB
}

type UserPersist interface {
	FetchUsers(ctx context.Context, req *pb.FetchUsersRequest) (*pb.FetchUsersResponse, error)
}
