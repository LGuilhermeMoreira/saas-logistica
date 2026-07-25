package usecase

import (
	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/authentication"
	"auth/pkg/authorization"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateUserInput) (*output.CreateUserOutput, error)
	Delete(ctx context.Context, dto input.DeleteUserInput) error
	Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error)
}

type UserUsecase struct {
	repository contract.UserRepositoryInterface
	jwt        authentication.JWTInterface
	casbin     authorization.RoleAssigner
}

func NewUserUsecase(repo contract.UserRepositoryInterface, jwt authentication.JWTInterface, casbin authorization.RoleAssigner) UserUsecaseInterface {
	return &UserUsecase{repository: repo, jwt: jwt, casbin: casbin}
}
func (c *UserUsecase) Create(ctx context.Context, dto input.CreateUserInput) (*output.CreateUserOutput, error) {
	result, err := entity.NewUser(dto.Name, dto.Email, dto.Password, dto.RoleID)
	if err != nil {
		return nil, err
	}

	err = c.repository.Create(ctx, result)
	if err != nil {
		return nil, err
	}

	err = c.casbin.AssignRoleToUser(result.ID.String(), result.RoleID.String())
	if err != nil {
		return nil, err
	}

	return &output.CreateUserOutput{
		ID:     result.ID,
		RoleID: result.RoleID,
		Email:  result.Email,
		Name:   result.Name,
	}, nil
}
func (c *UserUsecase) Delete(ctx context.Context, dto input.DeleteUserInput) error {
	ID, err := uuid.Parse(dto.ID)
	if err != nil {
		return err
	}

	result, err := c.repository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	return c.repository.Delete(ctx, result)
}
func (c *UserUsecase) Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error) {
	if strings.TrimSpace(dto.Email) == "" || strings.TrimSpace(dto.Password) == "" {
		return nil, errors.New("invalid credentials")
	}

	result, err := c.repository.Login(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}

	token, err := c.jwt.GenerateToken(result.ToMap())
	if err != nil {
		return nil, err
	}

	return &output.LoginOutput{
		Token: string(token),
	}, nil
}
