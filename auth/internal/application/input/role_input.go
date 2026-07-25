package input

type CreateRoleInput struct {
	Name        string
	Descritpion string
	Permissions map[string]any
}

type DeleteRoleInput struct {
	ID string
}
