package transport

import (
	"context"
	"errors"

	"delivery/internal/domain/errs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errorCodes = []struct {
	err  error
	code codes.Code
}{
	{errs.ErrInvalidStreet, codes.InvalidArgument},
	{errs.ErrInvalidNumber, codes.InvalidArgument},
	{errs.ErrInvalidNeighborhood, codes.InvalidArgument},
	{errs.ErrInvalidCity, codes.InvalidArgument},
	{errs.ErrInvalidZipCode, codes.InvalidArgument},
	{errs.ErrInvalidWeight, codes.InvalidArgument},
	{errs.ErrInvalidClientID, codes.InvalidArgument},
	{errs.ErrInvalidDriverID, codes.InvalidArgument},
	{errs.ErrInvalidDeliveryID, codes.InvalidArgument},

	{errs.ErrDeliveryNotFound, codes.NotFound},

	{errs.ErrDeliveryAlreadyAssigned, codes.FailedPrecondition},
	{errs.ErrConcurrencyConflict, codes.Aborted},

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
