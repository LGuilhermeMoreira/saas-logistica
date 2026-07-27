package transport

import (
	"auth/internal/application/input"
	"auth/internal/application/usecase"
	"context"
	authv1 "proto/gen/auth/v1"

	"gorm.io/gorm"
)

type AuthTransport struct {
	authv1.UnimplementedAuthServiceServer
	db  *gorm.DB
	uuc usecase.UserUsecaseInterface
	ruc usecase.RoleUsecaseInterface
}

func NewAuthTransport(db *gorm.DB, uuc usecase.UserUsecaseInterface, ruc usecase.RoleUsecaseInterface) *AuthTransport {
	return &AuthTransport{
		db:  db,
		uuc: uuc,
		ruc: ruc,
	}
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
	err := u.uuc.Delete(ctx, input.DeleteUserInput{ID: req.Id})

	if err != nil {
		return nil, err
	}

	return &authv1.DeleteUserResponse{
		Msg: "OK",
	}, nil
}

func (u *AuthTransport) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
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
	err := u.ruc.Delete(ctx, input.DeleteRoleInput{ID: req.Id})

	if err != nil {
		return nil, err
	}

	return &authv1.DeleteRoleResponse{
		Msg: "OK",
	}, nil
}
