package http

import (
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/location-service/internal/location"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(locationHandler *location.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "location-service",
			"status":  "ok",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		drivers := v1.Group("/drivers")
		{
			drivers.PUT("/:id/location", locationHandler.UpdateDriverLocation)
			drivers.POST("/:id/available", locationHandler.SetDriverAvailable)
			drivers.POST("/:id/unavailable", locationHandler.SetDriverUnavailable)
			drivers.POST("/:id/claim", locationHandler.ClaimDriver)
		}

		v1.GET("/drivers/nearby", locationHandler.NearbyDrivers)
	}

	return router
}

func requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startedAt := time.Now()
		ctx.Next()

		log.Printf(
			"%s %s %d %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Writer.Status(),
			time.Since(startedAt),
		)
	}
}
