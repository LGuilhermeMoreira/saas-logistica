package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func ExtractRequestValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Header("X-Request-ID", reqID)

		token := c.GetHeader("Authorization")

		httpPath := c.Request.URL.Path
		httpMethod := c.Request.Method

		md := metadata.New(map[string]string{
			"x-request-id":  reqID,
			"authorization": token,
			"http-path":     httpPath,
			"http-method":   httpMethod,
		})

		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
