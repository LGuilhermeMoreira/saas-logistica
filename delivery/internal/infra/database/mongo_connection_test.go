package database_test

import (
	"context"
	"delivery/internal/infra/database"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestNewMongoConnection(t *testing.T) {
	ctx := context.Background()

	container, err := mongodb.Run(
		ctx,
		"mongo:7",
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := database.NewMongoConnection(uri)

	require.NoError(t, err)
	require.NotNil(t, client)

	t.Cleanup(func() {
		disconnectCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		require.NoError(t, client.Disconnect(disconnectCtx))
	})

	err = client.Ping(ctx, nil)
	require.NoError(t, err)
}
