package outport

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
)

type FilmRemoting interface {
	FetchFilms(ctx context.Context, page *model.PageReqParams) (*model.FilmPage, error)
	FetchCategory(ctx context.Context, id int32) (*model.Category, error)
	FetchCategoryByFilmID(ctx context.Context, id int32) (*model.Category, error)
	FetchActor(ctx context.Context, id int32) (*model.Actor, error)
	FetchActorsByFilmID(ctx context.Context, id int32) ([]*model.Actor, error)
}

type UserRemoting interface {
	FetchUsers(ctx context.Context, page *model.PageReqParams) (*model.UserPage, error)
}
