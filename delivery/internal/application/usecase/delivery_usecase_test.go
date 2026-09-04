package usecase_test

import (
	"context"
	"errors"
	"testing"

	"delivery/internal/application/input"
	"delivery/internal/application/usecase"
	"delivery/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDeliveryRepository struct {
	mock.Mock
}

func (m *MockDeliveryRepository) Create(ctx context.Context, model *entity.Delivery) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockDeliveryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Delivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Delivery), args.Error(1)
}

func (m *MockDeliveryRepository) AssingToDriver(ctx context.Context, model *entity.Delivery) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockDeliveryRepository) FindUnassociatedDeliveries(ctx context.Context) ([]entity.Delivery, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Delivery), args.Error(1)
}

func newDelivery(t *testing.T) *entity.Delivery {
	t.Helper()

	to := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua A",
		Number:       "100",
		Neighborhood: "Centro",
		City:         "Quixada",
		ZipCode:      "63900-000",
	}
	from := entity.Address{
		ID:           uuid.New(),
		Street:       "Rua B",
		Number:       "200",
		Neighborhood: "Centro",
		City:         "Quixada",
		ZipCode:      "63900-000",
	}

	delivery, err := entity.NewDelivery(to, from, 10, uuid.New(), map[string]any{"fragile": true})
	require.NoError(t, err)
	return delivery
}

func validCreateInput(clientID uuid.UUID) input.CreateDeliveryInput {
	return input.CreateDeliveryInput{
		To: input.AddressInput{
			Street:       "Rua A",
			Number:       "100",
			Neighborhood: "Centro",
			City:         "Quixada",
			ZipCode:      "63900-000",
		},
		From: input.AddressInput{
			Street:       "Rua B",
			Number:       "200",
			Neighborhood: "Centro",
			City:         "Quixada",
			ZipCode:      "63900-000",
		},
		Weight:   10,
		ClientID: clientID.String(),
		Metadata: map[string]any{"fragile": true},
	}
}

func TestDeliveryUsecase_FindByID(t *testing.T) {
	ctx := context.Background()
	delivery := newDelivery(t)
	repository := new(MockDeliveryRepository)
	repository.On("FindByID", ctx, delivery.ID).Return(delivery, nil)

	result, err := usecase.NewDeliveryUsecase(repository).FindByID(ctx, input.FindByIDDeliveryInput{ID: delivery.ID.String()})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, delivery.ID, result.ID)
	assert.Equal(t, delivery.To.Street, result.To.Street)
	assert.Equal(t, delivery.From.City, result.From.City)
	assert.Equal(t, delivery.Weight, result.Weight)
	assert.Equal(t, delivery.ClientID, result.ClientID)
	assert.Equal(t, "created", result.Status)
	repository.AssertExpectations(t)
}

func TestDeliveryUsecase_FindByID_InvalidID(t *testing.T) {
	repository := new(MockDeliveryRepository)

	result, err := usecase.NewDeliveryUsecase(repository).FindByID(context.Background(), input.FindByIDDeliveryInput{ID: "invalid-id"})

	assert.Error(t, err)
	assert.Nil(t, result)
	repository.AssertNotCalled(t, "FindByID")
}

func TestDeliveryUsecase_AssignToDriver(t *testing.T) {
	ctx := context.Background()
	delivery := newDelivery(t)
	driverID := uuid.New()
	repository := new(MockDeliveryRepository)
	repository.On("FindByID", ctx, delivery.ID).Return(delivery, nil)
	repository.On("AssingToDriver", ctx, mock.MatchedBy(func(model *entity.Delivery) bool {
		return model.ID == delivery.ID && model.DriverID != nil && *model.DriverID == driverID
	})).Return(nil)

	result, err := usecase.NewDeliveryUsecase(repository).AssignToDriver(ctx, input.AssignDeliveryToDriverInput{
		DeliveryID: delivery.ID.String(),
		DriverID:   driverID.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.DriverID)
	assert.Equal(t, driverID, *result.DriverID)
	assert.Equal(t, "in_transit", result.Status)
	repository.AssertExpectations(t)
}

func TestDeliveryUsecase_AssignToDriver_RejectsAlreadyAssignedDelivery(t *testing.T) {
	ctx := context.Background()
	delivery := newDelivery(t)
	firstDriverID := uuid.New()
	require.NoError(t, delivery.AssingDriver(firstDriverID))
	repository := new(MockDeliveryRepository)
	repository.On("FindByID", ctx, delivery.ID).Return(delivery, nil)

	result, err := usecase.NewDeliveryUsecase(repository).AssignToDriver(ctx, input.AssignDeliveryToDriverInput{
		DeliveryID: delivery.ID.String(),
		DriverID:   uuid.New().String(),
	})

	assert.EqualError(t, err, "delivery already assigned by driver")
	assert.Nil(t, result)
	repository.AssertNotCalled(t, "AssingToDriver")
}

func TestDeliveryUsecase_FindUnassociated(t *testing.T) {
	ctx := context.Background()
	deliveryOne := newDelivery(t)
	deliveryTwo := newDelivery(t)
	repository := new(MockDeliveryRepository)
	repository.On("FindUnassociatedDeliveries", ctx).Return([]entity.Delivery{*deliveryOne, *deliveryTwo}, nil)

	result, err := usecase.NewDeliveryUsecase(repository).FindUnassociated(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, deliveryOne.ID, result[0].ID)
	assert.Equal(t, deliveryTwo.ID, result[1].ID)
	assert.Equal(t, deliveryOne.ClientID, result[0].ClientID)
	repository.AssertExpectations(t)
}

func TestDeliveryUsecase_FindUnassociated_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("database error")
	repository := new(MockDeliveryRepository)
	repository.On("FindUnassociatedDeliveries", ctx).Return(nil, expectedErr)

	result, err := usecase.NewDeliveryUsecase(repository).FindUnassociated(ctx)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	repository.AssertExpectations(t)
}

func TestDeliveryUsecase_Create(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New()
	repository := new(MockDeliveryRepository)
	repository.On("Create", ctx, mock.MatchedBy(func(model *entity.Delivery) bool {
		return model.ClientID == clientID && model.Weight == 10 && model.Status == "created"
	})).Return(nil)

	result, err := usecase.NewDeliveryUsecase(repository).Create(ctx, validCreateInput(clientID))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, clientID, result.ClientID)
	assert.Equal(t, "Rua A", result.To.Street)
	assert.Equal(t, "Rua B", result.From.Street)
	assert.Equal(t, map[string]any{"fragile": true}, result.Metadata)
	assert.Equal(t, "created", result.Status)
	repository.AssertExpectations(t)
}

func TestDeliveryUsecase_Create_RejectsInvalidInput(t *testing.T) {
	repository := new(MockDeliveryRepository)
	dto := validCreateInput(uuid.New())
	dto.Weight = 0

	result, err := usecase.NewDeliveryUsecase(repository).Create(context.Background(), dto)

	assert.EqualError(t, err, "weight is invalid")
	assert.Nil(t, result)
	repository.AssertNotCalled(t, "Create")
}
