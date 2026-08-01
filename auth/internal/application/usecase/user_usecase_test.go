package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"auth/internal/application/input"
	"auth/internal/application/usecase"
	"auth/internal/domain/entity"
	"auth/pkg/authentication"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, model *entity.User) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, model *entity.User) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockUserRepository) Login(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(data any) (authentication.JWTToken, error) {
	args := m.Called(data)
	return args.Get(0).(authentication.JWTToken), args.Error(1)
}

func (m *MockTokenGenerator) GenerateShortToken(data any) (authentication.JWTToken, error) {
	args := m.Called(data)
	return args.Get(0).(authentication.JWTToken), args.Error(1)
}

func TestUserUsecase_Create(t *testing.T) {
	ctx := context.Background()
	validRoleID := uuid.New().String()

	t.Run("should create a user with success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.CreateUserInput{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Password: "securepassword",
			RoleID:   validRoleID,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		result, err := uc.Create(ctx, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "John Doe", result.Name)
		assert.Equal(t, "john.doe@example.com", result.Email)
		assert.Equal(t, validRoleID, result.RoleID.String())
		assert.NotEqual(t, uuid.Nil, result.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when name is empty (Entity Rule)", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.CreateUserInput{
			Name:     "   ",
			Email:    "john.doe@example.com",
			Password: "securepassword",
			RoleID:   validRoleID,
		}

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "client must have a name", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when email is invalid (Entity Rule)", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.CreateUserInput{
			Name:     "John Doe",
			Email:    "invalid-email",
			Password: "securepassword",
			RoleID:   validRoleID,
		}

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email is invalid", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when password is too short (Entity Rule)", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "123", // Menos que 6 caracteres
			RoleID:   validRoleID,
		}

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "password does not meet the security requirements", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.CreateUserInput{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "securepassword",
			RoleID:   validRoleID,
		}

		expectedErr := errors.New("database error")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(expectedErr)

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestUserUsecase_Delete(t *testing.T) {
	ctx := context.Background()
	validUUID := uuid.New()

	t.Run("should delete user with success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.DeleteUserInput{ID: validUUID.String()}
		fakeUser := &entity.User{ID: validUUID}

		mockRepo.On("FindByID", ctx, validUUID).Return(fakeUser, nil)
		mockRepo.On("Delete", ctx, fakeUser).Return(nil)

		err := uc.Delete(ctx, dto)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid uuid", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.DeleteUserInput{ID: "invalid-uuid"}
		err := uc.Delete(ctx, dto)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "FindByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.DeleteUserInput{ID: validUUID.String()}
		expectedErr := errors.New("not found")

		mockRepo.On("FindByID", ctx, validUUID).Return(nil, expectedErr)

		err := uc.Delete(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertNotCalled(t, "Delete")
	})
}

func TestUserUsecase_Login(t *testing.T) {
	ctx := context.Background()
	validEmail := "john@example.com"
	validPassword := "securepassword"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	t.Run("should login and generate token with success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		uc := usecase.NewUserUsecase(mockRepo, mockTokenGen)

		dto := input.LoginInput{
			Email:    validEmail,
			Password: validPassword,
		}

		fakeUser := &entity.User{
			ID:       uuid.New(),
			Name:     "John",
			Email:    validEmail,
			Password: string(hashedPassword),
		}

		mockRepo.On("Login", ctx, validEmail).Return(fakeUser, nil)
		mockTokenGen.On("GenerateToken", mock.AnythingOfType("map[string]interface {}")).Return(authentication.JWTToken("fake-jwt-token"), nil)

		result, err := uc.Login(ctx, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "fake-jwt-token", result.Token)

		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("should return error for empty credentials", func(t *testing.T) {
		uc := usecase.NewUserUsecase(nil, nil)

		dto := input.LoginInput{
			Email:    "   ",
			Password: "",
		}

		result, err := uc.Login(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("should return error if user not found on repository", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.LoginInput{Email: validEmail, Password: validPassword}
		expectedErr := errors.New("user not found")

		mockRepo.On("Login", ctx, validEmail).Return(nil, expectedErr)

		result, err := uc.Login(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error for wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil)

		dto := input.LoginInput{Email: validEmail, Password: "wrongpassword"}

		fakeUser := &entity.User{
			ID:       uuid.New(),
			Password: string(hashedPassword), // A senha real aqui é "securepassword"
		}

		mockRepo.On("Login", ctx, validEmail).Return(fakeUser, nil)

		result, err := uc.Login(ctx, dto)

		assert.Error(t, err) // Esperamos um erro do próprio pacote do bcrypt
		assert.Nil(t, result)
	})

	t.Run("should return error if token generation fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		uc := usecase.NewUserUsecase(mockRepo, mockTokenGen)

		dto := input.LoginInput{Email: validEmail, Password: validPassword}
		fakeUser := &entity.User{
			ID:       uuid.New(),
			Password: string(hashedPassword),
		}

		expectedErr := errors.New("failed to generate token")

		mockRepo.On("Login", ctx, validEmail).Return(fakeUser, nil)
		mockTokenGen.On("GenerateToken", mock.AnythingOfType("map[string]interface {}")).Return(authentication.JWTToken(""), expectedErr)

		result, err := uc.Login(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}
