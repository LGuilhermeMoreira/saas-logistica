package transport

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"delivery/internal/domain/errs"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
		wantMsg  string
	}{
		{"validation", errs.ErrInvalidWeight, codes.InvalidArgument, errs.ErrInvalidWeight.Error()},
		{"not found", errs.ErrDeliveryNotFound, codes.NotFound, errs.ErrDeliveryNotFound.Error()},
		{"already assigned", errs.ErrDeliveryAlreadyAssigned, codes.FailedPrecondition, errs.ErrDeliveryAlreadyAssigned.Error()},
		{"concurrency conflict", errs.ErrConcurrencyConflict, codes.Aborted, errs.ErrConcurrencyConflict.Error()},
		{"context deadline", fmt.Errorf("failed to create delivery: %w", context.DeadlineExceeded), codes.DeadlineExceeded, context.DeadlineExceeded.Error()},
		{"unknown error is hidden", fmt.Errorf("failed to create delivery: %w", errors.New("server selection error: mongo:27017")), codes.Internal, errs.ErrInternal.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, ok := status.FromError(toGRPCError(tt.err))

			assert.True(t, ok)
			assert.Equal(t, tt.wantCode, st.Code())
			assert.Equal(t, tt.wantMsg, st.Message())
		})
	}

	assert.NoError(t, toGRPCError(nil))
}
