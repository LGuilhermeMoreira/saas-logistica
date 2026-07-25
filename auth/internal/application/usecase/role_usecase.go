package usecase

import (
	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/authorization"
	"context"
)

type RoleUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error)

	Delete(ctx context.Context, dto input.DeleteRoleInput) error
}

type RoleUsecase struct {
	repository contract.RoleRepositoryInterface
	casbin     authorization.PermissionEnforcer
}

func NewRoleUsecase(repo contract.RoleRepositoryInterface, casbin authorization.PermissionEnforcer) RoleUsecaseInterface {
	return &RoleUsecase{
		repository: repo,
		casbin:     casbin,
	}
}

func (r *RoleUsecase) Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error) {
	role, err := entity.NewRole(dto.Name, dto.Descritpion)
	if err != nil {
		return nil, err
	}

	err = r.repository.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return &output.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}, nil

}

func (r *RoleUsecase) Delete(ctx context.Context, dto input.DeleteRoleInput) error
