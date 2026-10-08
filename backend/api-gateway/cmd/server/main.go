package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/authclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/config"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	authConn, err := grpc.NewClient(cfg.AuthServiceGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create auth-service client: %v", err)
	}
	defer authConn.Close()

	authValidator := authclient.NewGRPCValidator(authConn)
	authProxy := proxy.New(cfg.AuthServiceHTTPURL)
	userProxy := proxy.New(cfg.UserServiceHTTPURL)
	rideProxy := proxy.New(cfg.RideServiceHTTPURL)
	locationProxy := proxy.New(cfg.LocationServiceHTTPURL)
	matchingProxy := proxy.New(cfg.MatchingServiceHTTPURL)
	pricingProxy := proxy.New(cfg.PricingServiceHTTPURL)
	notificationProxy := proxy.New(cfg.NotificationServiceHTTPURL)
	paymentProxy := proxy.New(cfg.PaymentServiceHTTPURL)
	router := httpapi.NewRouter(authValidator, authProxy, userProxy, rideProxy, locationProxy, matchingProxy, pricingProxy, notificationProxy, paymentProxy)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("starting api-gateway on port %s", cfg.Port)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("api-gateway stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown api-gateway: %v", err)
		}
	}
}
