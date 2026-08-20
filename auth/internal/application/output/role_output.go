package output

import "github.com/google/uuid"

type PermissionsOutput struct {
	Action string
	Path   string
}

type RoleOutput struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Permissions []PermissionsOutput
}
