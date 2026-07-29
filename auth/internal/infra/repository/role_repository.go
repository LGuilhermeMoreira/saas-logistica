package repository

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"

	"github.com/google/uuid"
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

func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var result entity.Role

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	var result []entity.Role

	err := r.db.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
