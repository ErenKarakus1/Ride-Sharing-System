package http

import (
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/authclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

func NewRouter(authValidator authclient.Validator, authProxy *proxy.Proxy, userProxy *proxy.Proxy, rideProxy *proxy.Proxy, locationProxy *proxy.Proxy, matchingProxy *proxy.Proxy, pricingProxy *proxy.Proxy, notificationProxy *proxy.Proxy) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "api-gateway",
			"status":  "ok",
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", gin.WrapH(authProxy.Handler()))
		v1.POST("/auth/login", gin.WrapH(authProxy.Handler()))

		protected := v1.Group("")
		protected.Use(authMiddleware(authValidator))
		{
			protected.Any("/users", gin.WrapH(userProxy.Handler()))
			protected.Any("/users/*path", gin.WrapH(userProxy.Handler()))
			protected.Any("/rides", gin.WrapH(rideProxy.Handler()))
			protected.Any("/rides/*path", gin.WrapH(rideProxy.Handler()))
			protected.Any("/drivers/*path", gin.WrapH(locationProxy.Handler()))
			protected.Any("/matches", gin.WrapH(matchingProxy.Handler()))
			protected.Any("/fare-estimates", gin.WrapH(pricingProxy.Handler()))
		}
	}

	router.GET("/ws/notifications", gin.WrapH(notificationProxy.Handler()))

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
