package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/database"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/events"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/ride"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	publisher := events.NewKafkaPublisher(cfg.KafkaBrokers)
	defer publisher.Close()

	repository := ride.NewPostgresRepository(db)
	service := ride.NewService(repository, publisher)
	handler := ride.NewHandler(service)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting ride-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ride-service stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown ride-service: %v", err)
		}
	}
}
