package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()

		args := []any{
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", status,
			"duration", duration,
			"client_ip", ctx.ClientIP(),
			"user_agent", ctx.Request.UserAgent(),
		}

		if len(ctx.Errors) > 0 {
			args = append(args, "error", ctx.Errors.String())

			slog.Error(
				"http request",
				args...,
			)

			return
		}

		slog.Info(
			"http request",
			args...,
		)
	}
}
