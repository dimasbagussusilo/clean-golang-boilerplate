// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"appsku-golang/app/controllers"
	"appsku-golang/app/global-utils/grpcclient"
	"appsku-golang/app/global-utils/kafka"
	"appsku-golang/app/global-utils/mongodb"
	"appsku-golang/app/global-utils/redisdb"
	"appsku-golang/app/grpcs"
	"appsku-golang/app/handlers"
	"appsku-golang/app/repositories"
	"appsku-golang/app/routes"
	"appsku-golang/app/services"
	"appsku-golang/app/usecases"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// Injectors from wire.go:

func InitializeHttpServer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) *gin.Engine {
	iMongoDB := mongodb.NewMongoDB(mongoDBParam)
	iRedis := redisdb.NewRedisConn(redisParam)
	iExampleRepository := repositories.NewExampleRepository(iMongoDB, iRedis)
	iExampleUseCase := usecases.NewExampleUseCase(iExampleRepository)
	exampleController := controllers.NewExampleController(iExampleUseCase)
	iStoreRepository := repositories.NewStoreRepository(iMongoDB, iRedis)
	iStoreUseCase := usecases.NewStoreUseCase(iStoreRepository)
	storeController := controllers.NewStoreController(iStoreUseCase)
	engine := routes.NewHttpRoute(exampleController, storeController)
	return engine
}

func InitializeFiberServer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) *fiber.App {
	iMongoDB := mongodb.NewMongoDB(mongoDBParam)
	iRedis := redisdb.NewRedisConn(redisParam)
	iExampleRepository := repositories.NewExampleRepository(iMongoDB, iRedis)
	iExampleUseCase := usecases.NewExampleUseCase(iExampleRepository)
	exampleController := controllers.NewExampleController(iExampleUseCase)
	iStoreRepository := repositories.NewStoreRepository(iMongoDB, iRedis)
	iStoreUseCase := usecases.NewStoreUseCase(iStoreRepository)
	storeController := controllers.NewStoreController(iStoreUseCase)
	app := routes.NewFiberRoute(exampleController, storeController)
	return app
}

func InitializeGrpcServer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) *grpc.Server {
	iMongoDB := mongodb.NewMongoDB(mongoDBParam)
	iRedis := redisdb.NewRedisConn(redisParam)
	iExampleRepository := repositories.NewExampleRepository(iMongoDB, iRedis)
	iExampleUseCase := usecases.NewExampleUseCase(iExampleRepository)
	exampleGrpc := grpcs.NewExampleGrpc(iExampleUseCase)
	iStoreRepository := repositories.NewStoreRepository(iMongoDB, iRedis)
	iStoreUseCase := usecases.NewStoreUseCase(iStoreRepository)
	storeGrpc := grpcs.NewStoreGrpc(iStoreUseCase)
	server := handlers.MainGrpcHandler(exampleGrpc, storeGrpc)
	return server
}

func InitializeConsumer(grpcParams []grpcclient.GRPCClientParam, redisParam redisdb.RedisParam, kafkaParam []string, mongoDBParam mongodb.MongoDBParam) handlers.IConsumerHandler {
	iConsumerUseCase := usecases.NewConsumerUseCase()
	consumerController := controllers.NewConsumerController(iConsumerUseCase)
	iKafkaPublisher := kafka.NewKafkaPublisher(kafkaParam)
	iConsumerHandler := handlers.MainConsumerHandler(consumerController, iKafkaPublisher)
	return iConsumerHandler
}

// wire.go:

var pkgSet = wire.NewSet(grpcclient.NewWrapGrpcClient, redisdb.NewRedisConn, kafka.NewKafkaPublisher, mongodb.NewMongoDB)

var pkgConsumerSet = wire.NewSet(grpcclient.NewWrapGrpcClient, redisdb.NewRedisConn, kafka.NewKafkaPublisher, mongodb.NewMongoDB)

var exampleSet = wire.NewSet(controllers.NewExampleController, usecases.NewExampleUseCase, usecases.NewExampleValidator, repositories.NewExampleRepository)

var storeSet = wire.NewSet(controllers.NewStoreController, usecases.NewStoreUseCase, repositories.NewStoreRepository)

var grpcServiceSet = wire.NewSet(services.NewExampleGrpcClient)

var grpcServerSet = wire.NewSet(grpcs.NewExampleGrpc, grpcs.NewStoreGrpc)

var consumerSet = wire.NewSet(controllers.NewConsumerController, usecases.NewConsumerUseCase)
