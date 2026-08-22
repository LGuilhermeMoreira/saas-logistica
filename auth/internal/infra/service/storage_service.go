package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"time"

	"auth/config"
	"auth/internal/domain/contract"
	"auth/pkg/logger"

	"github.com/minio/minio-go/v7"
)

type StorageService struct {
	client              *minio.Client
	bucketName          string
	signedUrlExpiration time.Duration
	log                 *slog.Logger
}

func NewStorageService(client *minio.Client, env *config.Env, log *slog.Logger) contract.StorageServiceInterface {
	return &StorageService{
		client:              client,
		bucketName:          env.STORAGE_BUCKET_NAME,
		signedUrlExpiration: env.STORAGE_SIGNED_URL_EXPIRATION,
		log:                 log,
	}
}

func (s *StorageService) UploadFile(
	ctx context.Context,
	filename, contentType string,
	reader io.Reader,
	size int64,
) (string, string, error) {
	reqID := logger.ExtractRequestID(ctx)
	s.log.Debug("starting file upload to storage", slog.String("request_id", reqID), slog.String("filename", filename), slog.Int64("size", size))

	opts := minio.PutObjectOptions{ContentType: contentType}

	if _, err := s.client.PutObject(ctx, s.bucketName, filename, reader, size, opts); err != nil {
		s.log.Error("failed to upload file to minio", slog.String("request_id", reqID), slog.String("filename", filename), slog.Any("error", err))
		return "", "", fmt.Errorf("storage: falha ao fazer upload de %q: %w", filename, err)
	}

	signedURL, err := s.signURL(ctx, filename)
	if err != nil {
		return "", "", err
	}

	s.log.Info("file uploaded successfully", slog.String("request_id", reqID), slog.String("filename", filename))

	return filename, signedURL, nil
}

func (s *StorageService) GetSignedUrl(ctx context.Context, filename string) (string, error) {
	reqID := logger.ExtractRequestID(ctx)
	s.log.Debug("requesting signed url", slog.String("request_id", reqID), slog.String("filename", filename))

	return s.signURL(ctx, filename)
}

func (s *StorageService) DeleteFile(ctx context.Context, filename string) error {
	reqID := logger.ExtractRequestID(ctx)
	s.log.Debug("starting file deletion from storage", slog.String("request_id", reqID), slog.String("filename", filename))

	err := s.client.RemoveObject(ctx, s.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		s.log.Error("failed to delete file from minio", slog.String("request_id", reqID), slog.String("filename", filename), slog.Any("error", err))
		return fmt.Errorf("storage: falha ao deletar %q: %w", filename, err)
	}

	s.log.Info("file deleted successfully", slog.String("request_id", reqID), slog.String("filename", filename))
	return nil
}

func (s *StorageService) signURL(ctx context.Context, filename string) (string, error) {
	reqID := logger.ExtractRequestID(ctx)

	signed, err := s.client.PresignedGetObject(
		ctx, s.bucketName, filename, s.signedUrlExpiration, url.Values{},
	)
	if err != nil {
		s.log.Error("failed to generate signed url", slog.String("request_id", reqID), slog.String("filename", filename), slog.Any("error", err))
		return "", fmt.Errorf("storage: falha ao gerar URL assinada para %q: %w", filename, err)
	}

	return signed.String(), nil
}
