package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/config"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/locationclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/matching"
)

func main() {
	cfg := config.Load()
	locationClient := locationclient.NewHTTPClient(cfg.LocationServiceURL)
	service := matching.NewService(locationClient)
	handler := matching.NewHandler(service)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting matching-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("matching-service stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown matching-service: %v", err)
		}
	}
}
