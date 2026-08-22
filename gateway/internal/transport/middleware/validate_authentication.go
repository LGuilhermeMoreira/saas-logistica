package middleware

import (
	"gateway/pkg/authentication"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateAuthentication(validator authentication.TokenValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx.GetHeader("Authorization"))

		err := validator.Validate(token)
		ctx.Set("token", token)

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

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
