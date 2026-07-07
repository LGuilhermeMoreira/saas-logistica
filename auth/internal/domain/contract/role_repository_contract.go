package contract

import (
	"auth/internal/domain/entity"
	"context"
)

type RoleRepositoryInterface interface {
	Create(ctx context.Context, model *entity.Role) error
	Delete(ctx context.Context, model *entity.Role) error
}
