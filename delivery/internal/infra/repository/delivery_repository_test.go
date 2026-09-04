package repository_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"delivery/internal/domain/entity"
	"delivery/internal/infra/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongo(t *testing.T) *mongo.Database {
	t.Helper()

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

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = client.Disconnect(context.Background())
	})

	db := client.Database("delivery_test")

	return db
}

func newRepository(t *testing.T) (*mongo.Database, *repository.DeliveryRepository) {
	t.Helper()

	db := setupMongo(t)

	log := slog.New(slog.NewTextHandler(nil, nil))

	repo := repository.NewDeliveryRepository(db, log)

	return db, repo.(*repository.DeliveryRepository)
}

func newDelivery(t *testing.T) *entity.Delivery {
	t.Helper()

	return &entity.Delivery{
		ID: uuid.New(),
		To: entity.Address{
			ID:           uuid.New(),
			Street:       "Rua A",
			Number:       "100",
			Neighborhood: "Centro",
			City:         "Quixadá",
			ZipCode:      "63900-000",
		},
		From: entity.Address{
			ID:           uuid.New(),
			Street:       "Rua B",
			Number:       "200",
			Neighborhood: "Centro",
			City:         "Quixadá",
			ZipCode:      "63900-000",
		},
		Weight:   10,
		ClientID: uuid.New(),
		Status:   entity.DeliveryStatusCreated,
	}
}

func TestDeliveryRepository_Create(t *testing.T) {
	db, repo := newRepository(t)

	ctx := context.Background()
	delivery := newDelivery(t)

	err := repo.Create(ctx, delivery)

	require.NoError(t, err)

	var result entity.Delivery

	err = db.Collection("deliveries").
		FindOne(ctx, bson.M{"_id": delivery.ID}).
		Decode(&result)

	require.NoError(t, err)

	assert.Equal(t, delivery.ID, result.ID)
	assert.Equal(t, delivery.To, result.To)
	assert.Equal(t, delivery.From, result.From)
	assert.Equal(t, delivery.Weight, result.Weight)
	assert.Equal(t, delivery.ClientID, result.ClientID)
	assert.Equal(t, delivery.Status, result.Status)
}

func TestDeliveryRepository_FindByID(t *testing.T) {
	t.Run("should find delivery", func(t *testing.T) {
		_, repo := newRepository(t)

		ctx := context.Background()
		delivery := newDelivery(t)

		err := repo.Create(ctx, delivery)
		require.NoError(t, err)

		result, err := repo.FindByID(ctx, delivery.ID)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, delivery.ID, result.ID)
		assert.Equal(t, delivery.Weight, result.Weight)
		assert.Equal(t, delivery.ClientID, result.ClientID)
		assert.Equal(t, delivery.Status, result.Status)
	})

	t.Run("should return error when delivery does not exist", func(t *testing.T) {
		_, repo := newRepository(t)

		result, err := repo.FindByID(
			context.Background(),
			uuid.New(),
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "delivery not found")
	})
}

func TestDeliveryRepository_FindUnassociatedDeliveries(t *testing.T) {
	t.Run("should return only created and unassociated deliveries", func(t *testing.T) {
		_, repo := newRepository(t)

		ctx := context.Background()

		delivery1 := newDelivery(t)
		delivery2 := newDelivery(t)
		delivery3 := newDelivery(t)
		delivery4 := newDelivery(t)

		// Deve aparecer.
		require.NoError(t, repo.Create(ctx, delivery1))

		// Deve aparecer.
		require.NoError(t, repo.Create(ctx, delivery2))

		// Não deve aparecer: já possui motorista.
		driverID := uuid.New()
		delivery3.DriverID = &driverID
		delivery3.Status = entity.DeliveryStatusInTransit

		require.NoError(t, repo.Create(ctx, delivery3))

		// Não deve aparecer: status diferente de created.
		delivery4.Status = entity.DeliveryStatusCancelled

		require.NoError(t, repo.Create(ctx, delivery4))

		results, err := repo.FindUnassociatedDeliveries(ctx)

		require.NoError(t, err)

		assert.Len(t, results, 2)

		ids := []uuid.UUID{
			results[0].ID,
			results[1].ID,
		}

		assert.ElementsMatch(
			t,
			[]uuid.UUID{
				delivery1.ID,
				delivery2.ID,
			},
			ids,
		)
	})
}

func TestDeliveryRepository_AssingToDriver(t *testing.T) {
	t.Run("should assign driver", func(t *testing.T) {
		db, repo := newRepository(t)

		ctx := context.Background()

		delivery := newDelivery(t)

		require.NoError(t, repo.Create(ctx, delivery))

		driverID := uuid.New()

		delivery.DriverID = &driverID
		delivery.Status = entity.DeliveryStatusInTransit

		err := repo.AssingToDriver(ctx, delivery)

		require.NoError(t, err)

		var result entity.Delivery

		err = db.Collection("deliveries").
			FindOne(
				ctx,
				bson.M{"_id": delivery.ID},
			).
			Decode(&result)

		require.NoError(t, err)

		assert.Equal(t, delivery.ID, result.ID)
		assert.Equal(t, delivery.Weight, result.Weight)
		assert.Equal(t, delivery.ClientID, result.ClientID)
		assert.Equal(t, delivery.Status, result.Status)
	})

	t.Run("should return concurrency conflict when delivery does not exist", func(t *testing.T) {
		_, repo := newRepository(t)

		delivery := newDelivery(t)

		driverID := uuid.New()
		delivery.DriverID = &driverID
		delivery.Status = entity.DeliveryStatusInTransit

		err := repo.AssingToDriver(
			context.Background(),
			delivery,
		)

		require.Error(t, err)
		assert.ErrorIs(
			t,
			err,
			repository.ErrConcurrencyConflict,
		)
	})

	t.Run("should return concurrency conflict when delivery already has driver", func(t *testing.T) {
		_, repo := newRepository(t)

		ctx := context.Background()

		delivery := newDelivery(t)

		firstDriverID := uuid.New()
		delivery.DriverID = &firstDriverID
		delivery.Status = entity.DeliveryStatusInTransit

		require.NoError(t, repo.Create(ctx, delivery))

		secondDriverID := uuid.New()
		delivery.DriverID = &secondDriverID

		err := repo.AssingToDriver(ctx, delivery)

		require.Error(t, err)
		assert.ErrorIs(
			t,
			err,
			repository.ErrConcurrencyConflict,
		)
	})

	t.Run("should return concurrency conflict when status is not created", func(t *testing.T) {
		_, repo := newRepository(t)

		ctx := context.Background()

		delivery := newDelivery(t)

		delivery.Status = entity.DeliveryStatusCancelled

		require.NoError(t, repo.Create(ctx, delivery))

		driverID := uuid.New()
		delivery.DriverID = &driverID
		delivery.Status = entity.DeliveryStatusInTransit

		err := repo.AssingToDriver(ctx, delivery)

		require.Error(t, err)
		assert.ErrorIs(
			t,
			err,
			repository.ErrConcurrencyConflict,
		)
	})
}

func TestDeliveryRepository_AssingToDriver_Concurrency(t *testing.T) {
	db, repo := newRepository(t)

	ctx := context.Background()

	delivery := newDelivery(t)

	require.NoError(t, repo.Create(ctx, delivery))

	driver1 := uuid.New()
	driver2 := uuid.New()

	delivery1 := *delivery
	delivery1.DriverID = &driver1
	delivery1.Status = entity.DeliveryStatusInTransit

	delivery2 := *delivery
	delivery2.DriverID = &driver2
	delivery2.Status = entity.DeliveryStatusInTransit

	errCh := make(chan error, 2)

	go func() {
		errCh <- repo.AssingToDriver(ctx, &delivery1)
	}()

	go func() {
		errCh <- repo.AssingToDriver(ctx, &delivery2)
	}()

	err1 := <-errCh
	err2 := <-errCh

	// Exatamente uma operação deve conseguir
	// atribuir o motorista.
	assert.True(
		t,
		(err1 == nil && err2 != nil) ||
			(err1 != nil && err2 == nil),
	)

	var result entity.Delivery

	err := db.Collection("deliveries").
		FindOne(ctx, bson.M{"_id": delivery.ID}).
		Decode(&result)

	require.NoError(t, err)

	require.NotNil(t, result.DriverID)

	assert.Equal(
		t,
		entity.DeliveryStatusInTransit,
		result.Status,
	)

	assert.Contains(
		t,
		[]uuid.UUID{driver1, driver2},
		*result.DriverID,
	)

	assert.Eventually(
		t,
		func() bool {
			return result.DriverID != nil
		},
		time.Second,
		10*time.Millisecond,
	)
}
