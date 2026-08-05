package main

import (
	"auth/config"
	"auth/di"
	"auth/internal/infra/database"
	"auth/internal/infra/s3"
	"fmt"
	"log"
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

	dbConn, err := database.NewPostgresConnection(cnfg.PostgresURI())

	if err != nil {
		panic(err)
	}

	minioClient, err := s3.NewS3Connection(cnfg)
	if err != nil {
		panic(err)
	}

	authTransport, err := di.InitGRPCAuthTransport(cnfg, dbConn, minioClient)
	if err != nil {
		panic(err)
	}

	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", cnfg.PORT))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	authv1.RegisterAuthServiceServer(server, authTransport)

	go func() {
		log.Printf("gRPC server running on port %s\n", cnfg.PORT)
		if err := server.Serve(conn); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received. Initiating graceful shutdown...")
	server.GracefulStop()
	log.Println("Server gracefully stopped.")
}
