package database_test

import (
	"auth/internal/infra/database"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewPostgresConnection(t *testing.T) {

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("erro ao terminar container: %v", err)
		}
	}()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	t.Run("successfully connects and configures the pool", func(t *testing.T) {
		db, err := database.NewPostgresConnection(dsn)
		require.NoError(t, err)
		require.NotNil(t, db)

		sqlDB, err := db.DB()
		require.NoError(t, err)
		defer sqlDB.Close()

		require.NoError(t, sqlDB.Ping())

		stats := sqlDB.Stats()
		require.Equal(t, 100, stats.MaxOpenConnections)
	})

	t.Run("fails and retries until attempts are exhausted", func(t *testing.T) {
		dsnInvalido := "postgres://testuser:testpass@localhost:1/testdb?sslmode=disable"

		start := time.Now()
		db, err := database.NewPostgresConnection(dsnInvalido)
		elapsed := time.Since(start)

		require.Error(t, err)
		require.Nil(t, db)
		require.EqualError(t, err, "Unable to connect to the database after 5 attempts.")

		require.GreaterOrEqual(t, elapsed, 8*time.Second)
	})
}
