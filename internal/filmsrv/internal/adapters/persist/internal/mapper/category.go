package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist/internal/repo"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

func CategoryEntityToResponse(entity *repo.CategoryEntity) *pb.GetCategoryResponse {
	return &pb.GetCategoryResponse{
		Id:         entity.ID,
		Name:       entity.Name,
		LastUpdate: entity.LastUpdate.Unix(),
	}
}
