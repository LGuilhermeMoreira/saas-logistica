package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"auth/internal/application/input"
	"auth/internal/application/usecase"
	"auth/internal/domain/entity"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, model *entity.Role) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, model *entity.Role) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Role), args.Error(1)
}

type MockOPASyncService struct {
	mock.Mock
}

func (m *MockOPASyncService) SyncPolicies(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRoleUsecase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a role with success", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.CreateRoleInput{
			Name:        "Admin",
			Description: "Administrador do sistema",
			Permissions: []input.PermissionInput{
				{Action: "POST", Path: "/alguma-coisa"},
			},
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Role")).Return(nil)
		mockOPA.On("SyncPolicies", ctx).Return(nil)

		result, err := uc.Create(ctx, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Admin", result.Name)
		assert.NotNil(t, result.Description)
		assert.Equal(t, "Administrador do sistema", *result.Description)
		assert.NotEqual(t, uuid.Nil, result.ID)

		mockRepo.AssertExpectations(t)
		mockOPA.AssertExpectations(t)
	})

	t.Run("should return error when name is empty (Entity Rule)", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.CreateRoleInput{
			Name:        "   ",
			Description: "Administrador do sistema",
			Permissions: []input.PermissionInput{
				{Action: "POST", Path: "/alguma-coisa"},
			},
		}

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "name is invalid", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
		mockOPA.AssertNotCalled(t, "SyncPolicies")
	})

	t.Run("should return error when permissions are empty (Entity Rule)", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.CreateRoleInput{
			Name:        "Admin",
			Description: "Administrador do sistema",
			Permissions: []input.PermissionInput{},
		}

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "role must have one permission", err.Error())

		mockRepo.AssertNotCalled(t, "Create")
		mockOPA.AssertNotCalled(t, "SyncPolicies")
	})

	t.Run("should return error when repository fails to persist", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.CreateRoleInput{
			Name:        "Admin",
			Description: "Administrador do sistema",
			Permissions: []input.PermissionInput{
				{Action: "POST", Path: "/alguma-coisa"},
			},
		}

		expectedErr := errors.New("database error")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Role")).Return(expectedErr)

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		mockOPA.AssertNotCalled(t, "SyncPolicies")
	})

	t.Run("should return error when OPA Sync fails", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.CreateRoleInput{
			Name:        "Admin",
			Description: "Administrador do sistema",
			Permissions: []input.PermissionInput{
				{Action: "POST", Path: "/alguma-coisa"},
			},
		}

		expectedErr := errors.New("opa sync error")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Role")).Return(nil)
		mockOPA.On("SyncPolicies", ctx).Return(expectedErr)

		result, err := uc.Create(ctx, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
		mockOPA.AssertExpectations(t)
	})
}

func TestRoleUsecase_Delete(t *testing.T) {
	ctx := context.Background()
	validUUID := uuid.New()

	t.Run("should delete the role with success", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.DeleteRoleInput{ID: validUUID.String()}
		fakeRole := &entity.Role{ID: validUUID, Name: "Admin"}

		mockRepo.On("FindByID", ctx, validUUID).Return(fakeRole, nil)
		mockRepo.On("Delete", ctx, fakeRole).Return(nil)

		err := uc.Delete(ctx, dto)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error if uuid is invalid", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.DeleteRoleInput{ID: "invalid-uuid"}

		err := uc.Delete(ctx, dto)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "FindByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("should return an error if role is not found", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.DeleteRoleInput{ID: validUUID.String()}
		expectedErr := errors.New("role not found")

		mockRepo.On("FindByID", ctx, validUUID).Return(nil, expectedErr)

		err := uc.Delete(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("should return an error if the repository fails to delete", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		mockOPA := new(MockOPASyncService)
		uc := usecase.NewRoleUsecase(mockRepo, mockOPA)

		dto := input.DeleteRoleInput{ID: validUUID.String()}
		fakeRole := &entity.Role{ID: validUUID, Name: "Admin"}
		expectedErr := errors.New("delete error")

		mockRepo.On("FindByID", ctx, validUUID).Return(fakeRole, nil)
		mockRepo.On("Delete", ctx, fakeRole).Return(expectedErr)

		err := uc.Delete(ctx, dto)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
