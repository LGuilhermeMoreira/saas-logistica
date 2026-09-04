package repository

import (
	"context"
	"delivery/internal/domain/contract"
	"delivery/internal/domain/entity"
	"errors"
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

var ErrConcurrencyConflict = errors.New("concurrency conflict: the delivery was modified by another process")

func (d *DeliveryRepository) AssingToDriver(ctx context.Context, model *entity.Delivery) error {

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
		return err
	}

	if result.MatchedCount == 0 {
		return ErrConcurrencyConflict
	}

	return nil
}

func (d *DeliveryRepository) Create(ctx context.Context, model *entity.Delivery) error {

	coll := d.db.Collection("deliveries")

	_, err := coll.InsertOne(ctx, model)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeliveryRepository) FindByID(ctx context.Context, ID uuid.UUID) (*entity.Delivery, error) {
	coll := d.db.Collection("deliveries")

	filter := bson.M{
		"_id": ID.String(),
	}

	var model entity.Delivery

	err := coll.FindOne(ctx, filter).Decode(&model)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("delivery not found")
		}
		return nil, err
	}

	return &model, nil
}
