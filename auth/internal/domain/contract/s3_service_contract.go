package contract

import (
	"context"
	"io"
)

type StorageService interface {
	UploadFile(ctx context.Context, filename, contentType string, reader io.Reader, size int64) (string, string, error)
	GetSignedUrl(ctx context.Context, filename string) (string, error)
	DeleteFile(ctx context.Context, oldURL string) error
}
