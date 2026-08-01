package repository

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func (c *UserRepository) Create(ctx context.Context, model *entity.User) error {

	// não gostei de utilizar generics, mesmo sabendo as vantagens
	// 	return gorm.G[entity.User](c.db).Create(ctx, model)
	err := c.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return fmt.Errorf("failed to persist User in database: %w", err)
	}
	return nil
}

func (c *UserRepository) Delete(ctx context.Context, model *entity.User) error {
	err := c.db.WithContext(ctx).Delete(model).Error
	if err != nil {
		return fmt.Errorf("failed to delete User in database: %w", err)
	}
	return nil
}

func (c *UserRepository) Login(ctx context.Context, email string) (*entity.User, error) {
	var User entity.User

	err := c.db.WithContext(ctx).Preload("Roles", nil).Where("email = ?", email).First(&User).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find User by email: %w", err)
	}

	return &User, nil
}

func (c *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var result entity.User
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find User by id: %w", err)
	}
	return &result, nil
}
