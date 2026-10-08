package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/user-service/internal/user"
)

func NewRouter(userHandler *user.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "user-service",
			"status":  "ok",
		})
	})

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.PUT("/:id", userHandler.Update)
		}
	}

	return router
}
