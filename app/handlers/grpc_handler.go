package handlers

import (
	"appsku-golang/app/grpcs"
	pb "appsku-golang/files/grpc-protos"

	"google.golang.org/grpc"
)

func MainGrpcHandler(
	exampleGrpc *grpcs.ExampleGrpc,
	storeGrpc *grpcs.StoreGrpc,
) *grpc.Server {

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, exampleGrpc)
	pb.RegisterStoreServiceServer(grpcServer, storeGrpc)

	return grpcServer
}
