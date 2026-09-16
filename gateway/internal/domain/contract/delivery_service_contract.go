package contract

import (
	"context"

	deliveryv1 "proto/gen/delivery/v1"
)

type DeliveryService interface {
	CreateDelivery(context.Context, *deliveryv1.CreateDeliveryRequest) (*deliveryv1.CreateDeliveryResponse, error)
	FindByID(context.Context, string) (*deliveryv1.FindByIDResponse, error)
	AssignToDriver(context.Context, string, string) (*deliveryv1.AssignToDriverResponse, error)
	FindUnassociated(context.Context) (*deliveryv1.FindUnassociatedResponse, error)
}
