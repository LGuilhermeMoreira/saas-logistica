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
	Permissions []Permission `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;" json:"permissions"`
}

type Permission struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	RoleID uuid.UUID `gorm:"type:uuid;index;not null" json:"-"`
	Action string    `gorm:"not null" json:"action"`
	Path   string    `gorm:"not null" json:"path"`
}

func NewPermission(action, path string) (*Permission, error) {
	permissionID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate permission id: %w", err)
	}

	result := Permission{
		ID:     permissionID,
		Action: action,
		Path:   path,
	}

	return &result, nil
}

func NewRole(name, description string, perm []Permission) (*Role, error) {
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

	if len(perm) == 0 {
		return nil, errors.New("role must have one permission")
	}

	for i := 0; i < len(perm); i++ {
		perm[i].RoleID = id
	}

	return &Role{
		ID:          id,
		Name:        name,
		Description: realDescription,
		Permissions: perm,
	}, nil
}
