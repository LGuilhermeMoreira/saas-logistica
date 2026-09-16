package delivery

type AddressInput struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	ZipCode      string `json:"zip_code"`
}

type CreateDeliveryInput struct {
	To       AddressInput   `json:"to"`
	From     AddressInput   `json:"from"`
	Weight   float32        `json:"weight"`
	ClientID string         `json:"client_id"`
	Metadata map[string]any `json:"metadata"`
}

type AssignToDriverInput struct {
	DeliveryID string `json:"delivery_id"`
	DriverID   string `json:"driver_id"`
}

type AddressOutput struct {
	ID           string `json:"id"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	ZipCode      string `json:"zip_code"`
}

type DeliveryOutput struct {
	ID       string         `json:"id"`
	To       AddressOutput  `json:"to"`
	From     AddressOutput  `json:"from"`
	Weight   float32        `json:"weight"`
	Metadata map[string]any `json:"metadata"`
	ClientID string         `json:"client_id"`
	DriverID string         `json:"driver_id"`
	Status   string         `json:"status"`
}

type DeliveryResponse struct {
	Delivery DeliveryOutput `json:"delivery"`
}

type DeliveriesResponse struct {
	Deliveries []DeliveryOutput `json:"deliveries"`
}
