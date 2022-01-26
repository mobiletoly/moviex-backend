package remoting

import (
	"context"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/remoting/mapper"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/model"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/outport"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type filmAdapter struct {
	config     *app.Config
	filmClient pb.FilmClient
}

func NewFilmAdapter(config *app.Config) (outport.FilmRemoting, error) {
	addr := fmt.Sprintf("%s:%d", config.FilmSrv.Host, config.FilmSrv.Port)
	logrus.Infoln("create Film service client:", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed creating Film client connection: %w", err)
	}
	filmClient := pb.NewFilmClient(conn)
	return &filmAdapter{
		config:     config,
		filmClient: filmClient,
	}, nil
}

func (impl filmAdapter) FetchFilms(ctx context.Context, page *model.PageReqParams) (*model.FilmPage, error) {
	req := pb.GetFilmsRequest{
		NumRecords:  page.Count(),
		FirstRecord: page.AfterOffset(),
	}
	callCtx := newCallOutgoingContext(ctx)
	resp, err := impl.filmClient.GetFilms(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching film page: %w", err)
	}
	ms := mapper.FetchFilmResponsesToModels(resp.GetFilms())
	pageInfo := model.NewPageInfoWithRecordIndexAndTotal(page.AfterOffset(), int32(len(ms)), resp.GetTotalRecords())
	return &model.FilmPage{
		TotalCount: resp.TotalRecords,
		PageInfo:   pageInfo,
		Edges:      mapper.WrapFilms(ms),
	}, nil
}

func (impl filmAdapter) FetchCategory(ctx context.Context, id int32) (*model.Category, error) {
	callCtx := newCallOutgoingContext(ctx)
	req := pb.GetCategoryRequest{Id: id}
	resp, err := impl.filmClient.GetCategory(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching category: %w", err)
	}
	m := mapper.FetchCategoryResponseToModel(resp)
	return m, nil
}

func (impl filmAdapter) FetchCategoryByFilmID(ctx context.Context, id int32) (*model.Category, error) {
	callCtx := newCallOutgoingContext(ctx)
	req := pb.GetByFilmIdRequest{Id: id}
	resp, err := impl.filmClient.GetCategoryByFilmID(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching category by film id: %w", err)
	}
	m := mapper.FetchCategoryResponseToModel(resp)
	return m, nil
}

func (impl filmAdapter) FetchActor(ctx context.Context, id int32) (*model.Actor, error) {
	callCtx := newCallOutgoingContext(ctx)
	req := pb.GetActorRequest{Id: id}
	resp, err := impl.filmClient.GetActor(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching actor: %w", err)
	}
	m := mapper.FetchActorResponseToModel(resp)
	return m, nil
}

func (impl filmAdapter) FetchActorsByFilmID(ctx context.Context, id int32) ([]*model.Actor, error) {
	callCtx := newCallOutgoingContext(ctx)
	req := pb.GetByFilmIdRequest{Id: id}
	resp, err := impl.filmClient.GetActorsByFilmID(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed fetching actors by film id=%d: %w", id, err)
	}
	m := mapper.FetchActorsResponseToModels(resp)
	return m, nil
}
