package auth

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

type CreateUserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteUserInput struct {
	ID string `json:"id"`
}

type DeleteUserOutput struct {
	Msg string `json:"msg"`
}

type PermissionInput struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type CreateRoleInput struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions []PermissionInput `json:"permissions"`
}

type CreateRoleOutput struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions []PermissionInput `json:"permissions"`
}

type DeleteRoleInput struct {
	ID string `json:"id"`
}

type DeleteRoleOutput struct {
	Msg string `json:"msg"`
}
