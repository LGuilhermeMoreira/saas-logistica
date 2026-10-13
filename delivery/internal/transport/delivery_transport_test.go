package transport

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"delivery/internal/application/input"
	"delivery/internal/application/output"
	"delivery/internal/domain/errs"
	"delivery/pkg/authentication"

	deliveryv1 "proto/gen/delivery/v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockDeliveryUsecase struct {
	mock.Mock
}

func (m *mockDeliveryUsecase) FindByID(ctx context.Context, dto input.FindByIDDeliveryInput) (*output.DeliveryOutput, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*output.DeliveryOutput), args.Error(1)
}

func (m *mockDeliveryUsecase) AssignToDriver(ctx context.Context, dto input.AssignDeliveryToDriverInput) (*output.DeliveryOutput, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*output.DeliveryOutput), args.Error(1)
}

func (m *mockDeliveryUsecase) FindUnassociated(ctx context.Context) ([]output.DeliveryOutput, error) {
	args := m.Called(ctx)
	return args.Get(0).([]output.DeliveryOutput), args.Error(1)
}

func (m *mockDeliveryUsecase) Create(ctx context.Context, dto input.CreateDeliveryInput) (*output.DeliveryOutput, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*output.DeliveryOutput), args.Error(1)
}

func newTestTransport(uc *mockDeliveryUsecase) *DeliveryTransport {
	return NewDeliveryTransport(uc, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func newAddress() *deliveryv1.Address {
	return &deliveryv1.Address{Street: "Rua A", Number: "1", Neighborhood: "Centro", City: "Quixadá", ZipCode: "63900-000"}
}

func assertGRPCError(t *testing.T, err error, code codes.Code, msg string) {
	t.Helper()
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, code, st.Code())
	assert.Equal(t, msg, st.Message())
}

func TestDeliveryTransport_CreateDelivery_MissingAddress(t *testing.T) {
	uc := new(mockDeliveryUsecase)
	tr := newTestTransport(uc)
	ctx := authentication.ContextWithUserID(context.Background(), uuid.New())

	tests := []struct {
		name string
		req  *deliveryv1.CreateDeliveryRequest
	}{
		{"missing to", &deliveryv1.CreateDeliveryRequest{From: newAddress(), Weight: 1}},
		{"missing from", &deliveryv1.CreateDeliveryRequest{To: newAddress(), Weight: 1}},
		{"missing both", &deliveryv1.CreateDeliveryRequest{Weight: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tr.CreateDelivery(ctx, tt.req)

			assert.Nil(t, resp)
			assertGRPCError(t, err, codes.InvalidArgument, errs.ErrMissingAddress.Error())
		})
	}

	uc.AssertNotCalled(t, "Create")
}

func TestDeliveryTransport_CreateDelivery_RequiresAuthenticatedUser(t *testing.T) {
	uc := new(mockDeliveryUsecase)
	tr := newTestTransport(uc)

	resp, err := tr.CreateDelivery(context.Background(), &deliveryv1.CreateDeliveryRequest{To: newAddress(), From: newAddress(), Weight: 1})

	assert.Nil(t, resp)
	assertGRPCError(t, err, codes.Unauthenticated, errs.ErrUnauthenticated.Error())
	uc.AssertNotCalled(t, "Create")
}

func TestDeliveryTransport_CreateDelivery_UsesAuthenticatedUserAsClient(t *testing.T) {
	uc := new(mockDeliveryUsecase)
	tr := newTestTransport(uc)
	clientID := uuid.New()
	ctx := authentication.ContextWithUserID(context.Background(), clientID)

	uc.On("Create", ctx, mock.MatchedBy(func(dto input.CreateDeliveryInput) bool {
		return dto.ClientID == clientID.String()
	})).Return(&output.DeliveryOutput{ID: uuid.New(), ClientID: clientID, Status: "created"}, nil)

	resp, err := tr.CreateDelivery(ctx, &deliveryv1.CreateDeliveryRequest{To: newAddress(), From: newAddress(), Weight: 1})

	require.NoError(t, err)
	assert.Equal(t, clientID.String(), resp.GetDelivery().GetClientId())
	assert.Empty(t, resp.GetDelivery().GetDriverId())
	uc.AssertExpectations(t)
}

func TestDeliveryTransport_AssignToDriver_RequiresAuthenticatedUser(t *testing.T) {
	uc := new(mockDeliveryUsecase)
	tr := newTestTransport(uc)

	resp, err := tr.AssignToDriver(context.Background(), &deliveryv1.AssignToDriverRequest{DeliveryId: uuid.NewString()})

	assert.Nil(t, resp)
	assertGRPCError(t, err, codes.Unauthenticated, errs.ErrUnauthenticated.Error())
	uc.AssertNotCalled(t, "AssignToDriver")
}

func TestDeliveryTransport_AssignToDriver_UsesAuthenticatedUserAsDriver(t *testing.T) {
	uc := new(mockDeliveryUsecase)
	tr := newTestTransport(uc)
	driverID := uuid.New()
	deliveryID := uuid.New()
	ctx := authentication.ContextWithUserID(context.Background(), driverID)

	uc.On("AssignToDriver", ctx, input.AssignDeliveryToDriverInput{
		DeliveryID: deliveryID.String(),
		DriverID:   driverID.String(),
	}).Return(&output.DeliveryOutput{ID: deliveryID, DriverID: &driverID, Status: "in_transit"}, nil)

	resp, err := tr.AssignToDriver(ctx, &deliveryv1.AssignToDriverRequest{DeliveryId: deliveryID.String()})

	require.NoError(t, err)
	assert.Equal(t, driverID.String(), resp.GetDelivery().GetDriverId())
	uc.AssertExpectations(t)
}
