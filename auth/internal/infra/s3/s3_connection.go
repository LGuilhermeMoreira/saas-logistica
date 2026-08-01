package s3

import (
	"auth/config"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewS3Connection(env *config.Env) (*minio.Client, error) {
	minioClient, err := minio.New(env.STORAGE_URI, &minio.Options{
		Creds:  credentials.NewStaticV4(env.STORAGE_ACCESS_KEY_ID, env.STORAGE_SECRECT_ACCESS_KEY, ""),
		Secure: env.STORAGE_USE_SSL,
	})

	if err != nil {
		return nil, fmt.Errorf("Error creating connection to S3 Bucket")
	}

	return minioClient, nil
}
