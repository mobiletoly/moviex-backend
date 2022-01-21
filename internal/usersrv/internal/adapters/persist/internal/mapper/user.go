package mapper

import (
	rpcmoviex "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/adapters/persist/internal/repo"
)

func UserEntityToResponse(entity *repo.UserEntity) *rpcmoviex.FetchUserResponse {
	return &rpcmoviex.FetchUserResponse{
		Id:       entity.ID,
		Email:    entity.Email,
		Password: entity.Password,
	}
}

func UserEntitiesToResponse(entities []*repo.UserEntity, totalRecords int32) *rpcmoviex.FetchUsersResponse {
	resps := make([]*rpcmoviex.FetchUserResponse, len(entities))
	for i := range entities {
		resps[i] = UserEntityToResponse(entities[i])
	}
	return &rpcmoviex.FetchUsersResponse{
		Users:        resps,
		TotalRecords: totalRecords,
	}
}
