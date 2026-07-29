package service

import (
	"auth/config"
	"auth/internal/domain/contract"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageService struct {
	client              *minio.Client
	bucketName          string
	signedUrlExpiration time.Duration
}

func NewStorageService(client *minio.Client, env *config.Env) contract.StorageService {
	return &StorageService{
		client:              client,
		bucketName:          env.STORAGE_BUCKET_NAME,
		signedUrlExpiration: env.STORAGE_SIGNED_URL_EXPIRATION,
	}
}

func (s *StorageService) UploadFile(
	ctx context.Context,
	filename, contentType string,
	reader io.Reader,
	size int64,
) (string, string, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}

	if _, err := s.client.PutObject(ctx, s.bucketName, filename, reader, size, opts); err != nil {
		return "", "", fmt.Errorf("storage: falha ao fazer upload de %q: %w", filename, err)
	}

	signedURL, err := s.signURL(ctx, filename)
	if err != nil {
		return "", "", err
	}

	return filename, signedURL, nil
}

func (s *StorageService) GetSignedUrl(ctx context.Context, filename string) (string, error) {
	return s.signURL(ctx, filename)
}

func (s *StorageService) DeleteFile(ctx context.Context, filename string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("storage: falha ao deletar %q: %w", filename, err)
	}
	return nil
}

func (s *StorageService) signURL(ctx context.Context, filename string) (string, error) {
	signed, err := s.client.PresignedGetObject(
		ctx, s.bucketName, filename, s.signedUrlExpiration, url.Values{},
	)
	if err != nil {
		return "", fmt.Errorf("storage: falha ao gerar URL assinada para %q: %w", filename, err)
	}
	return signed.String(), nil
}
