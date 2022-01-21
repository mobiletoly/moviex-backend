package apiserver

import (
	"context"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	pb "github.com/mobiletoly/moviex-backend/internal/proto"
)

func (s server) FetchUsers(ctx context.Context, req *pb.FetchUsersRequest) (*pb.FetchUsersResponse, error) {
	service.LogEntry(ctx).Infof("Fetch users: firstRecord=%d numRecords=%d", req.GetFirstRecord(), req.GetNumRecords())
	return s.di.UserUseCase.FetchUsers(ctx, req)
}
