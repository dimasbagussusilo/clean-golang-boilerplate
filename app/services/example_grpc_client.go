package services

import (
	dboconst "appsku-golang/app/constants"
	pb "appsku-golang/files/grpc-protos"
	"context"

	"appsku-golang/app/global-utils/grpcclient"
	"appsku-golang/app/global-utils/model"
)

type IExampleGrpcClient interface {
	GetExampleByID(ctx context.Context, id int64) *model.Response
}

type ExampleGrpcClient struct {
	ExampleService pb.ExampleServiceClient
}

func NewExampleGrpcClient(connection grpcclient.IGRPCClients) IExampleGrpcClient {
	conn := connection.Client(dboconst.ExampleService)

	return &ExampleGrpcClient{
		ExampleService: pb.NewExampleServiceClient(conn),
	}
}

func (s *ExampleGrpcClient) GetExampleByID(ctx context.Context, id int64) *model.Response {
	response := &model.Response{}

	return response
}
