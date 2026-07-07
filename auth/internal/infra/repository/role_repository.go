package repository

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) contract.RoleRepositoryInterface {
	return &RoleRepository{db: db}

}

func (r *RoleRepository) Create(ctx context.Context, model *entity.Role) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *RoleRepository) Delete(ctx context.Context, model *entity.Role) error {
	return r.db.WithContext(ctx).Delete(model).Error
}
