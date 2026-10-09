package http

import (
	"net/http"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/authclient"
	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	allowedOrigins := map[string]struct{}{
		"http://localhost:3000": {},
		"http://localhost:5173": {},
		"http://127.0.0.1:3000": {},
		"http://127.0.0.1:5173": {},
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Vary", "Origin")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func authMiddleware(validator authclient.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := bearerToken(ctx.GetHeader("Authorization"))
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := validator.Validate(ctx.Request.Context(), accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		ctx.Request.Header.Set("X-User-ID", claims.UserID)
		ctx.Request.Header.Set("X-User-Email", claims.Email)
		ctx.Request.Header.Set("X-User-Role", claims.Role)
		ctx.Next()
	}
}

func bearerToken(header string) string {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}
