package entity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Description *string
}

func NewRole(name, description string) (*Role, error) {
	description = strings.TrimSpace(description)
	name = strings.TrimSpace(name)

	var realDescription *string

	if name == "" {
		return nil, errors.New("name is invalid")
	}

	if description != "" {
		realDescription = &description
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate role id: %w", err)
	}

	return &Role{
		ID:          id,
		Name:        name,
		Description: realDescription,
	}, nil
}
