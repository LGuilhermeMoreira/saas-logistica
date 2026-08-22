package middleware

import (
	"gateway/pkg/authentication"
	"gateway/pkg/authorization"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateAuthorization(opa authorization.OPAInterface, jwt authentication.TokenValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exists := ctx.Get("token")
		if !exists {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString, ok := token.(string)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ExtractClaims(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		data, ok := claims["data"].(map[string]any)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		role, ok := data["role_name"].(string)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		input := authorization.OPAInput{
			Action: ctx.Request.Method,
			Path:   ctx.Request.URL.Path,
		}

		input.User.Role = role

		if err := opa.Validate(input); err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
