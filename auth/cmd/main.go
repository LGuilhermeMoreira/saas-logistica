package main

import (
	"auth/config"
	"auth/di"
	"auth/internal/domain/entity"
	"auth/internal/infra/database"
	"auth/internal/infra/s3"
	"auth/pkg/logger"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	authv1 "proto/gen/auth/v1"
)

func main() {

	godotenv.Load()

	cnfg, err := config.NewEnv()
	if err != nil {
		panic(err)
	}

	logger := logger.New(cnfg.LOG_MODE)

	logger.Info("starting application with configurations", slog.String("port", cnfg.PORT), slog.String("log_mode", cnfg.LOG_MODE))

	dbConn, err := database.NewPostgresConnection(cnfg.PostgresURI())
	if err != nil {
		logger.Error("error connecting to database", slog.String("error", err.Error()))
		panic(err)
	}

	err = dbConn.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
	)
	if err != nil {
		logger.Error("error migrating entities to database", slog.String("error", err.Error()))
		panic(err)
	}

	minioClient, err := s3.NewS3Connection(cnfg)
	if err != nil {
		logger.Error("error connecting to minio", slog.String("error", err.Error()))
		panic(err)
	}

	authTransport, err := di.InitGRPCAuthTransport(cnfg, dbConn, minioClient, logger)
	if err != nil {
		logger.Error("error initializing GRPC auth transport", slog.String("error", err.Error()))
		panic(err)
	}

	grpcMiddleware, err := di.InitGRPCMiddleware(cnfg, logger)
	if err != nil {

	}

	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", cnfg.PORT))
	if err != nil {
		logger.Error("error listening to tcp port", slog.String("error", err.Error()))
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcMiddleware.ValidateCredentials()),
	)

	authv1.RegisterAuthServiceServer(server, authTransport)

	go func() {
		logger.Info(fmt.Sprintf("gRPC server running on port %s", cnfg.PORT))
		if err := server.Serve(conn); err != nil {
			logger.Error("error serving grpc", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received. Initiating graceful shutdown...")
	server.GracefulStop()
	logger.Info("Server gracefully stopped.")
}
