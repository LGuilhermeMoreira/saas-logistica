package httperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Codes used by the services' sentinel errors: their messages are safe to
// forward. Any other code may be produced by the gRPC client itself and carry
// connection details, so the client only gets a generic message.
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

func Respond(ctx *gin.Context, err error) {
	httpStatus, msg := FromGRPC(err)
	ctx.JSON(httpStatus, gin.H{"error": msg})
}
