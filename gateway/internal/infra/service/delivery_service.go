package service

import (
	"context"

	"gateway/internal/domain/contract"
	deliveryv1 "proto/gen/delivery/v1"
)

type DeliveryService struct {
	deliveryClient deliveryv1.DeliveryServiceClient
}

func NewDeliveryService(deliveryClient deliveryv1.DeliveryServiceClient) contract.DeliveryService {
	return &DeliveryService{deliveryClient: deliveryClient}
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, req *deliveryv1.CreateDeliveryRequest) (*deliveryv1.CreateDeliveryResponse, error) {
	return s.deliveryClient.CreateDelivery(ctx, req)
}

func (s *DeliveryService) FindByID(ctx context.Context, id string) (*deliveryv1.FindByIDResponse, error) {
	return s.deliveryClient.FindByID(ctx, &deliveryv1.FindByIDRequest{Id: id})
}

func (s *DeliveryService) AssignToDriver(ctx context.Context, deliveryID string) (*deliveryv1.AssignToDriverResponse, error) {
	return s.deliveryClient.AssignToDriver(ctx, &deliveryv1.AssignToDriverRequest{
		DeliveryId: deliveryID,
	})
}

func (s *DeliveryService) FindUnassociated(ctx context.Context) (*deliveryv1.FindUnassociatedResponse, error) {
	return s.deliveryClient.FindUnassociated(ctx, &deliveryv1.FindUnassociatedRequest{})
}
