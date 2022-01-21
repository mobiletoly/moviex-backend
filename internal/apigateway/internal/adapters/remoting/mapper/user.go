package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

func FetchUserResponseToModel(r *pb.FetchUserResponse) *model.User {
	return &model.User{
		ID:       r.Id,
		Email:    r.Email,
		Password: r.Password,
	}
}

func FetchUserResponsesToModels(rs []*pb.FetchUserResponse) []*model.User {
	models := make([]*model.User, len(rs))
	for i := range rs {
		models[i] = FetchUserResponseToModel(rs[i])
	}
	return models
}

func WrapUsers(users []*model.User) []*model.UserEdge {
	edges := make([]*model.UserEdge, len(users))
	for i, user := range users {
		edges[i] = &model.UserEdge{Node: user}
	}
	return edges
}
