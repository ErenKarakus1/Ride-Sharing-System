package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/authclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/events"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	authConn, err := grpc.NewClient(cfg.AuthServiceGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create auth-service client: %v", err)
	}
	defer authConn.Close()

	validator := authclient.NewGRPCValidator(authConn)
	hub := notification.NewHub(cfg.AllowedOrigins)
	go hub.Run(ctx)

	service := notification.NewService(hub)
	consumer := events.NewKafkaConsumer(cfg.KafkaBrokers, cfg.ConsumerGroup, service)
	defer consumer.Close()

	router := httpapi.NewRouter(hub, validator)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 2)

	log.Printf("starting notification-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	log.Printf("starting ride event consumer")
	go func() {
		errCh <- consumer.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed && ctx.Err() == nil {
			log.Fatalf("notification-service stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown notification-service: %v", err)
		}
	}
}
