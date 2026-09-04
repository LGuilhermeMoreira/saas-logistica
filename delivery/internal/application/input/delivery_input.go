package input

type AddressInput struct {
	Street       string
	Number       string
	Neighborhood string
	City         string
	ZipCode      string
}

type CreateDeliveryInput struct {
	To       AddressInput
	From     AddressInput
	Weight   float32
	ClientID string
	Metadata map[string]any
}

type FindByIDDeliveryInput struct {
	ID string
}

type AssignDeliveryToDriverInput struct {
	DeliveryID string
	DriverID   string
}
