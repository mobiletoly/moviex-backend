package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/filmsrv/internal/adapters/persist/internal/repo"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

func FilmEntityToResponse(entity *repo.FilmEntity) *pb.GetFilmResponse {
	return &pb.GetFilmResponse{
		Id:              entity.ID,
		Title:           entity.Title,
		Description:     entity.Description,
		ReleaseYear:     uint32(entity.ReleaseYear),
		LanguageId:      entity.LanguageID,
		Length:          uint32(entity.Length),
		Rating:          entity.Rating,
		UpdateTime:      entity.LastUpdate.Unix(),
		SpecialFeatures: entity.SpecialFeatures,
	}
}

func FilmEntitiesToResponse(entities []*repo.FilmEntity, totalRecords int32) *pb.GetFilmsResponse {
	resps := make([]*pb.GetFilmResponse, len(entities))
	for i := range entities {
		resps[i] = FilmEntityToResponse(entities[i])
	}
	return &pb.GetFilmsResponse{
		Films:        resps,
		TotalRecords: totalRecords,
	}
}
