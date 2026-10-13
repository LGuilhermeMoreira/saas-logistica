package authentication

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "authenticated-user-id"

func ContextWithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// The only trusted source of the caller's identity: ids sent in request
// bodies must never decide who is acting.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok && id != uuid.Nil
}
