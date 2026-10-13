package entity

import (
	"auth/internal/domain/errs"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"uniqueIndex"`
	Password string    `gorm:"not null"`
	RoleID   uuid.UUID `gorm:"not null"`
	Role     Role      `gorm:"foreignKey:RoleID;references:ID"`
}

func NewUser(name, email, password, roleID string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if name == "" {
		return nil, errs.ErrInvalidUserName
	}

	_, mailErr := mail.ParseAddress(email)

	if email == "" || mailErr != nil {
		return nil, errs.ErrInvalidEmail
	}

	if password == "" || len(password) < 6 {
		return nil, errs.ErrWeakPassword
	}

	realRoleId, err := uuid.Parse(roleID)
	if err != nil || realRoleId == uuid.Nil {
		return nil, errs.ErrInvalidRoleID
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user id: %w", err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: string(hashPassword),
		RoleID:   realRoleId,
	}, nil
}

func (c User) ToMap() map[string]any {
	return map[string]any{
		"id":    c.ID,
		"name":  c.Name,
		"email": c.Email,
		// "password":   c.Password,
		"role_id":    c.RoleID,
		"role_name":  c.Role.Name,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
		"deleted_at": c.DeletedAt,
	}
}
