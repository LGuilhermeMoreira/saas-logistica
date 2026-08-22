package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func New(logMode string) *slog.Logger {
	var handler slog.Handler
	switch logMode {
	case "dev":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}

	logger := slog.New(handler)
	return logger
}

type contextKey string

const RequestIDKey contextKey = "request_id"

func ExtractRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		return reqID
	}
	return "fallback-" + uuid.New().String()
}
