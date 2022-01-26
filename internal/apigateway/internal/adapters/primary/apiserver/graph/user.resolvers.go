package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/model"
)

func (r *queryResolver) Users(ctx context.Context, first int32, after *string) (*model.UserPage, error) {
	page := model.NewPageReqParams(after, first)
	return r.DI.UserUseCase.FetchUsers(ctx, &page)
}
