package http

import (
	"net/http"
	"strings"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/api-gateway/internal/authclient"
	"github.com/gin-gonic/gin"
)

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
