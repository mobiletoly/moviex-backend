package outport

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

type DBConnector interface {
	DB() *sqlx.DB
}

type FilmPersist interface {
	FetchFilms(ctx context.Context, req *pb.GetFilmsRequest) (*pb.GetFilmsResponse, error)
	FetchCategory(ctx context.Context, id int32) (*pb.GetCategoryResponse, error)
	FetchCategoryByFilmID(ctx context.Context, id int32) (*pb.GetCategoryResponse, error)
	FetchActor(ctx context.Context, id int32) (*pb.GetActorResponse, error)
	FetchActorsByFilmID(ctx context.Context, id int32) (*pb.GetActorsResponse, error)
}
