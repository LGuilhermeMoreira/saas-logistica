package transport

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"delivery/internal/domain/errs"

	deliveryv1 "proto/gen/delivery/v1"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeliveryTransport_CreateDelivery_MissingAddress(t *testing.T) {
	// nil usecase: the request must be rejected before reaching it
	tr := NewDeliveryTransport(nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	address := &deliveryv1.Address{Street: "Rua A", Number: "1", Neighborhood: "Centro", City: "Quixadá", ZipCode: "63900-000"}

	tests := []struct {
		name string
		req  *deliveryv1.CreateDeliveryRequest
	}{
		{"missing to", &deliveryv1.CreateDeliveryRequest{From: address, Weight: 1}},
		{"missing from", &deliveryv1.CreateDeliveryRequest{To: address, Weight: 1}},
		{"missing both", &deliveryv1.CreateDeliveryRequest{Weight: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tr.CreateDelivery(context.Background(), tt.req)

			assert.Nil(t, resp)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
			assert.Equal(t, errs.ErrMissingAddress.Error(), st.Message())
		})
	}
}
