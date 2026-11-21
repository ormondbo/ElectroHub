package authclient

import (
	"context"
	"fmt"
	"os"

	authpb "electrohub/internal/proto/auth"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(lifecycle fx.Lifecycle) (authpb.AuthServiceClient, error) {
	address := os.Getenv("AUTH_GRPC_ADDR")
	if address == "" {
		address = "localhost:50051"
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create auth grpc connection: %w", err)
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})
	return authpb.NewAuthServiceClient(conn), nil
}
