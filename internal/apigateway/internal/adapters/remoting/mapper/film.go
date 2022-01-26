package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/model"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"strconv"
)

func FetchFilmResponseToModel(r *pb.GetFilmResponse) *model.Film {
	return &model.Film{
		ID:    strconv.Itoa(int(r.Id)),
		Title: r.Title,
	}
}

func FetchFilmResponsesToModels(rs []*pb.GetFilmResponse) []*model.Film {
	models := make([]*model.Film, len(rs))
	for i := range rs {
		models[i] = FetchFilmResponseToModel(rs[i])
	}
	return models
}

func WrapFilms(films []*model.Film) []*model.FilmEdge {
	edges := make([]*model.FilmEdge, len(films))
	for i, film := range films {
		edges[i] = &model.FilmEdge{Node: film}
	}
	return edges
}
