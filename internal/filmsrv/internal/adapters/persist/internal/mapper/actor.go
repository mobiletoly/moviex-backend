package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist/internal/repo"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

func ActorEntityToResponse(entity *repo.ActorEntity) *pb.GetActorResponse {
	return &pb.GetActorResponse{
		Id:         entity.ID,
		FirstName:  entity.FirstName,
		LastName:   entity.LastName,
		LastUpdate: entity.LastUpdate.Unix(),
	}
}

func ActorEntitiesToResponse(entities []*repo.ActorEntity, totalRecords int32) *pb.GetActorsResponse {
	resps := make([]*pb.GetActorResponse, len(entities))
	for i := range entities {
		resps[i] = ActorEntityToResponse(entities[i])
	}
	return &pb.GetActorsResponse{
		Actors:       resps,
		TotalRecords: totalRecords,
	}
}
