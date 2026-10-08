package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	userv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/user/v1"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/database"
	usergrpc "github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/grpc"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/user"
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

	repository := user.NewPostgresRepository(db)
	service := user.NewService(repository)
	handler := user.NewHandler(service)
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

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, usergrpc.NewUserServer(service))

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting user-service http on port %s", cfg.HTTPPort)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	log.Printf("starting user-service grpc on port %s", cfg.GRPCPort)
	go func() {
		errCh <- grpcServer.Serve(grpcListener)
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("user-service stopped: %v", err)
		}
	case <-ctx.Done():
		grpcServer.GracefulStop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown user-service: %v", err)
		}
	}
}
