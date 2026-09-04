package entity_test

import (
	"delivery/internal/domain/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDelivery(t *testing.T) {
	to := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua A",
		Number:       "100",
		Neighborhood: "Centro",
		City:         "Quixadá",
		ZipCode:      "63900-000",
	}

	from := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua B",
		Number:       "200",
		Neighborhood: "Centro",
		City:         "Quixadá",
		ZipCode:      "63900-000",
	}

	clientID := uuid.New()

	tests := []struct {
		name     string
		weight   float32
		clientID uuid.UUID
		wantErr  string
	}{
		{
			name:     "should create delivery",
			weight:   10,
			clientID: clientID,
		},
		{
			name:     "should reject zero weight",
			weight:   0,
			clientID: clientID,
			wantErr:  "weight is invalid",
		},
		{
			name:     "should reject negative weight",
			weight:   -1,
			clientID: clientID,
			wantErr:  "weight is invalid",
		},
		{
			name:     "should reject nil client id",
			weight:   10,
			clientID: uuid.Nil,
			wantErr:  "client id is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := map[string]any{
				"fragile": true,
			}

			delivery, err := entity.NewDelivery(
				to,
				from,
				tt.weight,
				tt.clientID,
				metadata,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Nil(t, delivery)
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, delivery)

			assert.NotEqual(t, uuid.Nil, delivery.ID)

			assert.Equal(t, to, delivery.To)
			assert.Equal(t, from, delivery.From)
			assert.Equal(t, tt.weight, delivery.Weight)
			assert.Equal(t, tt.clientID, delivery.ClientID)
			assert.Equal(t, metadata, delivery.Metadata)

			assert.Nil(t, delivery.DriverID)

			assert.Equal(t, "created", string(delivery.Status))
		})
	}
}

func TestDelivery_AssingDriver(t *testing.T) {
	to := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua A",
		Number:       "100",
		Neighborhood: "Centro",
		City:         "Quixadá",
		ZipCode:      "63900-000",
	}

	from := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua B",
		Number:       "200",
		Neighborhood: "Centro",
		City:         "Quixadá",
		ZipCode:      "63900-000",
	}

	clientID := uuid.New()

	t.Run("should assign driver", func(t *testing.T) {
		delivery, err := entity.NewDelivery(
			to,
			from,
			10,
			clientID,
			nil,
		)
		require.NoError(t, err)

		driverID := uuid.New()

		err = delivery.AssingDriver(driverID)

		require.NoError(t, err)

		require.NotNil(t, delivery.DriverID)
		assert.Equal(t, driverID, *delivery.DriverID)
		assert.Equal(t, "in_transit", string(delivery.Status))
	})

	t.Run("should reject nil driver id", func(t *testing.T) {
		delivery, err := entity.NewDelivery(
			to,
			from,
			10,
			clientID,
			nil,
		)
		require.NoError(t, err)

		err = delivery.AssingDriver(uuid.Nil)

		require.Error(t, err)
		assert.EqualError(t, err, "driver id is invalid")

		assert.Nil(t, delivery.DriverID)
		assert.Equal(t, "created", string(delivery.Status))
	})

	t.Run("should reject assigning another driver", func(t *testing.T) {
		delivery, err := entity.NewDelivery(
			to,
			from,
			10,
			clientID,
			nil,
		)
		require.NoError(t, err)

		firstDriverID := uuid.New()
		secondDriverID := uuid.New()

		err = delivery.AssingDriver(firstDriverID)
		require.NoError(t, err)

		err = delivery.AssingDriver(secondDriverID)

		require.Error(t, err)
		assert.EqualError(
			t,
			err,
			"delivery already assigned by driver",
		)

		require.NotNil(t, delivery.DriverID)
		assert.Equal(t, firstDriverID, *delivery.DriverID)
		assert.Equal(t, "in_transit", string(delivery.Status))
	})
}
