package usecase

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"auth/internal/application/input"
	"auth/internal/application/output"
	"auth/internal/domain/contract"
	"auth/internal/domain/entity"
	"auth/pkg/authentication"
	"auth/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseInterface interface {
	Create(ctx context.Context, dto input.CreateUserInput) (*output.UserOutput, error)
	Delete(ctx context.Context, dto input.DeleteUserInput) error
	Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error)
	FindByID(ctx context.Context, dto input.FindByIDUserInput) (*output.UserOutput, error)
}

type UserUsecase struct {
	repository contract.UserRepositoryInterface
	jwt        authentication.TokenGenerator
	log        *slog.Logger
}

func NewUserUsecase(repo contract.UserRepositoryInterface, jwt authentication.TokenGenerator, log *slog.Logger) UserUsecaseInterface {
	return &UserUsecase{repository: repo, jwt: jwt, log: log}
}

func (c *UserUsecase) Create(ctx context.Context, dto input.CreateUserInput) (*output.UserOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Info("starting user creation", slog.String("request_id", reqID), slog.String("email", dto.Email))

	result, err := entity.NewUser(dto.Name, dto.Email, dto.Password, dto.RoleID)
	if err != nil {
		c.log.Error("failed to create user entity", slog.String("request_id", reqID), slog.String("email", dto.Email), slog.Any("error", err))
		return nil, err
	}

	err = c.repository.Create(ctx, result)
	if err != nil {
		c.log.Error("failed to save user in repository", slog.String("request_id", reqID), slog.String("email", dto.Email), slog.Any("error", err))
		return nil, err
	}

	c.log.Info("user created successfully", slog.String("request_id", reqID), slog.String("user_id", result.ID.String()))

	return &output.UserOutput{
		ID:     result.ID,
		RoleID: result.RoleID,
		Email:  result.Email,
		Name:   result.Name,
	}, nil
}

func (c *UserUsecase) Delete(ctx context.Context, dto input.DeleteUserInput) error {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Info("starting user deletion", slog.String("request_id", reqID), slog.String("user_id", dto.ID))

	ID, err := uuid.Parse(dto.ID)
	if err != nil {
		c.log.Error("invalid user id format", slog.String("request_id", reqID), slog.String("user_id", dto.ID), slog.Any("error", err))
		return err
	}

	result, err := c.repository.FindByID(ctx, ID)
	if err != nil {
		c.log.Error("failed to find user for deletion", slog.String("request_id", reqID), slog.String("user_id", dto.ID), slog.Any("error", err))
		return err
	}

	err = c.repository.Delete(ctx, result)
	if err != nil {
		c.log.Error("failed to delete user in repository", slog.String("request_id", reqID), slog.String("user_id", dto.ID), slog.Any("error", err))
		return err
	}

	c.log.Info("user deleted successfully", slog.String("request_id", reqID), slog.String("user_id", dto.ID))
	return nil
}

func (c *UserUsecase) Login(ctx context.Context, dto input.LoginInput) (*output.LoginOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Info("starting user login process", slog.String("request_id", reqID), slog.String("email", dto.Email))

	if strings.TrimSpace(dto.Email) == "" || strings.TrimSpace(dto.Password) == "" {
		err := errors.New("invalid credentials")
		c.log.Warn("empty credentials provided", slog.String("request_id", reqID), slog.String("email", dto.Email))
		return nil, err
	}

	result, err := c.repository.Login(ctx, dto.Email)
	if err != nil {
		c.log.Error("failed to fetch user for login", slog.String("request_id", reqID), slog.String("email", dto.Email), slog.Any("error", err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(dto.Password))
	if err != nil {
		c.log.Warn("invalid password attempt", slog.String("request_id", reqID), slog.String("email", dto.Email))
		return nil, err
	}

	token, err := c.jwt.GenerateToken(result.ToMap())
	if err != nil {
		c.log.Error("failed to generate jwt token", slog.String("request_id", reqID), slog.String("email", dto.Email), slog.Any("error", err))
		return nil, err
	}

	c.log.Info("user logged in successfully", slog.String("request_id", reqID), slog.String("email", dto.Email))

	return &output.LoginOutput{
		Token: string(token),
	}, nil
}

func (c *UserUsecase) FindByID(ctx context.Context, dto input.FindByIDUserInput) (*output.UserOutput, error) {
	reqID := logger.ExtractRequestID(ctx)
	c.log.Info("starting find user by id", slog.String("request_id", reqID), slog.String("user_id", dto.ID))

	id, err := uuid.Parse(dto.ID)
	if err != nil {
		c.log.Error("invalid user id format", slog.String("request_id", reqID), slog.String("user_id", dto.ID), slog.Any("error", err))
		return nil, err
	}

	model, err := c.repository.FindByID(ctx, id)
	if err != nil {
		c.log.Error("failed to find user by id", slog.String("request_id", reqID), slog.String("user_id", dto.ID), slog.Any("error", err))
		return nil, err
	}

	c.log.Info("user found successfully", slog.String("request_id", reqID), slog.String("user_id", dto.ID))

	return &output.UserOutput{
		Name:   model.Name,
		Email:  model.Email,
		ID:     model.ID,
		RoleID: model.RoleID,
	}, nil
}
