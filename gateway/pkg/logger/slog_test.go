package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestNew_ReturnsNonNilLogger(t *testing.T) {
	modes := []string{"dev", "prod", "outro", ""}

	for _, mode := range modes {
		mode := mode
		t.Run(mode, func(t *testing.T) {
			logger := New(mode)
			if logger == nil {
				t.Fatalf("New(%q) retornou nil, esperado um *slog.Logger", mode)
			}
		})
	}
}

func TestNew_DevMode(t *testing.T) {
	logger := New("dev")
	handler := logger.Handler()

	if _, ok := handler.(*slog.TextHandler); !ok {
		t.Errorf("esperado *slog.TextHandler no modo dev, obtido %T", handler)
	}

	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("esperado nível Debug habilitado no modo dev")
	}
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("esperado nível Info habilitado no modo dev")
	}
}

func TestNew_ProdMode(t *testing.T) {
	logger := New("prod")
	handler := logger.Handler()

	if _, ok := handler.(*slog.JSONHandler); !ok {
		t.Errorf("esperado *slog.JSONHandler no modo prod, obtido %T", handler)
	}

	ctx := context.Background()
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("esperado nível Debug desabilitado no modo prod")
	}
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("esperado nível Info habilitado no modo prod")
	}
	if !handler.Enabled(ctx, slog.LevelWarn) {
		t.Error("esperado nível Warn habilitado no modo prod")
	}
	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("esperado nível Error habilitado no modo prod")
	}
}

func TestNew_DefaultMode(t *testing.T) {
	unknownModes := []string{"staging", "TEST", "123", " "}

	for _, mode := range unknownModes {
		mode := mode
		t.Run(mode, func(t *testing.T) {
			logger := New(mode)
			handler := logger.Handler()

			if _, ok := handler.(*slog.TextHandler); !ok {
				t.Errorf("esperado *slog.TextHandler para logMode %q, obtido %T", mode, handler)
			}

			ctx := context.Background()
			if !handler.Enabled(ctx, slog.LevelDebug) {
				t.Errorf("esperado nível Debug habilitado para logMode %q", mode)
			}
		})
	}
}

func TestNew_EmptyMode(t *testing.T) {
	logger := New("")
	handler := logger.Handler()

	if _, ok := handler.(*slog.TextHandler); !ok {
		t.Errorf("esperado *slog.TextHandler para logMode vazio, obtido %T", handler)
	}
}

func TestNew_CaseSensitivity(t *testing.T) {
	logger := New("PROD")
	handler := logger.Handler()

	if _, ok := handler.(*slog.TextHandler); !ok {
		t.Errorf("esperado *slog.TextHandler para logMode \"PROD\" (case-sensitive), obtido %T", handler)
	}
}
