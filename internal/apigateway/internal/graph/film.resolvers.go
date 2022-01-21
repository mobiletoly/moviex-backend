package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/generated"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
)

func (r *filmResolver) Category(ctx context.Context, obj *model.Film) (*model.Category, error) {
	return r.DI.FilmUseCase.FetchCategoryByFilmID(ctx, obj.ID)
}

func (r *filmResolver) Actors(ctx context.Context, obj *model.Film) ([]*model.Actor, error) {
	return r.DI.FilmUseCase.FetchActorsByFilmID(ctx, obj.ID)
}

func (r *queryResolver) Films(ctx context.Context, first int32, after *string) (*model.FilmPage, error) {
	page := model.NewPageReqParams(after, first)
	return r.DI.FilmUseCase.FetchFilms(ctx, &page)
}

func (r *queryResolver) Category(ctx context.Context, id string) (*model.Category, error) {
	return r.DI.FilmUseCase.FetchCategory(ctx, id)
}

func (r *queryResolver) CategoryByFilmID(ctx context.Context, id string) (*model.Category, error) {
	return r.DI.FilmUseCase.FetchCategoryByFilmID(ctx, id)
}

func (r *queryResolver) Actor(ctx context.Context, id string) (*model.Actor, error) {
	return r.DI.FilmUseCase.FetchActor(ctx, id)
}

func (r *queryResolver) ActorsByFilmID(ctx context.Context, id string) ([]*model.Actor, error) {
	return r.DI.FilmUseCase.FetchActorsByFilmID(ctx, id)
}

// Film returns generated.FilmResolver implementation.
func (r *Resolver) Film() generated.FilmResolver { return &filmResolver{r} }

type filmResolver struct{ *Resolver }
