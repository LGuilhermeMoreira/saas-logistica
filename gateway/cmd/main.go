package main

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/infra/service"
	"gateway/internal/transport/api/auth"
	"gateway/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv1 "proto/gen/auth/v1"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()
	cnfg := config.NewENV()

	logger := logger.New(cnfg.LOG_MODE)
	slog.SetDefault(logger)

	slog.Info("starting gateway application", slog.String("port", cnfg.PORT), slog.String("log_mode", cnfg.LOG_MODE))

	authConn, err := grpc.NewClient(cnfg.AUTH_BASE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to connect to auth service", "error", err)
		panic(err)
	}
	defer authConn.Close()

	authClient := authv1.NewAuthServiceClient(authConn)
	authService := service.NewAuthService(authClient)
	authController := auth.NewAuthController(authService)

	mux := gin.New()
	mux.POST("auth/login", authController.Login)
	authGroup := mux.Group("/auth")
	{
		authGroup.POST("/create-user", authController.CreateUser)
		authGroup.POST("/create-role", authController.CreateRole)
		authGroup.DELETE("/delete-user/:id", authController.DeleteUser)
		authGroup.DELETE("/delete-role/:id", authController.DeleteRole)

	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cnfg.PORT),
		Handler: mux,
	}

	go func() {
		slog.Info("gateway server is listening", slog.String("port", cnfg.PORT))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("gateway server failed", "error", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown signal received, shutting down gateway server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("gateway server shutdown failed", "error", err)
		panic(err)
	}

	slog.Info("gateway server gracefully stopped")
}
