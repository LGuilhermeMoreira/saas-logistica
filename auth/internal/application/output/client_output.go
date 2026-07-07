package output

import "github.com/google/uuid"

type CreateClientOutput struct {
	ID     uuid.UUID
	Name   string
	Email  string
	RoleID uuid.UUID
}

type LoginOutput struct {
	Token string
}
