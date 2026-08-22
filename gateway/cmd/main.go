package main

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/infra/service"
	"gateway/internal/transport/api/auth"
	"gateway/internal/transport/middleware"
	"gateway/pkg/authentication"
	"gateway/pkg/authorization"
	"gateway/pkg/logger"
	"io"
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
	// slog.SetDefault(logger)

	logger.Info("starting gateway application", slog.String("port", cnfg.PORT), slog.String("log_mode", cnfg.LOG_MODE))

	authConn, err := grpc.NewClient(cnfg.AUTH_BASE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to auth service", "error", err)
		panic(err)
	}
	defer authConn.Close()

	authClient := authv1.NewAuthServiceClient(authConn)
	authService := service.NewAuthService(authClient)
	authController := auth.NewAuthController(authService, logger)

	jwt := authentication.NewJWT(cnfg)
	opa := authorization.NewOPA(cnfg)

	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	mux := gin.New()

	mux.Use(middleware.Logger())
	mux.Use(middleware.ExtractRequestValues())

	mux.POST("auth/login", authController.Login)

	authGroup := mux.Group("/auth")
	authGroup.Use(middleware.ValidateAuthentication(jwt))
	authGroup.Use(middleware.ValidateAuthorization(opa, jwt))
	{
		authGroup.POST("/create-user", authController.CreateUser)
		authGroup.POST("/create-role", authController.CreateRole)
		authGroup.DELETE("/delete-user/:id", authController.DeleteUser)
		authGroup.DELETE("/delete-role/:id", authController.DeleteRole)
		authGroup.GET("/find-role/:id", authController.FindRoleByID)
		authGroup.GET("/find-user/:id", authController.FindUserByID)
		authGroup.GET("/find-all-roles", authController.FindAllRoles)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cnfg.PORT),
		Handler: mux,
	}

	go func() {
		logger.Info("gateway server is listening", slog.String("port", cnfg.PORT))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("gateway server failed", "error", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown signal received, shutting down gateway server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("gateway server shutdown failed", "error", err)
		panic(err)
	}

	logger.Info("gateway server gracefully stopped")
}
