package grpcmiddleware_test

import (
	"context"
	"testing"

	"delivery/pkg/authentication"
	"delivery/pkg/authorization"
	grpcmiddleware "delivery/pkg/grpc_middleware"

	jwt_go "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeJWT struct {
	claims jwt_go.MapClaims
}

func (f fakeJWT) Validate(string) error { return nil }

func (f fakeJWT) ExtractClaims(string) (jwt_go.MapClaims, error) { return f.claims, nil }

type allowOPA struct{}

func (allowOPA) Validate(authorization.OPAInput) error { return nil }

func incomingContext() context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "Bearer token",
		"http-path":     "/delivery/create",
		"http-method":   "POST",
	}))
}

func TestValidateCredentials_PutsTokenUserIDInContext(t *testing.T) {
	userID := uuid.New()
	interceptor := grpcmiddleware.NewGRPCMiddleware(
		fakeJWT{claims: jwt_go.MapClaims{"data": map[string]any{"id": userID.String(), "role_name": "client"}}},
		allowOPA{},
	).ValidateCredentials()

	var got uuid.UUID
	handler := func(ctx context.Context, req any) (any, error) {
		got, _ = authentication.UserIDFromContext(ctx)
		return nil, nil
	}

	_, err := interceptor(incomingContext(), nil, &grpc.UnaryServerInfo{FullMethod: "/delivery.v1.DeliveryService/CreateDelivery"}, handler)

	require.NoError(t, err)
	assert.Equal(t, userID, got)
}

func TestValidateCredentials_RejectsTokenWithoutUserID(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
	}{
		{"missing id", map[string]any{"role_name": "client"}},
		{"invalid id", map[string]any{"id": "not-a-uuid", "role_name": "client"}},
		{"nil id", map[string]any{"id": uuid.Nil.String(), "role_name": "client"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := grpcmiddleware.NewGRPCMiddleware(
				fakeJWT{claims: jwt_go.MapClaims{"data": tt.data}},
				allowOPA{},
			).ValidateCredentials()

			called := false
			handler := func(ctx context.Context, req any) (any, error) {
				called = true
				return nil, nil
			}

			_, err := interceptor(incomingContext(), nil, &grpc.UnaryServerInfo{FullMethod: "/delivery.v1.DeliveryService/CreateDelivery"}, handler)

			assert.False(t, called)
			assert.Equal(t, codes.Unauthenticated, status.Code(err))
		})
	}
}
