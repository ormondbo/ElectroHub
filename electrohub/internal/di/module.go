package di

import (
	"electrohub/internal"
	"electrohub/internal/client/authclient"
	"electrohub/internal/handler"
	"electrohub/internal/repository/postgres"
	"electrohub/internal/service"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"electrohub_api",
	fx.Provide(
		//PostgressPool
		postgres.NewPostgresPool,

		//Repositories <= PostgresPool
		postgres.NewClientRepo,
		postgres.NewProductRepo,
		postgres.NewSupplierRepo,
		postgres.NewImageRepo,
		postgres.NewAddressRepo,

		//Services <= Repositories
		service.NewClientService,
		service.NewProductService,
		service.NewSupplierService,
		service.NewImageService,

		//AuthClient
		authclient.NewAuthClient,
		//AuthService <= AuthClient (Services+)
		service.NewAuthService,

		//API <= Services
		handler.NewAPI,
	),
	fx.Invoke(
		internal.RegisterHTTPServer,
	),
)
