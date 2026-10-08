package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/auth"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/database"
	authgrpc "github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/grpc"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/http"
	userclient "github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/userclient"
	authv1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/auth/v1"
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

	userConn, err := grpc.NewClient(cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create user-service client: %v", err)
	}
	defer userConn.Close()

	repository := auth.NewPostgresRepository(db)
	tokenIssuer := auth.NewTokenIssuer(cfg.JWTSecret, cfg.TokenTTL)
	profiles := userclient.NewGRPCProfileClient(userConn)
	service := auth.NewService(repository, tokenIssuer, profiles)
	handler := auth.NewHandler(service)
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
	authv1.RegisterAuthServiceServer(grpcServer, authgrpc.NewAuthServer(service))

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting auth-service http on port %s", cfg.HTTPPort)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	log.Printf("starting auth-service grpc on port %s", cfg.GRPCPort)
	go func() {
		errCh <- grpcServer.Serve(grpcListener)
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("auth-service stopped: %v", err)
		}
	case <-ctx.Done():
		grpcServer.GracefulStop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown auth-service: %v", err)
		}
	}
}
