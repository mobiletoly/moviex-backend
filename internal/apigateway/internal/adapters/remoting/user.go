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
	"google.golang.org/grpc/credentials/insecure"
)

type userAdapter struct {
	config       *app.Config
	userClient   pb.UserClient
	closeHandler func() error
}

func NewUserAdapter(config *app.Config) (outport.UserRemoting, error) {
	addr := fmt.Sprintf("%s:%d", config.UserSrv.Host, config.UserSrv.Port)
	logrus.Infoln("create User service client:", addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error creating User client connection: %w", err)
	}
	userClient := pb.NewUserClient(conn)
	ua := &userAdapter{
		config:     config,
		userClient: userClient,
		closeHandler: func() error {
			logrus.Infoln("close User service")
			return conn.Close()
		},
	}
	return ua, nil
}

func (impl userAdapter) Close() error {
	return impl.closeHandler()
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
