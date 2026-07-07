package application

import (
	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/authentication"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ClientUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateClientInput) (*output.CreateClientOutput, error)
	Delete(ctx context.Context, dto input.DeleteClientInput) error
	Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error)
}

type ClientUsecase struct {
	repository contract.ClientRepositoryInterface
	jwt        authentication.JWTInterface
}

func NewClientUsecase(repo contract.ClientRepositoryInterface, jwt authentication.JWTInterface) ClientUsecaseInterface {
	return &ClientUsecase{repository: repo, jwt: jwt}
}
func (c *ClientUsecase) Create(ctx context.Context, dto input.CreateClientInput) (*output.CreateClientOutput, error) {
	cliente, err := entity.NewClient(dto.Name, dto.Email, dto.Password, dto.RoleID)
	if err != nil {
		return nil, err
	}

	err = c.repository.Create(ctx, cliente)
	if err != nil {
		return nil, err
	}

	return &output.CreateClientOutput{
		ID:     cliente.ID,
		RoleID: cliente.RoleID,
		Email:  cliente.Email,
		Name:   cliente.Name,
	}, nil
}
func (c *ClientUsecase) Delete(ctx context.Context, dto input.DeleteClientInput) error {
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
func (c *ClientUsecase) Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error) {
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
