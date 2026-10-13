package transport

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"auth/internal/domain/errs"

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
		{"validation", errs.ErrInvalidEmail, codes.InvalidArgument, errs.ErrInvalidEmail.Error()},
		{"invalid credentials", errs.ErrInvalidCredentials, codes.Unauthenticated, errs.ErrInvalidCredentials.Error()},
		{"wrapped not found keeps only sentinel message", fmt.Errorf("failed to find User by id: %w", errs.ErrUserNotFound), codes.NotFound, errs.ErrUserNotFound.Error()},
		{"already exists", fmt.Errorf("failed to persist User in database: %w", errs.ErrUserAlreadyExists), codes.AlreadyExists, errs.ErrUserAlreadyExists.Error()},
		{"role in use", errs.ErrRoleInUse, codes.FailedPrecondition, errs.ErrRoleInUse.Error()},
		{"context canceled", fmt.Errorf("query: %w", context.Canceled), codes.Canceled, context.Canceled.Error()},
		{"unknown error is hidden", errors.New("pq: connection refused at 10.0.0.1:5432"), codes.Internal, errs.ErrInternal.Error()},
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
