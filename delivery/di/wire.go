//go:build wireinject
// +build wireinject

package di

import (
	"log/slog"

	"delivery/config"
	"delivery/internal/application/usecase"
	"delivery/internal/infra/repository"
	"delivery/internal/transport"
	"delivery/pkg/authentication"
	"delivery/pkg/authorization"
	grpcmiddleware "delivery/pkg/grpc_middleware"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var JWTProviderSet = wire.NewSet(
	authentication.NewJWT,
	wire.Bind(new(authentication.TokenValidator), new(*authentication.JWT)),
)

func InitDeliveryTransport(env *config.Env, db *mongo.Database, log *slog.Logger) (*transport.DeliveryTransport, error) {
	wire.Build(
		repository.NewDeliveryRepository,
		usecase.NewDeliveryUsecase,
		authorization.NewOPA,
		transport.NewDeliveryTransport,
	)
	return nil, nil
}

func InitGRPCMiddleware(env *config.Env, log *slog.Logger) (*grpcmiddleware.GRPCMiddleware, error) {
	wire.Build(
		JWTProviderSet,
		authorization.NewOPA,
		grpcmiddleware.NewGRPCMiddleware,
	)
	return nil, nil
}
