package input

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	RoleID   string
}

type LoginInput struct {
	Email    string
	Password string
}

type DeleteUserInput struct {
	ID string
}
