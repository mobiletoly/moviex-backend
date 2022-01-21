package persist

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/sqlhelp"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist/internal/mapper"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist/internal/repo"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/outport"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

type filmAdapter struct {
	config *app.Config
	repo   repo.DBRepo
}

func NewFilmAdapter(config *app.Config, dbc outport.DBConnector) outport.FilmPersist {
	return &filmAdapter{
		config: config,
		repo:   repo.NewDBRepo(dbc.DB()),
	}
}

func (impl filmAdapter) FetchFilms(ctx context.Context, req *pb.GetFilmsRequest) (*pb.GetFilmsResponse, error) {
	page := sqlhelp.NewDBPageReq(req.GetNumRecords(), req.GetFirstRecord())
	entities, total, err := impl.repo.FetchFilms(ctx, page)
	if err != nil {
		return nil, err
	}
	return mapper.FilmEntitiesToResponse(entities, total), err
}

func (impl filmAdapter) FetchCategory(ctx context.Context, id int32) (*pb.GetCategoryResponse, error) {
	entity, err := impl.repo.FetchCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.CategoryEntityToResponse(entity), nil
}

func (impl filmAdapter) FetchCategoryByFilmID(ctx context.Context, id int32) (*pb.GetCategoryResponse, error) {
	entity, err := impl.repo.FetchCategoryByFilmID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.CategoryEntityToResponse(entity), nil
}

func (impl filmAdapter) FetchActor(ctx context.Context, id int32) (*pb.GetActorResponse, error) {
	entity, err := impl.repo.FetchActorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ActorEntityToResponse(entity), nil
}

func (impl filmAdapter) FetchActorsByFilmID(ctx context.Context, id int32) (*pb.GetActorsResponse, error) {
	entities, err := impl.repo.FetchActorsByFilmID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ActorEntitiesToResponse(entities, int32(len(entities))), nil
}
