package log

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	boldWhite := color.New(color.FgWhite, color.Bold).SprintFunc()

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		p, _ := peer.FromContext(ctx)

		// Process the request
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err)
		colorCode := green(code)
		if code.String() != "OK" {
			colorCode = red(code)
		}

		method := info.FullMethod
		ip := "-"
		if p != nil {
			ip = p.Addr.String()
		}

		fmt.Printf("[gRPC] %s | %s | %s | %s | %s\n",
			boldWhite(time.Now().Format("2006/01/02 - 15:04:05")),
			colorCode,
			yellow(fmt.Sprintf("%v", duration)),
			cyan(ip),
			boldWhite(method),
		)

		return resp, err
	}
}
