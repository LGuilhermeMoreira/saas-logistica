package transport

import (
	"context"
	"errors"

	"auth/internal/domain/errs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errorCodes = []struct {
	err  error
	code codes.Code
}{
	{errs.ErrInvalidUserName, codes.InvalidArgument},
	{errs.ErrInvalidEmail, codes.InvalidArgument},
	{errs.ErrWeakPassword, codes.InvalidArgument},
	{errs.ErrInvalidUserID, codes.InvalidArgument},
	{errs.ErrInvalidRoleID, codes.InvalidArgument},
	{errs.ErrInvalidRoleName, codes.InvalidArgument},
	{errs.ErrRoleWithoutPermission, codes.InvalidArgument},

	{errs.ErrInvalidCredentials, codes.Unauthenticated},

	{errs.ErrUserNotFound, codes.NotFound},
	{errs.ErrRoleNotFound, codes.NotFound},

	{errs.ErrUserAlreadyExists, codes.AlreadyExists},
	{errs.ErrRoleAlreadyExists, codes.AlreadyExists},
	{errs.ErrRoleInUse, codes.FailedPrecondition},

	{context.Canceled, codes.Canceled},
	{context.DeadlineExceeded, codes.DeadlineExceeded},
}

// toGRPCError converts a domain sentinel into a gRPC status carrying only the
// sentinel message. Anything else is reported as Internal with a generic
// message, so infrastructure details never leak to the gateway.
func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			return status.Error(e.code, e.err.Error())
		}
	}

	return status.Error(codes.Internal, errs.ErrInternal.Error())
}
