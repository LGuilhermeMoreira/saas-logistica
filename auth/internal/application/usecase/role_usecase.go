package application

import (
	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"context"
)

type RoleUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateRoleInput) (*output.RoleOutput, error)

	Delete(ctx context.Context, dto input.DeleteRoleInput) error
}

type RoleUsecase struct {
	repository contract.RoleRepositoryInterface
}

func NewRoleUsecase(repo contract.RoleRepositoryInterface) RoleUsecaseInterface {
	return &RoleUsecase{
		repository: repo,
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
