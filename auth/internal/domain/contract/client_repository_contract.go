package contract

import (
	"auth/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type ClientRepositoryInterface interface {
	Create(ctx context.Context, model *entity.Client) error
	Delete(ctx context.Context, model *entity.Client) error
	Login(ctx context.Context, email string) (*entity.Client, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error)
}
