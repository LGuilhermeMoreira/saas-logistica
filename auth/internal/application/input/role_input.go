package input

type PermissionInput struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type CreateRoleInput struct {
	Name        string
	Description string
	Permissions []PermissionInput
}

func (c CreateRoleInput) PermissionToMap() []map[string]any {
	var response []map[string]any

	for _, v := range c.Permissions {
		i := map[string]any{
			"action": v.Action,
			"path":   v.Path,
		}
		response = append(response, i)
	}

	return response
}

type DeleteRoleInput struct {
	ID string
}
