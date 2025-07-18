package main

import (
	"appsku-golang/app/config"
	"appsku-golang/app/constants"
	"appsku-golang/app/handlers"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"appsku-golang/app/global-utils/grpcclient"
	"appsku-golang/app/global-utils/helper"
	"appsku-golang/app/global-utils/log"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfig()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	log.InitLog(config.Get().LogLevel)
}

func main() {
	name := flag.String("name", "", "Consumer name")
	flag.Parse()

	cfg := config.Get()

	grpcClientParams := []grpcclient.GRPCClientParam{
		config.BuildGrpcClientParam(constants.ExampleService)}

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, "main")
	} else if args[0] != "main" && args[0] != "consumer" && args[0] != "grpc" && args[0] != "fiber" {
		args[0] = "main"
	}

	switch args[0] {
	case "main":
		// Default to Gin
		ginEngine := InitializeHttpServer(
			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		grpcServer := InitializeGrpcServer(
			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		signalExit := make(chan os.Signal, 1)

		signal.Notify(signalExit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		//go startGrpc(context.Background(), grpcServer, cfg)
		go startServer(context.Background(), ginEngine, cfg)

		<-signalExit

		helper.ExitHTTP <- true

		grpcServer.GracefulStop()
		time.Sleep(2 * time.Second)
	case "fiber":
		fiberApp := InitializeFiberServer(
			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		grpcServer := InitializeGrpcServer(
			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		signalExit := make(chan os.Signal, 1)

		signal.Notify(signalExit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		//go startGrpc(context.Background(), grpcServer, cfg)
		go startServer(context.Background(), fiberApp, cfg)

		<-signalExit

		helper.ExitHTTP <- true

		grpcServer.GracefulStop()
		time.Sleep(2 * time.Second)
	case "grpc":
		grpcServer := InitializeGrpcServer(
			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		signalExit := make(chan os.Signal, 1)
		signal.Notify(signalExit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		startGrpc(context.Background(), grpcServer, cfg)

		<-signalExit

		grpcServer.GracefulStop()
		time.Sleep(2 * time.Second)
	case "consumer":
		consumer := InitializeConsumer(

			grpcClientParams,
			config.BuildRedisParam(),
			config.BuildKafkaParam(),
			config.BuildMongoDBParam(),
		)

		signal.Notify(helper.ExitKafka, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		startConsumer(context.Background(), consumer, cfg, *name)
	}
}

func startServer(ctx context.Context, server interface{}, cfg config.Configuration) {
	// Check the type of server
	switch s := server.(type) {
	case *fiber.App:
		go func() {
			logrus.Info("Running Fiber server service on port ", cfg.MainPort)
			addr := fmt.Sprintf(":%d", cfg.MainPort)
			if cfg.UseSSL {
				if err := s.ListenTLS(addr, cfg.PublicSSLPath, cfg.PrivateSSLPath); err != nil {
					logrus.Error(err)
					logrus.Fatal("shutting down Fiber server tls")
				}
			} else {
				if err := s.Listen(addr); err != nil {
					logrus.Error(err)
					logrus.Fatal("shutting down Fiber server")
				}
			}
		}()

		// Wait for exit signal
		<-helper.ExitHTTP
		logrus.WithField(helper.GetRequestIDContext(ctx)).Infoln("Wait for http process done")

		if err := s.Shutdown(); err != nil {
			logrus.Fatal(err)
		}

		logrus.WithField(helper.GetRequestIDContext(ctx)).Infoln("http already exited")
	case *gin.Engine:
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.MainPort),
			Handler: s,
		}

		go func() {
			logrus.Info("Running Gin server service on port ", srv.Addr)
			if cfg.UseSSL {
				if err := srv.ListenAndServeTLS(cfg.PublicSSLPath, cfg.PrivateSSLPath); err != nil && err != http.ErrServerClosed {
					logrus.Error(err)
					logrus.Fatal("shutting down Gin server tls")
				}
			} else {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.Error(err)
					logrus.Fatal("shutting down Gin server")
				}
			}
		}()

		// Wait for exit signal
		<-helper.ExitHTTP
		logrus.WithField(helper.GetRequestIDContext(ctx)).Infoln("Wait for http process done")

		ctx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Fatal(err)
		}

		logrus.WithField(helper.GetRequestIDContext(ctx)).Infoln("http already exited")
	default:
		logrus.Fatal("Unknown server type")
	}
}

func startGrpc(ctx context.Context, g *grpc.Server, cfg config.Configuration) {
	addr := fmt.Sprintf(":%d", cfg.Grpc.Server.Port)
	listen, err := net.Listen("tcp", addr)

	if err != nil {
		logrus.WithField(helper.GetRequestIDContext(ctx)).Fatalln("Error when listen server address : ", err.Error())
	}

	go func() {
		logrus.Info("Running server GRPC service on port ", addr)
		if err := g.Serve(listen); err != nil {
			logrus.Error(err)
			logrus.Fatal("shutting down grpc server")
		}
	}()

	// <-helper.ExitGrpcClient
}

func startConsumer(ctx context.Context, consumerHandler handlers.IConsumerHandler, cfg config.Configuration, name string) {
	consumerHandler.BindConsumer(ctx, cfg, name)
}
