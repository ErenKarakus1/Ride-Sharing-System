package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	ridev1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/ride/v1"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/database"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/events"
	ridegrpc "github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/grpc"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/outbox"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/ride"
	"google.golang.org/grpc"
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
	outboxStore := outbox.NewEventStore(db)
	dispatcher := outbox.NewDispatcher(outboxStore, publisher)

	repository := ride.NewPostgresRepository(db)
	service := ride.NewService(repository, outboxStore)
	handler := ride.NewHandler(service)
	router := httpapi.NewRouter(handler)

	httpServer := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(ridegrpc.InternalAuthInterceptor(cfg.InternalToken)))
	ridev1.RegisterRideServiceServer(grpcServer, ridegrpc.NewRideServer(service))

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	log.Printf("starting ride-service http on port %s", cfg.HTTPPort)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	log.Printf("starting ride-service grpc on port %s", cfg.GRPCPort)
	go func() {
		errCh <- grpcServer.Serve(grpcListener)
	}()

	log.Printf("starting ride outbox dispatcher")
	go dispatcher.Run(ctx)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ride-service stopped: %v", err)
		}
	case <-ctx.Done():
		grpcServer.GracefulStop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown ride-service: %v", err)
		}
	}
}
