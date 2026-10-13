package entity

import (
	"delivery/internal/domain/errs"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Address struct {
	ID           uuid.UUID `bson:"_id"`
	Street       string    `bson:"street"`
	Number       string    `bson:"number"`
	Neighborhood string    `bson:"neighborhood"`
	City         string    `bson:"city"`
	ZipCode      string    `bson:"zip_code"`
}

func NewAddress(street, number, neighborhood, city, zipCode string) (*Address, error) {
	if strings.TrimSpace(street) == "" {
		return nil, errs.ErrInvalidStreet
	}

	if strings.TrimSpace(number) == "" {
		return nil, errs.ErrInvalidNumber
	}

	if strings.TrimSpace(neighborhood) == "" {
		return nil, errs.ErrInvalidNeighborhood
	}

	if strings.TrimSpace(city) == "" {
		return nil, errs.ErrInvalidCity
	}

	if strings.TrimSpace(zipCode) == "" {
		return nil, errs.ErrInvalidZipCode
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate address id: %w", err)
	}

	return &Address{
		ID:           id,
		Street:       street,
		Number:       number,
		Neighborhood: neighborhood,
		City:         city,
		ZipCode:      zipCode,
	}, nil
}
