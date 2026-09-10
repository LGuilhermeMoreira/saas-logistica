package main

import (
	"context"
	"delivery/config"
	"delivery/di"
	"delivery/internal/infra/database"
	"delivery/pkg/logger"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	deliveryv1 "proto/gen/delivery/v1"
)

func main() {
	godotenv.Load()

	env := config.NewENV()

	logger := logger.New(env.LOG_MODE)

	logger.Info("starting delivery application", slog.String("port", env.PORT), slog.String("log_mode", env.LOG_MODE))

	// Connect to MongoDB
	mongoClient, err := database.NewMongoConnection(mongoURI(env))
	if err != nil {
		logger.Error("error connecting to MongoDB", slog.String("error", err.Error()))
		panic(err)
	}

	db := mongoClient.Database(env.DATABASE_NAME)

	deliveryTransport, err := di.InitDeliveryTransport(env, db, logger)
	if err != nil {
		logger.Error("error initializing delivery transport", slog.String("error", err.Error()))
		panic(err)
	}

	// Initialize gRPC middleware
	grpcMiddleware, err := di.InitGRPCMiddleware(env, logger)
	if err != nil {
		logger.Error("error initializing gRPC middleware", slog.String("error", err.Error()))
		panic(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", env.PORT))
	if err != nil {
		logger.Error("error listening to tcp port", slog.String("error", err.Error()), slog.String("port", env.PORT))
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcMiddleware.ValidateCredentials()),
	)

	deliveryv1.RegisterDeliveryServiceServer(server, deliveryTransport)

	go func() {
		logger.Info("gRPC server running", slog.String("port", env.PORT))
		if err := server.Serve(listener); err != nil {
			logger.Error("error serving gRPC", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown signal received, initiating graceful shutdown")
	server.GracefulStop()
	logger.Info("server gracefully stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Disconnect(ctx); err != nil {
		logger.Error("error disconnecting from MongoDB", slog.String("error", err.Error()))
	}
}

func mongoURI(env *config.Env) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=%s",
		env.DATABASE_USER,
		env.DATABASE_PASSWORD,
		env.DATABASE_HOST,
		env.DATABASE_PORT,
		env.DATABASE_NAME,
	)
}
