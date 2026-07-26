package contract

import (
	"auth/internal/domain/entity"
	"context"

	"github.com/google/uuid"
)

type RoleRepositoryInterface interface {
	Create(ctx context.Context, model *entity.Role) error
	Delete(ctx context.Context, model *entity.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
}
