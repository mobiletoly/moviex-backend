package remoting

import (
	"context"
	"fmt"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/adapters/remoting/mapper"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/core/outport"
	"github.com/mobiletoly/moviex-backend/internal/apigateway/internal/graph/model"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type userAdapter struct {
	config     *app.Config
	userClient pb.UserClient
}

func NewUserAdapter(config *app.Config) (outport.UserRemoting, error) {
	addr := fmt.Sprintf("%s:%d", config.UserSrv.Host, config.UserSrv.Port)
	logrus.Infoln("create User service client:", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("error creating User client connection: %w", err)
	}
	userClient := pb.NewUserClient(conn)
	return &userAdapter{
		config:     config,
		userClient: userClient,
	}, nil
}

func (impl userAdapter) FetchUsers(ctx context.Context, page *model.PageReqParams) (*model.UserPage, error) {
	req := pb.FetchUsersRequest{
		NumRecords:  page.Count(),
		FirstRecord: page.AfterOffset(),
	}
	callCtx := newCallOutgoingContext(ctx)
	resp, err := impl.userClient.FetchUsers(callCtx, &req)
	if err != nil {
		return nil, fmt.Errorf("error fetching user page: %w", err)
	}
	models := mapper.FetchUserResponsesToModels(resp.GetUsers())
	pageInfo := model.NewPageInfoWithRecordIndexAndTotal(page.AfterOffset(), int32(len(models)), resp.GetTotalRecords())
	return &model.UserPage{
		TotalCount: resp.TotalRecords,
		PageInfo:   pageInfo,
		Edges:      mapper.WrapUsers(models),
	}, nil
}
