package repository

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) contract.ClientRepositoryInterface {
	return &ClientRepository{
		db: db,
	}
}

func (c *ClientRepository) Create(ctx context.Context, model *entity.Client) error {

	// não gostei de utilizar generics, mesmo sabendo as vantagens
	// 	return gorm.G[entity.Client](c.db).Create(ctx, model)
	err := c.db.WithContext(ctx).Create(model).Error
	return fmt.Errorf("failed to persist client in database: %w", err)
}

func (c *ClientRepository) Delete(ctx context.Context, model *entity.Client) error {
	err := c.db.WithContext(ctx).Delete(model).Error

	return fmt.Errorf("failed to delete client in database: %w", err)
}

func (c *ClientRepository) Login(ctx context.Context, email string) (*entity.Client, error) {
	var client entity.Client

	err := c.db.WithContext(ctx).Where("email = ?", email).First(&client).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find client by email: %w", err)
	}

	return &client, nil
}

func (c *ClientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	var result entity.Client
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find client by id: %w", err)
	}
	return &result, nil
}
