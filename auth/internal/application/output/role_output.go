package output

import "github.com/google/uuid"

type RoleOutput struct {
	ID          uuid.UUID
	Name        string
	Description *string
}
