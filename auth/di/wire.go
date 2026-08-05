//go:build wireinject
// +build wireinject

package di

import (
	"auth/config"
	"auth/internal/application/usecase"
	"auth/internal/infra/repository"
	"auth/internal/infra/service"
	"auth/internal/transport"
	"auth/pkg/authentication"
	"auth/pkg/authorization"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

var AuthProviderSet = wire.NewSet(
	authentication.NewJWT,
	wire.Bind(new(authentication.TokenValidator), new(*authentication.JWT)),
	wire.Bind(new(authentication.TokenGenerator), new(*authentication.JWT)),
)

func InitGRPCAuthTransport(env *config.Env, db *gorm.DB, s3Client *minio.Client) (*transport.AuthTransport, error) {
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
		authorization.NewOPA,
		// gRPC
		transport.NewAuthTransport,
	)
	return nil, nil
}
