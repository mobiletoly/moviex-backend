package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"strconv"
)

func FetchActorResponseToModel(r *pb.GetActorResponse) *model.Actor {
	return &model.Actor{
		ID:         strconv.Itoa(int(r.Id)),
		FirstName:  r.GetFirstName(),
		LastName:   r.GetLastName(),
		LastUpdate: r.GetLastUpdate(),
	}
}

func FetchActorsResponseToModels(rs *pb.GetActorsResponse) []*model.Actor {
	models := make([]*model.Actor, len(rs.Actors))
	for i := range rs.Actors {
		models[i] = FetchActorResponseToModel(rs.Actors[i])
	}
	return models
}
