package di

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"shop/internal/api/httpapi"
	"shop/internal/client/authclient"
	"shop/internal/repository"
	"shop/internal/repository/postgres"
	"shop/internal/service"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"shop",
	fx.Provide(
		postgres.NewPostgresPool,

		fx.Annotate(
			func() (*authclient.Client, error) {
				address := os.Getenv("AUTH_GRPC_ADDR")
				if address == "" {
					address = "localhost:50051"
				}
				return authclient.New(address)
			},
			fx.As(new(httpapi.AuthClient)),
		),
		fx.Annotate(
			postgres.NewClientRepo,
			fx.As(new(repository.ClientRepository)),
		),
		fx.Annotate(
			postgres.NewProductRepo,
			fx.As(new(repository.ProductRepository)),
		),
		fx.Annotate(
			postgres.NewSupplierRepo,
			fx.As(new(repository.SupplierRepository)),
		),
		fx.Annotate(
			postgres.NewImageRepo,
			fx.As(new(repository.ImageRepository)),
		),
		fx.Annotate(
			postgres.NewAddressRepo,
			fx.As(new(repository.AddressRepository)),
		),

		service.NewClientService,
		service.NewProductService,
		service.NewSupplierService,
		service.NewImageService,

		httpapi.New,
		newHTTPServer,
	),
	fx.Invoke(
		registerHTTPServer,
		registerAuthClient,
	),
)

func newHTTPServer(api *httpapi.API) *http.Server {
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}

	return &http.Server{
		Addr:              address,
		Handler:           api,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func registerHTTPServer(lifecycle fx.Lifecycle, server *http.Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			listener, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go func() {
				log.Printf("HTTP server listening on %s", server.Addr)
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("HTTP server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func registerAuthClient(lc fx.Lifecycle, client httpapi.AuthClient) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			if closer, ok := client.(interface{ Close() error }); ok {
				return closer.Close()
			}
			return nil
		},
	})
}
