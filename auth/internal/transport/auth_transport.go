package transport

import (
	"auth/internal/application/input"
	"auth/internal/application/usecase"
	"auth/pkg/authentication"
	"auth/pkg/authorization"
	"context"
	"fmt"
	authv1 "proto/gen/auth/v1"
)

type AuthTransport struct {
	authv1.UnimplementedAuthServiceServer
	uuc usecase.UserUsecaseInterface
	ruc usecase.RoleUsecaseInterface
	jwt authentication.TokenValidator
	opa authorization.OPAInterface
}

func NewAuthTransport(
	userUsecase usecase.UserUsecaseInterface,
	roleUsecase usecase.RoleUsecaseInterface,
	tokenValidator authentication.TokenValidator,
	opa authorization.OPAInterface,
) *AuthTransport {
	return &AuthTransport{
		uuc: userUsecase,
		ruc: roleUsecase,
		jwt: tokenValidator,
		opa: opa,
	}
}

func (a *AuthTransport) validateCredentials(token, action, path string) error {
	if err := a.jwt.Validate(token); err != nil {
		return err
	}

	claims, err := a.jwt.ExtractClaims(token)
	if err != nil {
		return err
	}

	data, ok := claims["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("jwt: claim 'data' inválida")
	}

	role, ok := data["role_name"].(string)
	if !ok {
		return fmt.Errorf("jwt: claim 'role_name' inválida")
	}

	opaInput := authorization.OPAInput{
		Action: action,
		Path:   path,
	}
	opaInput.User.Role = role

	if err := a.opa.Validate(opaInput); err != nil {
		return err
	}

	return nil
}

func (u *AuthTransport) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	login, err := u.uuc.Login(ctx, input.LoginInput{Email: req.Email, Password: req.Password})

	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		AccessToken: login.Token,
	}, nil

}

func (u *AuthTransport) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	err := u.validateCredentials(req.AuthField.Token, req.AuthField.Action, req.AuthField.Path)
	if err != nil {
		return nil, err
	}

	result, err := u.uuc.Create(ctx, input.CreateUserInput{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		RoleID:   req.RoleId,
	})

	if err != nil {
		return nil, err
	}

	return &authv1.CreateUserResponse{
		Id:    result.ID.String(),
		Name:  result.Email,
		Email: result.Email,
	}, nil
}

func (u *AuthTransport) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	err := u.validateCredentials(req.AuthField.Token, req.AuthField.Action, req.AuthField.Path)
	if err != nil {
		return nil, err
	}

	err = u.uuc.Delete(ctx, input.DeleteUserInput{ID: req.Id})

	if err != nil {
		return nil, err
	}

	return &authv1.DeleteUserResponse{
		Msg: "OK",
	}, nil
}

func (u *AuthTransport) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	err := u.validateCredentials(req.AuthField.Token, req.AuthField.Action, req.AuthField.Path)
	if err != nil {
		return nil, err
	}

	result, err := u.ruc.Create(ctx, input.CreateRoleInput{Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}

	var desc string
	if result.Description != nil {
		desc = *result.Description
	}

	return &authv1.CreateRoleResponse{
		Id:          result.ID.String(),
		Name:        result.Name,
		Description: desc,
	}, nil
}
func (u *AuthTransport) DeleteRole(ctx context.Context, req *authv1.DeleteRoleRequest) (*authv1.DeleteRoleResponse, error) {
	err := u.validateCredentials(req.AuthField.Token, req.AuthField.Action, req.AuthField.Path)
	if err != nil {
		return nil, err
	}

	err = u.ruc.Delete(ctx, input.DeleteRoleInput{ID: req.Id})

	if err != nil {
		return nil, err
	}

	return &authv1.DeleteRoleResponse{
		Msg: "OK",
	}, nil
}
