package internal

import (
	"context"
	"electrohub/internal/handler"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"
)

func RegisterHTTPServer(lifecycle fx.Lifecycle, api *handler.API) {
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}

	server := http.Server{
		Addr:              address,
		Handler:           api,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

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
