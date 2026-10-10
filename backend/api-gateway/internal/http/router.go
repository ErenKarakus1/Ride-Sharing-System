package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/authclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(authValidator authclient.Validator, authProxy *proxy.Proxy, userProxy *proxy.Proxy, rideProxy *proxy.Proxy, locationProxy *proxy.Proxy, matchingProxy *proxy.Proxy, pricingProxy *proxy.Proxy, notificationProxy *proxy.Proxy, paymentProxy *proxy.Proxy) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestID())
	router.Use(cors())
	router.Use(requestTimeout(10 * time.Second))
	router.Use(bodyLimit(1 << 20))
	router.Use(requestLogger())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "api-gateway",
			"status":  "ok",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
			protected.Any("/payments/*path", gin.WrapH(paymentProxy.Handler()))
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
			"%s %s %d %s request_id=%s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Writer.Status(),
			time.Since(startedAt),
			ctx.Writer.Header().Get("X-Request-ID"),
		)
	}
}

func requestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = time.Now().UTC().Format("20060102150405.000000000")
		}

		ctx.Request.Header.Set("X-Request-ID", requestID)
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Next()
	}
}

func requestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/ws/notifications" {
			ctx.Next()
			return
		}

		requestContext, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(requestContext)
		ctx.Next()
	}
}

func bodyLimit(limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body != nil {
			ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, limit)
		}
		ctx.Next()
	}
}
