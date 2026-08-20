package contract

import (
	authv1 "proto/gen/auth/v1"
)

type AuthService interface {
	Login(email, password string) (string, error)
	CreateUser(name, email, password, roleID string, authField *authv1.AuthField) (*authv1.CreateUserResponse, error)
	DeleteUser(id string, authField *authv1.AuthField) (*authv1.DeleteUserResponse, error)
	CreateRole(name, description string, permissions []*authv1.Permission, authField *authv1.AuthField) (*authv1.CreateRoleResponse, error)
	DeleteRole(id string, authField *authv1.AuthField) (*authv1.DeleteRoleResponse, error)
}
