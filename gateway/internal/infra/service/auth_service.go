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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := s.authClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

func (s *AuthService) CreateUser(ctx context.Context, name, email, password, roleID string) (*authv1.CreateUserResponse, error) {
	resp, err := s.authClient.CreateUser(ctx, &authv1.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		RoleId:   roleID,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, id string) (*authv1.DeleteUserResponse, error) {
	resp, err := s.authClient.DeleteUser(ctx, &authv1.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) CreateRole(ctx context.Context, name, description string, permissions []*authv1.Permission) (*authv1.CreateRoleResponse, error) {
	resp, err := s.authClient.CreateRole(ctx, &authv1.CreateRoleRequest{
		Name:        name,
		Description: description,
		Permissions: permissions,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) DeleteRole(ctx context.Context, id string) (*authv1.DeleteRoleResponse, error) {
	resp, err := s.authClient.DeleteRole(ctx, &authv1.DeleteRoleRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) FindUserByID(ctx context.Context, id string) (*authv1.FindUserByIDResponse, error) {
	resp, err := s.authClient.FindUserByID(ctx, &authv1.FindUserByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) FindRoleByID(ctx context.Context, id string) (*authv1.FindRoleByIDResponse, error) {
	resp, err := s.authClient.FindRoleByID(ctx, &authv1.FindRoleByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) FindAllRoles(ctx context.Context) (*authv1.FindAllRolesResponse, error) {
	resp, err := s.authClient.FindAllRoles(ctx, &authv1.FindAllRolesRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
