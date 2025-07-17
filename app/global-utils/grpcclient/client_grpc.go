package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientParam struct {
	Name               string
	Host               string
	Port               string
	IsUseSSL           bool
	ProxyPath          string
	InsecureSkipVerify bool
}

func NewGrpcClient(param GRPCClientParam) *grpc.ClientConn {
	grpcHost := fmt.Sprintf("%s:%s", param.Host, param.Port)
	ctx := context.Background()

	var (
		grpcConnection *grpc.ClientConn
		err            error
	)

	if net.ParseIP(param.Host) == nil {
		param.IsUseSSL = true
	}

	if !param.IsUseSSL {
		grpcConnection, err = grpc.DialContext(ctx, grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Panic(err)
			return &grpc.ClientConn{}
		}

	} else {
		f, err := os.ReadFile(param.ProxyPath)
		if err != nil {
			log.Panic(err)
			return &grpc.ClientConn{}
		}

		p := x509.NewCertPool()
		p.AppendCertsFromPEM(f)
		tlsConfig := &tls.Config{
			RootCAs:            p,
			InsecureSkipVerify: param.InsecureSkipVerify,
		}

		grpcConnection, err = grpc.DialContext(ctx, grpcHost, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			log.Panic(err)
			return &grpc.ClientConn{}
		}
	}

	// Check if the connection not ready
	if grpcConnection.GetState() != connectivity.Ready && grpcConnection.GetState() != connectivity.Idle {
		log.Printf("Connection to grpc server %s is not ready.", param.Name)
	}

	ticker := time.NewTicker(30 * time.Second)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for a signal
		sig := <-sigCh
		log.Printf("Received signal: %v. Closing the ticker.", sig)

		// Stop the ticker
		ticker.Stop()
	}()

	go func(grpcConnection *grpc.ClientConn) {
		for range ticker.C {
			select {
			case <-ctx.Done():
				log.Println("Stopping the reconnection routine.")
				ticker.Stop()
			case <-ticker.C:
				// Check if the connection is still alive
				if grpcConnection.GetState() == connectivity.Ready || grpcConnection.GetState() == connectivity.Idle {
					fmt.Println("Connection is still alive.")
					continue
				}

				var newConn *grpc.ClientConn

				// Attempt to create a new connection
				if !param.IsUseSSL {
					newConn, err = grpc.DialContext(ctx, grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						log.Panic(err)
					}

				} else {
					f, err := os.ReadFile(param.ProxyPath)
					if err != nil {
						log.Panic(err)
					}

					p := x509.NewCertPool()
					p.AppendCertsFromPEM(f)
					tlsConfig := &tls.Config{
						RootCAs:            p,
						InsecureSkipVerify: param.InsecureSkipVerify,
					}

					newConn, err = grpc.DialContext(ctx, grpcHost, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
					if err != nil {
						log.Panic(err)
					}
				}

				// Replace the old connection with the new one
				if grpcConnection.GetState() != connectivity.Ready && grpcConnection.GetState() != connectivity.Idle {
					log.Printf("Failed reconnecting to grpc server %s.", param.Name)
					continue
				}

				log.Println("Successfully reconnected.")

				grpcConnection.Close()
				grpcConnection = newConn
			}
		}
	}(grpcConnection)

	return grpcConnection
}

func NewWrapGrpcClient(params []GRPCClientParam) IGRPCClients {
	if len(params) == 0 {
		log.Fatalf("Params grpc connection cannot empty")
		return nil
	}
	grpcConnection := make(map[string]*grpc.ClientConn, len(params))

	for _, param := range params {
		grpcConnection[param.Name] = NewGrpcClient(param)
	}

	return &GRPCClients{
		GRPC: grpcConnection,
	}
}

type IGRPCClients interface {
	Client(name string) *grpc.ClientConn
}

type GRPCClients struct {
	GRPC map[string]*grpc.ClientConn
}

func (gc *GRPCClients) Client(name string) *grpc.ClientConn {
	return gc.GRPC[name]
}
