package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/database"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/payment"
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

	repository := payment.NewPostgresRepository(db)
	service := payment.NewService(repository)
	handler := payment.NewHandler(service)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting payment-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("payment-service stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown payment-service: %v", err)
		}
	}
}
