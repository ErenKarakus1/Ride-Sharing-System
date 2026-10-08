package main

import (
	"context"
	"log"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/config"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/database"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/user"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	repository := user.NewPostgresRepository(db)
	handler := user.NewHandler(repository)
	router := httpapi.NewRouter(handler)

	log.Printf("starting user-service on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("user-service stopped: %v", err)
	}
}
