package repository

import (
	"context"
	"delivery/internal/domain/contract"
	"delivery/internal/domain/entity"
	"delivery/internal/domain/errs"
	"delivery/pkg/logger"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeliveryRepository struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewDeliveryRepository(db *mongo.Database, log *slog.Logger) contract.DeliveryRepositoryInterface {
	return &DeliveryRepository{
		db: db, log: log,
	}
}

func (d *DeliveryRepository) AssingToDriver(ctx context.Context, model *entity.Delivery) error {
	d.log.Debug("assigning driver to delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("driver_id", model.DriverID.String()))

	coll := d.db.Collection("deliveries")

	filter := bson.M{
		"_id":       model.ID,
		"driver_id": bson.M{"$exists": false},
		"status":    entity.DeliveryStatusCreated,
	}

	update := bson.M{
		"$set": bson.M{
			"driver_id":  model.DriverID,
			"status":     model.Status,
			"updated_at": time.Now(),
		},
	}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		d.log.Error("failed to assign driver to delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", model.ID.String()), slog.String("driver_id", model.DriverID.String()))
		return fmt.Errorf("failed to assign driver to delivery: %w", err)
	}

	if result.MatchedCount == 0 {
		d.log.Error("concurrency conflict when assigning driver", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("driver_id", model.DriverID.String()))
		return errs.ErrConcurrencyConflict
	}

	d.log.Info("driver assigned to delivery successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("driver_id", model.DriverID.String()))
	return nil
}

func (d *DeliveryRepository) Create(ctx context.Context, model *entity.Delivery) error {
	d.log.Debug("creating delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("client_id", model.ClientID.String()))

	coll := d.db.Collection("deliveries")

	_, err := coll.InsertOne(ctx, model)
	if err != nil {
		d.log.Error("failed to create delivery", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", model.ID.String()), slog.String("client_id", model.ClientID.String()))
		return fmt.Errorf("failed to create delivery: %w", err)
	}

	d.log.Info("delivery created successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", model.ID.String()), slog.String("client_id", model.ClientID.String()))
	return nil
}

func (d *DeliveryRepository) FindByID(ctx context.Context, ID uuid.UUID) (*entity.Delivery, error) {
	d.log.Debug("finding delivery by ID", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", ID.String()))

	coll := d.db.Collection("deliveries")

	filter := bson.M{
		"_id": ID,
	}

	var model entity.Delivery

	err := coll.FindOne(ctx, filter).Decode(&model)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			d.log.Error("delivery not found", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", ID.String()))
			return nil, errs.ErrDeliveryNotFound
		}
		d.log.Error("failed to find delivery by ID", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()), slog.String("delivery_id", ID.String()))
		return nil, fmt.Errorf("failed to find delivery by id: %w", err)
	}

	d.log.Info("delivery found successfully", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("delivery_id", ID.String()))
	return &model, nil
}

func (d *DeliveryRepository) FindUnassociatedDeliveries(ctx context.Context) ([]entity.Delivery, error) {
	d.log.Debug("finding unassociated deliveries", slog.String("request-id", logger.ExtractRequestID(ctx)))

	coll := d.db.Collection("deliveries")

	filter := bson.M{
		"status":    entity.DeliveryStatusCreated,
		"driver_id": bson.M{"$exists": false},
	}

	cursor, err := coll.Find(ctx, filter)

	if err != nil {
		d.log.Error("failed to find unassociated deliveries", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to find unassociated deliveries: %w", err)
	}
	defer cursor.Close(ctx)

	var results []entity.Delivery

	if err := cursor.All(ctx, &results); err != nil {
		d.log.Error("failed to decode unassociated deliveries", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode unassociated deliveries: %w", err)
	}

	d.log.Info("unassociated deliveries found", slog.String("request-id", logger.ExtractRequestID(ctx)), slog.Int("count", len(results)))
	return results, nil
}
