package internal

import (
	"auth-service/internal/api"
	"context"

	authpb "contracts/proto/auth"

	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewGRPCServer(lifecycle fx.Lifecycle, handler *api.Handler) *grpc.Server {
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, handler)
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", ":50051")
			if err != nil {
				return fmt.Errorf("listen on port 50051: %w", err)
			}
			go func() {
				if err := server.Serve(listener); err != nil {
					log.Printf("gRPC authentication server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			return nil
		},
	})
	return server
}
