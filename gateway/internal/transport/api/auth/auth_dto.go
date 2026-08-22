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

type UserOutput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID string `json:"role_id"`
}

type DeleteUserInput struct {
	ID string `json:"id"`
}

type DeleteUserOutput struct {
	Msg string `json:"msg"`
}

type Permission struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type CreateRoleInput struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

type RoleOutput struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

type DeleteRoleInput struct {
	ID string `json:"id"`
}

type DeleteRoleOutput struct {
	Msg string `json:"msg"`
}

type FindAllRolesOutput struct {
	Roles []RoleOutput `json:"roles"`
}
