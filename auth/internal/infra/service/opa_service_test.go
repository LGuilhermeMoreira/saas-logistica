package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"auth/internal/domain/entity"
	"auth/internal/infra/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockRoleRepository struct {
	FindAllFunc func(ctx context.Context) ([]entity.Role, error)
}

func (m *mockRoleRepository) Create(ctx context.Context, model *entity.Role) error {
	return nil
}

func (m *mockRoleRepository) Delete(ctx context.Context, model *entity.Role) error {
	return nil
}

func (m *mockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	return nil, nil
}

func (m *mockRoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	return m.FindAllFunc(ctx)
}

type mockStorageService struct {
	UploadFileFunc func(ctx context.Context, filename, contentType string, reader io.Reader, size int64) (string, string, error)
}

func (m *mockStorageService) UploadFile(ctx context.Context, filename, contentType string, reader io.Reader, size int64) (string, string, error) {
	return m.UploadFileFunc(ctx, filename, contentType, reader, size)
}

func (m *mockStorageService) GetSignedUrl(ctx context.Context, filename string) (string, error) {
	return "", nil
}

func (m *mockStorageService) DeleteFile(ctx context.Context, filename string) error {
	return nil
}

func TestOPASyncService_SyncPolicies(t *testing.T) {
	ctx := context.Background()

	t.Run("success syncing policies", func(t *testing.T) {
		repo := &mockRoleRepository{
			FindAllFunc: func(ctx context.Context) ([]entity.Role, error) {
				return []entity.Role{
					{
						ID:   uuid.New(),
						Name: "Admin",
					},
				}, nil
			},
		}

		storage := &mockStorageService{
			UploadFileFunc: func(ctx context.Context, filename, contentType string, reader io.Reader, size int64) (string, string, error) {
				assert.Equal(t, "bundle.tar.gz", filename)
				assert.Equal(t, "application/gzip", contentType)
				assert.Greater(t, size, int64(0))

				buf := new(bytes.Buffer)
				_, err := io.Copy(buf, reader)
				assert.NoError(t, err)
				assert.Greater(t, buf.Len(), 0)

				return filename, "http://signed-url.com", nil
			},
		}

		svc := service.NewOPASyncService(repo, storage)
		err := svc.SyncPolicies(ctx)

		assert.NoError(t, err)
	})

	t.Run("error fetching roles", func(t *testing.T) {
		repo := &mockRoleRepository{
			FindAllFunc: func(ctx context.Context) ([]entity.Role, error) {
				return nil, errors.New("database connection failed")
			},
		}

		storage := &mockStorageService{}

		svc := service.NewOPASyncService(repo, storage)
		err := svc.SyncPolicies(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao buscar roles do banco")
	})

	t.Run("error uploading bundle", func(t *testing.T) {
		repo := &mockRoleRepository{
			FindAllFunc: func(ctx context.Context) ([]entity.Role, error) {
				return []entity.Role{}, nil
			},
		}

		storage := &mockStorageService{
			UploadFileFunc: func(ctx context.Context, filename, contentType string, reader io.Reader, size int64) (string, string, error) {
				return "", "", errors.New("storage unavailable")
			},
		}

		svc := service.NewOPASyncService(repo, storage)
		err := svc.SyncPolicies(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao fazer upload do bundle para o storage")
	})
}
