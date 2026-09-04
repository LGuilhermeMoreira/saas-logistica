package output

import "github.com/google/uuid"

type AddressOutput struct {
	ID           uuid.UUID
	Street       string
	Number       string
	Neighborhood string
	City         string
	ZipCode      string
}

type DeliveryOutput struct {
	ID       uuid.UUID
	To       AddressOutput
	From     AddressOutput
	Weight   float32
	Metadata map[string]any
	ClientID uuid.UUID
	DriverID *uuid.UUID
	Status   string
}
