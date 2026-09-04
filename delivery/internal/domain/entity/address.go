package entity

import (
	"errors"
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
		return nil, errors.New("street is invalid")
	}

	if strings.TrimSpace(number) == "" {
		return nil, errors.New("number is invalid")
	}

	if strings.TrimSpace(neighborhood) == "" {
		return nil, errors.New("neighborhood is invalid")
	}

	if strings.TrimSpace(city) == "" {
		return nil, errors.New("city is invalid")
	}

	if strings.TrimSpace(zipCode) == "" {
		return nil, errors.New("zip code is invalid")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
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
