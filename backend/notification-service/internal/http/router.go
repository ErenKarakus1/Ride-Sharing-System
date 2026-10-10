package http

import (
	"net/http"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/authclient"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/notification"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(hub *notification.Hub, validator authclient.Validator) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": "notification-service",
			"status":  "ok",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/ws/notifications", func(ctx *gin.Context) {
		accessToken := accessToken(ctx)
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := validator.Validate(ctx.Request.Context(), accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		hub.HandleWebSocket(ctx, claims.UserID)
	})

	return router
}

func accessToken(ctx *gin.Context) string {
	token := bearerToken(ctx.GetHeader("Authorization"))
	if token != "" {
		return token
	}

	return strings.TrimSpace(ctx.Query("token"))
}

func bearerToken(header string) string {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}
