package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/events"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/locationclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/matching"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/matching-service/internal/rideclient"
)

func main() {
	cfg := config.Load()
	locationClient := locationclient.NewHTTPClient(cfg.LocationServiceURL)
	rideClient := rideclient.NewHTTPClient(cfg.RideServiceURL, cfg.InternalToken)
	service := matching.NewService(locationClient, rideClient)
	handler := matching.NewHandler(service)
	router := httpapi.NewRouter(handler)
	consumer := events.NewKafkaConsumer(cfg.KafkaBrokers, cfg.ConsumerGroup, service)
	defer consumer.Close()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	log.Printf("starting matching-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	log.Printf("starting matching-service payment event consumer")
	go func() {
		errCh <- consumer.Run(ctx)
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
