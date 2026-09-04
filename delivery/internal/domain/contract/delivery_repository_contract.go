package contract

import (
	"context"
	"delivery/internal/domain/entity"

	"github.com/google/uuid"
)

type DeliveryRepositoryInterface interface {
	Create(ctx context.Context, model *entity.Delivery) error
	FindByID(ctx context.Context, ID uuid.UUID) (*entity.Delivery, error)
	AssingToDriver(ctx context.Context, model *entity.Delivery) error
	FindUnassociatedDeliveries(ctx context.Context) ([]entity.Delivery, error)
}
