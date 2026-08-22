package repository

import (
	"context"
	"fmt"
	"log/slog"

	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) contract.UserRepositoryInterface {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (c *UserRepository) Create(ctx context.Context, model *entity.User) error {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Debug("executing create user query", slog.String("request_id", reqID), slog.String("user_id", model.ID.String()))

	err := c.db.WithContext(ctx).Create(model).Error
	if err != nil {
		c.log.Error("database error creating user", slog.String("request_id", reqID), slog.String("user_id", model.ID.String()), slog.Any("error", err))
		return fmt.Errorf("failed to persist User in database: %w", err)
	}
	return nil
}

func (c *UserRepository) Delete(ctx context.Context, model *entity.User) error {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Debug("executing delete user query", slog.String("request_id", reqID), slog.String("user_id", model.ID.String()))

	err := c.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		c.log.Error("database error deleting user", slog.String("request_id", reqID), slog.String("user_id", model.ID.String()), slog.Any("error", err))
		return fmt.Errorf("failed to delete User in database: %w", err)
	}
	return nil
}

func (c *UserRepository) Login(ctx context.Context, email string) (*entity.User, error) {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Debug("executing find user by email query for login", slog.String("request_id", reqID), slog.String("email", email))

	var User entity.User

	err := c.db.WithContext(ctx).Preload("Role", nil).Where("email = ?", email).First(&User).Error
	if err != nil {
		c.log.Error("database error finding user by email", slog.String("request_id", reqID), slog.String("email", email), slog.Any("error", err))
		return nil, fmt.Errorf("failed to find User by email: %w", err)
	}

	return &User, nil
}

func (c *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Debug("executing find user by id query", slog.String("request_id", reqID), slog.String("user_id", id.String()))

	var result entity.User

	err := c.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		c.log.Error("database error finding user by id", slog.String("request_id", reqID), slog.String("user_id", id.String()), slog.Any("error", err))
		return nil, fmt.Errorf("failed to find User by id: %w", err)
	}

	return &result, nil
}
