package di

import (
	"auth-service/internal"
	"auth-service/internal/api"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/jwtauth"

	"go.uber.org/fx"
)

var Module = fx.Module("auth",
	fx.Provide(
		postgres.NewPostgresPool,
		postgres.NewUserRepository,
		jwtauth.NewJwtProvider,
		service.NewAuthService,
		api.NewHandler,
		internal.NewGRPCServer,
	),
)
