package httperror_test

import (
	"errors"
	"net/http"
	"testing"

	"gateway/internal/transport/api/httperror"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromGRPC(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"invalid argument", status.Error(codes.InvalidArgument, "weight is invalid"), http.StatusBadRequest, "weight is invalid"},
		{"unauthenticated", status.Error(codes.Unauthenticated, "invalid credentials"), http.StatusUnauthorized, "invalid credentials"},
		{"permission denied", status.Error(codes.PermissionDenied, "access denied by OPA"), http.StatusForbidden, "access denied by OPA"},
		{"not found", status.Error(codes.NotFound, "delivery not found"), http.StatusNotFound, "delivery not found"},
		{"already exists", status.Error(codes.AlreadyExists, "user already exists"), http.StatusConflict, "user already exists"},
		{"failed precondition", status.Error(codes.FailedPrecondition, "delivery already assigned by driver"), http.StatusConflict, "delivery already assigned by driver"},
		{"aborted", status.Error(codes.Aborted, "concurrency conflict"), http.StatusConflict, "concurrency conflict"},
		{"unavailable hides connection details", status.Error(codes.Unavailable, "connection error: dial tcp 10.0.0.5:50051: connect: connection refused"), http.StatusServiceUnavailable, "Service Unavailable"},
		{"deadline exceeded", status.Error(codes.DeadlineExceeded, "context deadline exceeded"), http.StatusGatewayTimeout, "Gateway Timeout"},
		{"internal hides message", status.Error(codes.Internal, "internal server error"), http.StatusInternalServerError, "Internal Server Error"},
		{"unknown code hides message", status.Error(codes.Unknown, "pq: relation users does not exist"), http.StatusInternalServerError, "Internal Server Error"},
		{"non gRPC error", errors.New("boom"), http.StatusInternalServerError, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotMsg := httperror.FromGRPC(tt.err)

			assert.Equal(t, tt.wantStatus, gotStatus)
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}
