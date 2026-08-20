package service

import (
	"context"
	"gateway/internal/domain/contract"
	authv1 "proto/gen/auth/v1"
)

type AuthService struct {
	authClient authv1.AuthServiceClient
}

func NewAuthService(authClient authv1.AuthServiceClient) contract.AuthService {
	return &AuthService{
		authClient: authClient,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	resp, err := s.authClient.Login(context.Background(), &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

func (s *AuthService) CreateUser(name, email, password, roleID string, authField *authv1.AuthField) (*authv1.CreateUserResponse, error) {
	resp, err := s.authClient.CreateUser(context.Background(), &authv1.CreateUserRequest{
		Name:      name,
		Email:     email,
		Password:  password,
		RoleId:    roleID,
		AuthField: authField,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) DeleteUser(id string, authField *authv1.AuthField) (*authv1.DeleteUserResponse, error) {
	resp, err := s.authClient.DeleteUser(context.Background(), &authv1.DeleteUserRequest{
		Id:        id,
		AuthField: authField,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) CreateRole(name, description string, permissions []*authv1.Permission, authField *authv1.AuthField) (*authv1.CreateRoleResponse, error) {
	resp, err := s.authClient.CreateRole(context.Background(), &authv1.CreateRoleRequest{
		Name:        name,
		Description: description,
		Permissions: permissions,
		AuthField:   authField,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) DeleteRole(id string, authField *authv1.AuthField) (*authv1.DeleteRoleResponse, error) {
	resp, err := s.authClient.DeleteRole(context.Background(), &authv1.DeleteRoleRequest{
		Id:        id,
		AuthField: authField,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
