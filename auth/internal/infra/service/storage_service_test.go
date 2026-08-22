package service_test

import (
	"auth/config"
	"auth/internal/infra/service"
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	miniotc "github.com/testcontainers/testcontainers-go/modules/minio"
)

func setupMinioContainer(t *testing.T) (*minio.Client, *config.Env, func()) {
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

	endpoint, err := minioContainer.ConnectionString(ctx)
	require.NoError(t, err)

	username := minioContainer.Username
	password := minioContainer.Password

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false, // Since it's a local container, no TLS
	})
	require.NoError(t, err, "failed to initialize minio client")

	bucketName := "test-bucket"
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	require.NoError(t, err, "failed to create test bucket")

	env := &config.Env{
		STORAGE_BUCKET_NAME:           bucketName,
		STORAGE_SIGNED_URL_EXPIRATION: 1 * time.Hour,
	}

	cleanup := func() {
		err := minioContainer.Terminate(ctx)
		require.NoError(t, err, "failed to terminate minio container")
	}

	return client, env, cleanup
}

func TestStorageService(t *testing.T) {
	client, env, cleanup := setupMinioContainer(t)
	defer cleanup()

	svc := service.NewStorageService(client, env, slog.Default())
	ctx := context.Background()

	t.Run("success uploading file", func(t *testing.T) {
		content := []byte("hello minio testcontainers")
		reader := bytes.NewReader(content)
		filename := "test-upload.txt"
		contentType := "text/plain"

		name, signedURL, err := svc.UploadFile(ctx, filename, contentType, reader, int64(len(content)))

		assert.NoError(t, err)
		assert.Equal(t, filename, name)
		assert.NotEmpty(t, signedURL)

		assert.Contains(t, signedURL, filename)
		assert.Contains(t, signedURL, "X-Amz-Signature")
	})

	t.Run("success getting signed url", func(t *testing.T) {
		filename := "test-signed-url.png"

		signedURL, err := svc.GetSignedUrl(ctx, filename)

		assert.NoError(t, err)
		assert.NotEmpty(t, signedURL)
		assert.Contains(t, signedURL, filename)
		assert.Contains(t, signedURL, "X-Amz-Credential")
	})

	t.Run("success deleting file", func(t *testing.T) {
		filename := "test-delete.pdf"
		content := []byte("dummy pdf content")
		reader := bytes.NewReader(content)

		_, _, err := svc.UploadFile(ctx, filename, "application/pdf", reader, int64(len(content)))
		require.NoError(t, err)

		_, err = client.StatObject(ctx, env.STORAGE_BUCKET_NAME, filename, minio.StatObjectOptions{})
		require.NoError(t, err, "file should exist before deletion")

		err = svc.DeleteFile(ctx, filename)
		assert.NoError(t, err)

		_, err = client.StatObject(ctx, env.STORAGE_BUCKET_NAME, filename, minio.StatObjectOptions{})
		assert.Error(t, err)

		errResp, ok := err.(minio.ErrorResponse)
		require.True(t, ok)
		assert.Equal(t, "NoSuchKey", errResp.Code)
	})
}
