package transport

import (
	"auth/internal/application/input"
	"auth/internal/application/usecase"
	"auth/pkg/logger"
	"context"
	"log/slog"
	authv1 "proto/gen/auth/v1"
)

type AuthTransport struct {
	authv1.UnimplementedAuthServiceServer
	uuc usecase.UserUsecaseInterface
	ruc usecase.RoleUsecaseInterface
	log *slog.Logger
}

func NewAuthTransport(
	userUsecase usecase.UserUsecaseInterface,
	roleUsecase usecase.RoleUsecaseInterface,
	logger *slog.Logger,
) *AuthTransport {
	return &AuthTransport{
		uuc: userUsecase,
		ruc: roleUsecase,
		log: logger,
	}
}

func (u *AuthTransport) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing login request", slog.String("request_id", reqID), slog.String("email", req.Email))

	login, err := u.uuc.Login(ctx, input.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		u.log.Error("login failed", slog.String("request_id", reqID), slog.String("email", req.Email), slog.Any("error", err))
		return nil, err
	}

	u.log.Info("login successful", slog.String("request_id", reqID), slog.String("email", req.Email))
	return &authv1.LoginResponse{
		AccessToken: login.Token,
	}, nil
}

func (u *AuthTransport) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing create user request", slog.String("request_id", reqID), slog.String("email", req.Email))

	result, err := u.uuc.Create(ctx, input.CreateUserInput{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		RoleID:   req.RoleId,
	})
	if err != nil {
		u.log.Error("failed to create user", slog.String("request_id", reqID), slog.String("email", req.Email), slog.Any("error", err))
		return nil, err
	}

	u.log.Info("user created successfully", slog.String("request_id", reqID), slog.String("user_id", result.ID.String()))
	return &authv1.CreateUserResponse{
		Id:    result.ID.String(),
		Name:  result.Name,
		Email: result.Email,
	}, nil
}

func (u *AuthTransport) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing delete user request", slog.String("request_id", reqID), slog.String("user_id", req.Id))

	err := u.uuc.Delete(ctx, input.DeleteUserInput{ID: req.Id})
	if err != nil {
		u.log.Error("failed to delete user", slog.String("request_id", reqID), slog.String("user_id", req.Id), slog.Any("error", err))
		return nil, err
	}

	u.log.Info("user deleted successfully", slog.String("request_id", reqID), slog.String("user_id", req.Id))
	return &authv1.DeleteUserResponse{
		Msg: "OK",
	}, nil
}

func (u *AuthTransport) CreateRole(ctx context.Context, req *authv1.CreateRoleRequest) (*authv1.CreateRoleResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing create role request", slog.String("request_id", reqID), slog.String("role_name", req.Name))

	permissions := make([]input.PermissionInput, len(req.Permissions))
	for index, permission := range req.Permissions {
		permissions[index] = input.PermissionInput{
			Action: permission.Action,
			Path:   permission.Path,
		}
	}

	result, err := u.ruc.Create(ctx, input.CreateRoleInput{Name: req.Name, Description: req.Description, Permissions: permissions})
	if err != nil {
		u.log.Error("failed to create role", slog.String("request_id", reqID), slog.String("role_name", req.Name), slog.Any("error", err))
		return nil, err
	}

	var desc string
	if result.Description != nil {
		desc = *result.Description
	}

	perm := make([]*authv1.Permission, len(result.Permissions))
	for i, v := range result.Permissions {
		perm[i] = &authv1.Permission{
			Action: v.Action,
			Path:   v.Path,
		}
	}

	u.log.Info("role created successfully", slog.String("request_id", reqID), slog.String("role_id", result.ID.String()))
	return &authv1.CreateRoleResponse{
		Id:          result.ID.String(),
		Name:        result.Name,
		Description: desc,
		Permissions: perm,
	}, nil
}

func (u *AuthTransport) DeleteRole(ctx context.Context, req *authv1.DeleteRoleRequest) (*authv1.DeleteRoleResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing delete role request", slog.String("request_id", reqID), slog.String("role_id", req.Id))

	err := u.ruc.Delete(ctx, input.DeleteRoleInput{ID: req.Id})
	if err != nil {
		u.log.Error("failed to delete role", slog.String("request_id", reqID), slog.String("role_id", req.Id), slog.Any("error", err))
		return nil, err
	}

	u.log.Info("role deleted successfully", slog.String("request_id", reqID), slog.String("role_id", req.Id))
	return &authv1.DeleteRoleResponse{
		Msg: "OK",
	}, nil
}

func (u *AuthTransport) FindUserByID(ctx context.Context, req *authv1.FindUserByIDRequest) (*authv1.FindUserByIDResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing find user by id request", slog.String("request_id", reqID), slog.String("user_id", req.Id))

	result, err := u.uuc.FindByID(ctx, input.FindByIDUserInput{ID: req.Id})
	if err != nil {
		u.log.Error("failed to find user", slog.String("request_id", reqID), slog.String("user_id", req.Id), slog.Any("error", err))
		return nil, err
	}

	return &authv1.FindUserByIDResponse{
		Name:   result.Name,
		Email:  result.Email,
		Id:     result.ID.String(),
		RoleId: result.RoleID.String(),
	}, nil
}

func (u *AuthTransport) FindRoleByID(ctx context.Context, req *authv1.FindRoleByIDRequest) (*authv1.FindRoleByIDResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing find role by id request", slog.String("request_id", reqID), slog.String("role_id", req.Id))

	result, err := u.ruc.FindByID(ctx, input.FindByIDRoleInput{ID: req.Id})
	if err != nil {
		u.log.Error("failed to find role", slog.String("request_id", reqID), slog.String("role_id", req.Id), slog.Any("error", err))
		return nil, err
	}

	var desc string
	if result.Description != nil {
		desc = *result.Description
	}

	perm := make([]*authv1.Permission, len(result.Permissions))
	for i, v := range result.Permissions {
		perm[i] = &authv1.Permission{
			Action: v.Action,
			Path:   v.Path,
		}
	}

	return &authv1.FindRoleByIDResponse{
		Id:          result.ID.String(),
		Description: desc,
		Name:        result.Name,
		Permissions: perm,
	}, nil
}

func (u *AuthTransport) FindAllRoles(ctx context.Context, req *authv1.FindAllRolesRequest) (*authv1.FindAllRolesResponse, error) {
	reqID := logger.ExtractRequestID(ctx)
	u.log.Info("processing find all roles request", slog.String("request_id", reqID))

	result, err := u.ruc.FindAll(ctx)
	if err != nil {
		u.log.Error("failed to find all roles", slog.String("request_id", reqID), slog.Any("error", err))
		return nil, err
	}

	roles := make([]*authv1.Role, len(result))
	for i, role := range result {
		permissions := make([]*authv1.Permission, len(role.Permissions))
		for j, permission := range role.Permissions {
			permissions[j] = &authv1.Permission{
				Action: permission.Action,
				Path:   permission.Path,
			}
		}

		var desc string
		if role.Description != nil {
			desc = *role.Description
		}

		roles[i] = &authv1.Role{
			Id:          role.ID.String(),
			Name:        role.Name,
			Description: desc,
			Permissions: permissions,
		}
	}

	return &authv1.FindAllRolesResponse{
		Roles: roles,
	}, nil
}
