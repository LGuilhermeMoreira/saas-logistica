package middleware

import (
	"gateway/pkg/authentication"
	"gateway/pkg/authorization"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorization(opa authorization.OPAInterface, jwt authentication.TokenValidator, log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.Writer.Header().Get("X-Request-ID")
		if reqID == "" {
			reqID = ctx.GetHeader("X-Request-ID")
		}

		logger := log.With("request_id", reqID)

		token, exists := ctx.Get("token")
		if !exists {
			logger.Error("token not found in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString, ok := token.(string)
		if !ok {
			logger.Error("token is not a valid string")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ExtractClaims(tokenString)
		if err != nil {
			logger.Error("failed to extract token claims", "error", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		data, ok := claims["data"].(map[string]any)
		if !ok {
			logger.Error("data claim is missing or invalid type")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		role, ok := data["role_name"].(string)
		if !ok {
			logger.Error("role_name claim is missing or invalid type")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		input := authorization.OPAInput{
			Action: ctx.Request.Method,
			Path:   ctx.Request.URL.Path,
		}

		input.User.Role = role

		if err := opa.Validate(input); err != nil {
			logger.Error("OPA validation failed", "error", err.Error())
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
