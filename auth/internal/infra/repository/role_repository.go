package repository

import (
	"context"
	"log/slog"

	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewRoleRepository(db *gorm.DB, log *slog.Logger) contract.RoleRepositoryInterface {
	return &RoleRepository{
		db:  db,
		log: log,
	}
}

func (r *RoleRepository) Create(ctx context.Context, model *entity.Role) error {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Debug("executing create role query", slog.String("request_id", reqID), slog.String("role_id", model.ID.String()))

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		r.log.Error("database error creating role", slog.String("request_id", reqID), slog.String("role_id", model.ID.String()), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, model *entity.Role) error {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Debug("executing delete role query", slog.String("request_id", reqID), slog.String("role_id", model.ID.String()))

	err := r.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		r.log.Error("database error deleting role", slog.String("request_id", reqID), slog.String("role_id", model.ID.String()), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Debug("executing find role by id query", slog.String("request_id", reqID), slog.String("role_id", id.String()))

	var result entity.Role

	err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ?", id).First(&result).Error
	if err != nil {
		r.log.Error("database error finding role by id", slog.String("request_id", reqID), slog.String("role_id", id.String()), slog.Any("error", err))
		return nil, err
	}

	return &result, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Debug("executing find all roles query", slog.String("request_id", reqID))

	var result []entity.Role

	err := r.db.WithContext(ctx).Preload("Permissions").Find(&result).Error
	if err != nil {
		r.log.Error("database error finding all roles", slog.String("request_id", reqID), slog.Any("error", err))
		return nil, err
	}

	return result, nil
}
