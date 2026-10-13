package entity

import (
	"delivery/internal/domain/errs"
	"fmt"

	"github.com/google/uuid"
)

type deliveryStatus string

const (
	DeliveryStatusCreated   deliveryStatus = "created"
	DeliveryStatusInTransit deliveryStatus = "in_transit"
	DeliveryStatusDelivered deliveryStatus = "delivered"
	DeliveryStatusCancelled deliveryStatus = "cancelled"
)

type Delivery struct {
	ID       uuid.UUID      `bson:"_id"`
	To       Address        `bson:"to"`
	From     Address        `bson:"from"`
	Weight   float32        `bson:"weight"`
	Metadata map[string]any `bson:"metadata,omitempty"`
	ClientID uuid.UUID      `bson:"client_id"`
	DriverID *uuid.UUID     `bson:"driver_id,omitempty"`
	Status   deliveryStatus `bson:"status"`
}

func NewDelivery(to, from Address, weight float32, clientID uuid.UUID, metadata map[string]any) (*Delivery, error) {
	if weight <= 0 {
		return nil, errs.ErrInvalidWeight
	}

	if clientID == uuid.Nil {
		return nil, errs.ErrInvalidClientID
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate delivery id: %w", err)
	}

	return &Delivery{
		ID:       id,
		To:       to,
		From:     from,
		Weight:   weight,
		Metadata: metadata,
		ClientID: clientID,
		DriverID: nil,
		Status:   DeliveryStatusCreated,
	}, nil
}

func (d *Delivery) AssingDriver(id uuid.UUID) error {
	if id == uuid.Nil {
		return errs.ErrInvalidDriverID
	}
	if d.DriverID != nil {
		return errs.ErrDeliveryAlreadyAssigned
	}
	d.DriverID = &id
	d.Status = DeliveryStatusInTransit
	return nil
}
