package main

import (
	"log"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/config"
	httpapi "github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/http"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/user"
)

func main() {
	cfg := config.Load()
	repository := user.NewMemoryRepository()
	handler := user.NewHandler(repository)
	router := httpapi.NewRouter(handler)

	log.Printf("starting user-service on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("user-service stopped: %v", err)
	}
}
