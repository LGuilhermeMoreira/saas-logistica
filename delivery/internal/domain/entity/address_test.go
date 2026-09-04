package entity_test

import (
	"delivery/internal/domain/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name         string
		street       string
		number       string
		neighborhood string
		city         string
		zipCode      string
		wantErr      string
	}{
		{
			name:         "should create address",
			street:       "Rua A",
			number:       "123",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "63900-000",
		},
		{
			name:         "should reject empty street",
			street:       "",
			number:       "123",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "street is invalid",
		},
		{
			name:         "should reject whitespace street",
			street:       "   ",
			number:       "123",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "street is invalid",
		},
		{
			name:         "should reject empty number",
			street:       "Rua A",
			number:       "",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "number is invalid",
		},
		{
			name:         "should reject whitespace number",
			street:       "Rua A",
			number:       "   ",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "number is invalid",
		},
		{
			name:         "should reject empty neighborhood",
			street:       "Rua A",
			number:       "123",
			neighborhood: "",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "neighborhood is invalid",
		},
		{
			name:         "should reject whitespace neighborhood",
			street:       "Rua A",
			number:       "123",
			neighborhood: "   ",
			city:         "Quixadá",
			zipCode:      "63900-000",
			wantErr:      "neighborhood is invalid",
		},
		{
			name:         "should reject empty city",
			street:       "Rua A",
			number:       "123",
			neighborhood: "Centro",
			city:         "",
			zipCode:      "63900-000",
			wantErr:      "city is invalid",
		},
		{
			name:         "should reject whitespace city",
			street:       "Rua A",
			number:       "123",
			neighborhood: "Centro",
			city:         "   ",
			zipCode:      "63900-000",
			wantErr:      "city is invalid",
		},
		{
			name:         "should reject empty zip code",
			street:       "Rua A",
			number:       "123",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "",
			wantErr:      "zip code is invalid",
		},
		{
			name:         "should reject whitespace zip code",
			street:       "Rua A",
			number:       "123",
			neighborhood: "Centro",
			city:         "Quixadá",
			zipCode:      "   ",
			wantErr:      "zip code is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := entity.NewAddress(
				tt.street,
				tt.number,
				tt.neighborhood,
				tt.city,
				tt.zipCode,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Nil(t, address)
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, address)

			assert.NotEqual(t, uuid.Nil, address.ID)
			assert.Equal(t, tt.street, address.Street)
			assert.Equal(t, tt.number, address.Number)
			assert.Equal(t, tt.neighborhood, address.Neighborhood)
			assert.Equal(t, tt.city, address.City)
			assert.Equal(t, tt.zipCode, address.ZipCode)
		})
	}
}
