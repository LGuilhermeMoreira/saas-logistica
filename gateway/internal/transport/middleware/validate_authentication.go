package middleware

import (
	"gateway/pkg/authentication"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateAuthentication(validator authentication.TokenValidator, log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.Writer.Header().Get("X-Request-ID")
		if reqID == "" {
			reqID = ctx.GetHeader("X-Request-ID")
		}

		logger := log.With("request_id", reqID)

		token := extractToken(ctx.GetHeader("Authorization"))
		if token == "" {
			logger.Error("authorization token is missing or empty")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err := validator.Validate(token)
		if err != nil {
			logger.Error("failed to validate token", "error", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("token", token)

		ctx.Next()
	}
}

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
