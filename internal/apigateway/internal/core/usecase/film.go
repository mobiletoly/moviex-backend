package usecase

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/model"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/outport"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
)

type Film struct {
	appcfg       *app.Config
	filmRemoting outport.FilmRemoting
}

func NewFilmUseCase(appcfg *app.Config, fr outport.FilmRemoting) *Film {
	return &Film{
		appcfg:       appcfg,
		filmRemoting: fr,
	}
}

func (impl *Film) FetchFilms(ctx context.Context, page *model.PageReqParams) (*model.FilmPage, error) {
	service.LogEntry(ctx).Debug("fetch films")
	return impl.filmRemoting.FetchFilms(ctx, page)
}

func (impl *Film) FetchCategory(ctx context.Context, id string) (*model.Category, error) {
	service.LogEntry(ctx).Debugf("fetch category by id=%s", id)
	idval, err := idToInt32(id)
	if err != nil {
		return nil, err
	}
	return impl.filmRemoting.FetchCategory(ctx, idval)
}

func (impl *Film) FetchCategoryByFilmID(ctx context.Context, id string) (*model.Category, error) {
	service.LogEntry(ctx).Debugf("fetch category by film id=%s", id)
	idval, err := idToInt32(id)
	if err != nil {
		return nil, err
	}
	return impl.filmRemoting.FetchCategoryByFilmID(ctx, idval)
}

func (impl *Film) FetchActor(ctx context.Context, id string) (*model.Actor, error) {
	service.LogEntry(ctx).Debugf("fetch actor by id=%s", id)
	idval, err := idToInt32(id)
	if err != nil {
		return nil, err
	}
	return impl.filmRemoting.FetchActor(ctx, idval)
}

func (impl *Film) FetchActorsByFilmID(ctx context.Context, id string) ([]*model.Actor, error) {
	service.LogEntry(ctx).Debugf("fetch actors by film id=%s", id)
	idval, err := idToInt32(id)
	if err != nil {
		return nil, err
	}
	return impl.filmRemoting.FetchActorsByFilmID(ctx, idval)
}
