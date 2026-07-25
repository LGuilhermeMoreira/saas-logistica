package output

import "github.com/google/uuid"

type CreateUserOutput struct {
	ID     uuid.UUID
	Name   string
	Email  string
	RoleID uuid.UUID
}

type LoginOutput struct {
	Token string
}
