package errs

import "errors"

// Validation errors
var (
	ErrInvalidStreet       = errors.New("street is invalid")
	ErrInvalidNumber       = errors.New("number is invalid")
	ErrInvalidNeighborhood = errors.New("neighborhood is invalid")
	ErrInvalidCity         = errors.New("city is invalid")
	ErrInvalidZipCode      = errors.New("zip code is invalid")
	ErrInvalidWeight       = errors.New("weight is invalid")
	ErrInvalidClientID     = errors.New("client id is invalid")
	ErrInvalidDriverID     = errors.New("driver id is invalid")
	ErrInvalidDeliveryID   = errors.New("delivery id is invalid")
	ErrMissingAddress      = errors.New("origin and destination addresses are required")
)

// Not found errors
var (
	ErrDeliveryNotFound = errors.New("delivery not found")
)

// Conflict errors
var (
	ErrDeliveryAlreadyAssigned = errors.New("delivery already assigned by driver")
	ErrConcurrencyConflict     = errors.New("concurrency conflict: the delivery was modified by another process")
)

// Internal errors: the original cause must never reach the client
var (
	ErrInternal = errors.New("internal server error")
)
