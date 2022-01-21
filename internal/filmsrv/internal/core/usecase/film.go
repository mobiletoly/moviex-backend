package usecase

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/core/outport"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

type Film struct {
	pb.UnsafeFilmServer
	appcfg      *app.Config
	filmPersist outport.FilmPersist
}

func NewFilmUseCase(appcfg *app.Config, filmPersist outport.FilmPersist) *Film {
	return &Film{
		appcfg:      appcfg,
		filmPersist: filmPersist,
	}
}

func (impl *Film) GetFilms(ctx context.Context, req *pb.GetFilmsRequest) (*pb.GetFilmsResponse, error) {
	service.LogEntry(ctx).Debugf("fetch films: numRecords=%d firstRecord=%d", req.NumRecords, req.FirstRecord)
	return impl.filmPersist.FetchFilms(ctx, req)
}

func (impl *Film) GetCategory(ctx context.Context, request *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	service.LogEntry(ctx).Debugf("fetch category by id=%d", request.GetId())
	return impl.filmPersist.FetchCategory(ctx, request.GetId())
}

func (impl *Film) GetCategoryByFilmID(ctx context.Context, request *pb.GetByFilmIdRequest) (*pb.GetCategoryResponse, error) {
	service.LogEntry(ctx).Debugf("fetch category by film id=%d", request.GetId())
	return impl.filmPersist.FetchCategoryByFilmID(ctx, request.GetId())
}

func (impl *Film) GetActor(ctx context.Context, request *pb.GetActorRequest) (*pb.GetActorResponse, error) {
	service.LogEntry(ctx).Debugf("fetch actor by id=%d", request.GetId())
	return impl.filmPersist.FetchActor(ctx, request.GetId())
}

func (impl *Film) GetActorsByFilmID(ctx context.Context, request *pb.GetByFilmIdRequest) (*pb.GetActorsResponse, error) {
	service.LogEntry(ctx).Debugf("fetch actors by film id=%d", request.GetId())
	return impl.filmPersist.FetchActorsByFilmID(ctx, request.GetId())
}
