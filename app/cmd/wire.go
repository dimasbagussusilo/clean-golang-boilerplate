//go:build wireinject
// +build wireinject

package main

import (
	"appsku-golang/app/controllers"
	"appsku-golang/app/grpcs"
	"appsku-golang/app/handlers"
	"appsku-golang/app/repositories"
	"appsku-golang/app/routes"
	"appsku-golang/app/services"
	"appsku-golang/app/usecases"

	"appsku-golang/app/global-utils/grpcclient"
	kafkadbo "appsku-golang/app/global-utils/kafka"
	"appsku-golang/app/global-utils/mongodb"
	"appsku-golang/app/global-utils/redisdb"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

var pkgSet = wire.NewSet(
	grpcclient.NewWrapGrpcClient,
	redisdb.NewRedisConn,
	kafkadbo.NewKafkaPublisher,
	mongodb.NewMongoDB,
)

var pkgConsumerSet = wire.NewSet(
	grpcclient.NewWrapGrpcClient,
	redisdb.NewRedisConn,
	kafkadbo.NewKafkaPublisher,
	mongodb.NewMongoDB,
)

var exampleSet = wire.NewSet(
	controllers.NewExampleController,
	usecases.NewExampleUseCase,
	usecases.NewExampleValidator,
	repositories.NewExampleRepository,
)

var storeSet = wire.NewSet(
	controllers.NewStoreController,
	usecases.NewStoreUseCase,
	//usecases.NewStoreValidator,
	repositories.NewStoreRepository,
)

var grpcServiceSet = wire.NewSet(
	services.NewExampleGrpcClient,
)

var grpcServerSet = wire.NewSet(
	grpcs.NewExampleGrpc,
	grpcs.NewStoreGrpc,
)

var consumerSet = wire.NewSet(
	controllers.NewConsumerController,
	usecases.NewConsumerUseCase,
)

func InitializeHttpServer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) *gin.Engine {
	wire.Build(
		pkgSet,
		exampleSet,
		storeSet,
		//grpcServiceSet,
		routes.NewHttpRoute,
	)
	return nil
}

func InitializeGrpcServer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) *grpc.Server {
	wire.Build(
		pkgSet,
		exampleSet,
		storeSet,
		//grpcServiceSet,
		grpcServerSet,
		handlers.MainGrpcHandler,
	)
	return nil
}

func InitializeConsumer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) handlers.IConsumerHandler {
	wire.Build(
		pkgConsumerSet,
		consumerSet,
		//grpcServiceSet,
		handlers.MainConsumerHandler,
	)
	return nil
}
