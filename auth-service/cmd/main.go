package main

import (
	"auth-service/internal/di"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	app := fx.New(di.Module, fx.Invoke(func(*grpc.Server) {}))
	app.Run()
}
