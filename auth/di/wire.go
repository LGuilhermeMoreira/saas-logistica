//go:build wireinject
// +build wireinject

package di

import (
	"auth/config"
	"auth/pkg/authentication"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func InitGRPCAuthTransport(env *config.Env, db *gorm.DB, s3Client *minio.Client, jwt *authentication.JWT)
