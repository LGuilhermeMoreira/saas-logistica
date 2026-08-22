package contract

import (
	"context"
	authv1 "proto/gen/auth/v1"
)

type AuthService interface {
	Login(context context.Context, email, password string) (string, error)
	CreateUser(context context.Context, name, email, password, roleID string) (*authv1.CreateUserResponse, error)
	DeleteUser(context context.Context, id string) (*authv1.DeleteUserResponse, error)
	CreateRole(context context.Context, name, description string, permissions []*authv1.Permission) (*authv1.CreateRoleResponse, error)
	DeleteRole(context context.Context, id string) (*authv1.DeleteRoleResponse, error)
	FindUserByID(ctx context.Context, id string) (*authv1.FindUserByIDResponse, error)
	FindRoleByID(ctx context.Context, id string) (*authv1.FindRoleByIDResponse, error)
	FindAllRoles(ctx context.Context) (*authv1.FindAllRolesResponse, error)
}
