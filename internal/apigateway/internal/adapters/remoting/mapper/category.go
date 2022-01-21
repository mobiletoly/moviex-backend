package mapper

import (
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"strconv"
)

func FetchCategoryResponseToModel(r *pb.GetCategoryResponse) *model.Category {
	return &model.Category{
		ID:         strconv.Itoa(int(r.Id)),
		Name:       r.Name,
		LastUpdate: r.LastUpdate,
	}
}
