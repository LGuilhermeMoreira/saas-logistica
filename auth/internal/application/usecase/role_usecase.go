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
	var perms []entity.Permission

	for _, v := range dto.Permissions {
		perm, err := entity.NewPermission(v.Action, v.Path)
		if err != nil {
			return nil, err
		}

		perms = append(perms, *perm)
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

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
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
