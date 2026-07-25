package contract

import (
	"auth/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, model *entity.User) error
	Delete(ctx context.Context, model *entity.User) error
	Login(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
