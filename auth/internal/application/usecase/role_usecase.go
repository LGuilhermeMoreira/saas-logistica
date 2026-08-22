package usecase

import (
	"context"
	"log/slog"

	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/logger"

	"github.com/google/uuid"
)

type RoleUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error)
	Delete(ctx context.Context, dto input.DeleteRoleInput) error
	FindByID(ctx context.Context, dto input.FindByIDRoleInput) (*output.RoleOutput, error)
	FindAll(ctx context.Context) ([]output.RoleOutput, error)
}

type RoleUsecase struct {
	repository contract.RoleRepositoryInterface
	opaService contract.OPASyncServiceInterface
	log        *slog.Logger
}

func NewRoleUsecase(repo contract.RoleRepositoryInterface, opa contract.OPASyncServiceInterface, log *slog.Logger) RoleUsecaseInterface {
	return &RoleUsecase{
		repository: repo,
		opaService: opa,
		log:        log,
	}
}

func (r *RoleUsecase) Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Info("starting role creation", slog.String("request_id", reqID), slog.String("role_name", dto.Name))

	perms := make([]entity.Permission, len(dto.Permissions))

	for i, v := range dto.Permissions {
		perm, err := entity.NewPermission(v.Action, v.Path)
		if err != nil {
			r.log.Error("failed to create permission entity", slog.String("request_id", reqID), slog.Any("error", err))
			return nil, err
		}
		perms[i] = *perm
	}

	role, err := entity.NewRole(dto.Name, dto.Description, perms)
	if err != nil {
		r.log.Error("failed to create role entity", slog.String("request_id", reqID), slog.Any("error", err))
		return nil, err
	}

	err = r.repository.Create(ctx, role)
	if err != nil {
		r.log.Error("failed to save role in repository", slog.String("request_id", reqID), slog.Any("error", err))
		return nil, err
	}

	err = r.opaService.SyncPolicies(ctx)
	if err != nil {
		r.log.Error("failed to sync policies with OPA", slog.String("request_id", reqID), slog.String("role_id", role.ID.String()), slog.Any("error", err))
		return nil, err
	}

	permOutput := make([]output.PermissionsOutput, len(perms))

	for i := 0; i < len(perms); i++ {
		permOutput[i] = output.PermissionsOutput{
			Action: perms[i].Action,
			Path:   perms[i].Path,
		}
	}

	r.log.Info("role created successfully", slog.String("request_id", reqID), slog.String("role_id", role.ID.String()))

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permOutput,
	}, nil
}

func (r *RoleUsecase) Delete(ctx context.Context, dto input.DeleteRoleInput) error {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Info("starting role deletion", slog.String("request_id", reqID), slog.String("role_id", dto.ID))

	id, err := uuid.Parse(dto.ID)
	if err != nil {
		r.log.Error("invalid role id format", slog.String("request_id", reqID), slog.String("role_id", dto.ID), slog.Any("error", err))
		return err
	}

	result, err := r.repository.FindByID(ctx, id)
	if err != nil {
		r.log.Error("failed to find role for deletion", slog.String("request_id", reqID), slog.String("role_id", dto.ID), slog.Any("error", err))
		return err
	}

	err = r.repository.Delete(ctx, result)
	if err != nil {
		r.log.Error("failed to delete role in repository", slog.String("request_id", reqID), slog.String("role_id", dto.ID), slog.Any("error", err))
		return err
	}

	r.log.Info("role deleted successfully", slog.String("request_id", reqID), slog.String("role_id", dto.ID))
	return nil
}

func (r *RoleUsecase) FindByID(ctx context.Context, dto input.FindByIDRoleInput) (*output.RoleOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Info("starting find role by id", slog.String("request_id", reqID), slog.String("role_id", dto.ID))

	id, err := uuid.Parse(dto.ID)
	if err != nil {
		r.log.Error("invalid role id format", slog.String("request_id", reqID), slog.String("role_id", dto.ID), slog.Any("error", err))
		return nil, err
	}

	role, err := r.repository.FindByID(ctx, id)
	if err != nil {
		r.log.Error("failed to find role by id", slog.String("request_id", reqID), slog.String("role_id", dto.ID), slog.Any("error", err))
		return nil, err
	}

	permOutput := make([]output.PermissionsOutput, len(role.Permissions))

	for i, permission := range role.Permissions {
		permOutput[i] = output.PermissionsOutput{
			Action: permission.Action,
			Path:   permission.Path,
		}
	}

	r.log.Info("role found successfully", slog.String("request_id", reqID), slog.String("role_id", dto.ID))

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permOutput,
	}, nil
}

func (r *RoleUsecase) FindAll(ctx context.Context) ([]output.RoleOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	r.log.Info("starting find all roles", slog.String("request_id", reqID))

	roles, err := r.repository.FindAll(ctx)
	if err != nil {
		r.log.Error("failed to find all roles", slog.String("request_id", reqID), slog.Any("error", err))
		return nil, err
	}

	roleOutputs := make([]output.RoleOutput, len(roles))

	for i, role := range roles {
		permOutputs := make([]output.PermissionsOutput, len(role.Permissions))

		for j, permission := range role.Permissions {
			permOutputs[j] = output.PermissionsOutput{
				Action: permission.Action,
				Path:   permission.Path,
			}
		}

		roleOutputs[i] = output.RoleOutput{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permOutputs,
		}
	}

	r.log.Info("roles retrieved successfully", slog.String("request_id", reqID), slog.Int("count", len(roles)))

	return roleOutputs, nil
}
