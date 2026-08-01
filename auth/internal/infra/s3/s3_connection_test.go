package s3_test

import (
	"auth/config"
	"auth/internal/infra/s3"
	"context"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	miniotc "github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestNewS3Connection(t *testing.T) {
	ctx := context.Background()

	const (
		accessKey = "testuser"
		secretKey = "testpassword"
	)

	minioContainer, err := miniotc.Run(ctx,
		"minio/minio:RELEASE.2024-01-16T16-07-38Z",
		miniotc.WithUsername(accessKey),
		miniotc.WithPassword(secretKey),
	)
	require.NoError(t, err)
	defer func() {
		if err := minioContainer.Terminate(ctx); err != nil {
			t.Fatalf("erro ao terminar container: %v", err)
		}
	}()

	endpoint, err := minioContainer.ConnectionString(ctx)
	require.NoError(t, err)

	newEnv := func() *config.Env {
		return &config.Env{
			STORAGE_URI:                endpoint,
			STORAGE_ACCESS_KEY_ID:      accessKey,
			STORAGE_SECRECT_ACCESS_KEY: secretKey,
			STORAGE_USE_SSL:            false,
		}
	}

	t.Run("connects successfully", func(t *testing.T) {
		client, err := s3.NewS3Connection(newEnv())
		require.NoError(t, err)
		require.NotNil(t, client)

		ctxCheck, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		buckets, err := client.ListBuckets(ctxCheck)
		require.NoError(t, err)
		require.Empty(t, buckets)
	})

	t.Run("Creates a bucket and checks for its existence.", func(t *testing.T) {
		client, err := s3.NewS3Connection(newEnv())
		require.NoError(t, err)

		bucketName := "test-bucket"
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		require.NoError(t, err)

		exists, err := client.BucketExists(ctx, bucketName)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("Fails with an invalid endpoint.", func(t *testing.T) {
		env := &config.Env{
			STORAGE_URI:                "endpoint-com-espaco invalido",
			STORAGE_ACCESS_KEY_ID:      accessKey,
			STORAGE_SECRECT_ACCESS_KEY: secretKey,
			STORAGE_USE_SSL:            false,
		}

		client, err := s3.NewS3Connection(env)
		require.Error(t, err)
		require.Nil(t, client)
	})
}
