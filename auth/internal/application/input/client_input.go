package input

type CreateClientInput struct {
	Name     string
	Email    string
	Password string
	RoleID   string
}

type LoginInput struct {
	Email    string
	Password string
}

type DeleteClientInput struct {
	ID string
}
