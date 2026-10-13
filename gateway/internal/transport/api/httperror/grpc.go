package httperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// exposedCodes are the gRPC codes our services use for domain sentinel errors.
// Their messages are safe to forward to the client; every other code gets a
// generic message, since it may carry transport or infrastructure details.
var exposedCodes = map[codes.Code]int{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.FailedPrecondition: http.StatusConflict,
	codes.Aborted:            http.StatusConflict,
}

var hiddenCodes = map[codes.Code]int{
	codes.Unavailable:      http.StatusServiceUnavailable,
	codes.DeadlineExceeded: http.StatusGatewayTimeout,
	codes.Canceled:         http.StatusRequestTimeout,
}

// FromGRPC maps an error returned by a gRPC client to an HTTP status code and
// a message that is safe to send to the client.
func FromGRPC(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}

	if httpStatus, ok := exposedCodes[st.Code()]; ok {
		return httpStatus, st.Message()
	}

	if httpStatus, ok := hiddenCodes[st.Code()]; ok {
		return httpStatus, http.StatusText(httpStatus)
	}

	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
}

// Respond writes the HTTP error response for an error returned by a gRPC client.
func Respond(ctx *gin.Context, err error) {
	httpStatus, msg := FromGRPC(err)
	ctx.JSON(httpStatus, gin.H{"error": msg})
}
