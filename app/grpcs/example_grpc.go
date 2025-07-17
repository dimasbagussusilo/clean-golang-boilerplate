package grpcs

import (
	"appsku-golang/app/usecases"
	pb "appsku-golang/files/grpc-protos"
)

type ExampleGrpc struct {
	pb.ExampleServiceServer
	ExampleUseCase usecases.IExampleUseCase
}

func NewExampleGrpc(ExampleUseCase usecases.IExampleUseCase) *ExampleGrpc {
	return &ExampleGrpc{
		ExampleUseCase: ExampleUseCase,
	}
}
