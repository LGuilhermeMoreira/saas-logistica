//go:build wireinject
// +build wireinject

package di

import (
	"log/slog" // Importe o slog aqui

	"auth/config"
	"auth/internal/application/usecase"
	"auth/internal/infra/repository"
	"auth/internal/infra/service"
	"auth/internal/transport"
	"auth/pkg/authentication"
	"auth/pkg/authorization"
	grpcmiddleware "auth/pkg/grpc_middleware"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger" (Pode remover se não for usar o logger do gorm aqui)
)

var AuthProviderSet = wire.NewSet(
	authentication.NewJWT,
	wire.Bind(new(authentication.TokenValidator), new(*authentication.JWT)),
	wire.Bind(new(authentication.TokenGenerator), new(*authentication.JWT)),
)

// Adicionei log *slog.Logger nos parâmetros
func InitGRPCAuthTransport(env *config.Env, db *gorm.DB, s3Client *minio.Client, log *slog.Logger) (*transport.AuthTransport, error) {
	wire.Build(
		// JWT
		AuthProviderSet,
		// User
		repository.NewUserRepository, usecase.NewUserUsecase,
		// Role
		repository.NewRoleRepository,
		service.NewStorageService,
		service.NewOPASyncService,
		usecase.NewRoleUsecase,
		// OPA
		// authorization.NewOPA,
		// gRPC
		transport.NewAuthTransport,
	)
	return nil, nil
}

var ValidateCredentialsProviderSet = wire.NewSet(
	authentication.NewJWT,
	wire.Bind(new(authentication.TokenValidator), new(*authentication.JWT)),
)

// Se o seu middleware gRPC também precisar do logger, adicione o parâmetro aqui também!
func InitGRPCMiddleware(env *config.Env, log *slog.Logger) (*grpcmiddleware.GRPCMiddleware, error) {
	wire.Build(
		ValidateCredentialsProviderSet,
		authorization.NewOPA,
		grpcmiddleware.NewGRPCMiddleware,
	)
	return nil, nil
}
