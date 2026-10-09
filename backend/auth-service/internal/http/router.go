package http

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(authHandler *auth.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "auth-service",
			"status":  "ok",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		authRoutes.Use(rateLimit(10, time.Minute))
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}
	}

	return router
}

type rateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	requests map[string][]time.Time
}

func rateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := &rateLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}

	return func(ctx *gin.Context) {
		key := clientIP(ctx) + ":" + ctx.FullPath()
		if !limiter.allow(key, time.Now()) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		ctx.Next()
	}
}

func (l *rateLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := now.Add(-l.window)
	recent := l.requests[key][:0]
	for _, timestamp := range l.requests[key] {
		if timestamp.After(cutoff) {
			recent = append(recent, timestamp)
		}
	}
	if len(recent) >= l.limit {
		l.requests[key] = recent
		return false
	}

	l.requests[key] = append(recent, now)
	return true
}

func clientIP(ctx *gin.Context) string {
	forwardedFor := ctx.GetHeader("X-Forwarded-For")
	if forwardedFor != "" {
		ip, _, _ := strings.Cut(forwardedFor, ",")
		return strings.TrimSpace(ip)
	}

	return ctx.ClientIP()
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
