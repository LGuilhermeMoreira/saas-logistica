package usecase

import (
	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"

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
}

func NewRoleUsecase(repo contract.RoleRepositoryInterface, opa contract.OPASyncServiceInterface) RoleUsecaseInterface {
	return &RoleUsecase{
		repository: repo,
		opaService: opa,
	}
}

func (r *RoleUsecase) Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error) {
	perms := make([]entity.Permission, len(dto.Permissions))

	for i, v := range dto.Permissions {
		perm, err := entity.NewPermission(v.Action, v.Path)
		if err != nil {
			return nil, err
		}

		perms[i] = *perm
	}

	role, err := entity.NewRole(dto.Name, dto.Description, perms)
	if err != nil {
		return nil, err
	}

	err = r.repository.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	err = r.opaService.SyncPolicies(ctx)
	if err != nil {
		return nil, err
	}

	permOutput := make([]output.PermissionsOutput, len(perms))

	for i := 0; i < len(perms); i++ {
		permOutput[i] = output.PermissionsOutput{
			Action: perms[i].Action,
			Path:   perms[i].Path,
		}
	}

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permOutput,
	}, nil

}

func (r *RoleUsecase) Delete(ctx context.Context, dto input.DeleteRoleInput) error {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return err
	}

	result, err := r.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return r.repository.Delete(ctx, result)

}

func (r *RoleUsecase) FindByID(
	ctx context.Context,
	dto input.FindByIDRoleInput,
) (*output.RoleOutput, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, err
	}

	role, err := r.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	permOutput := make([]output.PermissionsOutput, len(role.Permissions))

	for i, permission := range role.Permissions {
		permOutput[i] = output.PermissionsOutput{
			Action: permission.Action,
			Path:   permission.Path,
		}
	}

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permOutput,
	}, nil
}

func (r *RoleUsecase) FindAll(ctx context.Context) ([]output.RoleOutput, error) {
	roles, err := r.repository.FindAll(ctx)
	if err != nil {
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

	return roleOutputs, nil
}
