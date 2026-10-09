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
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/outbox"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/payment"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/rideclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	rideConn, err := grpc.NewClient(cfg.RideServiceGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create ride-service client: %v", err)
	}
	defer rideConn.Close()
	rideClient := rideclient.NewGRPCClient(rideConn, cfg.InternalToken)

	repository := payment.NewPostgresRepository(db)
	publisher := events.NewKafkaPublisher(cfg.KafkaBrokers)
	defer publisher.Close()
	outboxStore := outbox.NewEventStore(db)
	dispatcher := outbox.NewDispatcher(outboxStore, publisher)

	service := payment.NewService(repository, outboxStore, rideClient)
	consumer := events.NewKafkaConsumer(cfg.KafkaBrokers, cfg.ConsumerGroup, service)
	defer consumer.Close()

	handler := payment.NewHandler(service)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	log.Printf("starting payment-service on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	log.Printf("starting payment ride event consumer")
	go func() {
		errCh <- consumer.Run(ctx)
	}()

	log.Printf("starting payment outbox dispatcher")
	go dispatcher.Run(ctx)

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
