package entity

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"uniqueIndex"`
	Password string    `gorm:"not null"`
	RoleID   uuid.UUID `gorm:"not null"`
	Role     Role      `gorm:"foreignKey:RoleID"`
}

func NewClient(name, email, password, roleID string) (*Client, error) {

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("client must have a name")
	}

	_, mailErr := mail.ParseAddress(email)

	if strings.TrimSpace(email) == "" || mailErr != nil {
		return nil, errors.New("email is invalid")
	}

	if strings.TrimSpace(password) == "" || len(password) < 6 {
		return nil, errors.New("password does not meet the security requirements")
	}

	realRoleId, err := uuid.Parse(roleID)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Client{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: string(hashPassword),
		RoleID:   realRoleId,
	}, nil
}
