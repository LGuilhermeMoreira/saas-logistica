package repository_test

import (
	"auth/internal/domain/entity"
	"auth/internal/infra/repository"
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "failed to connect to in-memory database")

	err = db.AutoMigrate(&entity.User{}, &entity.Role{}, &entity.Permission{})
	require.NoError(t, err, "failed to migrate entities")

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	userID, _ := uuid.NewV7()
	user := &entity.User{
		ID:    userID,
		Email: "test@example.com",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	var savedUser entity.User
	err = db.First(&savedUser, "id = ?", user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, savedUser.ID)
	assert.Equal(t, user.Email, savedUser.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	userID, _ := uuid.NewV7()
	user := &entity.User{
		ID:    userID,
		Email: "delete@example.com",
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	err = repo.Delete(ctx, user)
	assert.NoError(t, err)

	var count int64
	db.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestUserRepository_Login(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	perm, _ := entity.NewPermission("read", "/dashboard")
	role, _ := entity.NewRole("Admin", "", []entity.Permission{*perm})

	userID, _ := uuid.NewV7()
	expectedUser := &entity.User{
		ID:    userID,
		Email: "login@example.com",
		Role:  *role,
	}

	err := db.Create(expectedUser).Error
	require.NoError(t, err)

	t.Run("success logging in and fetching roles", func(t *testing.T) {
		result, err := repo.Login(ctx, "login@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)

		// Verifies if Preload("Roles") worked
		assert.NotNil(t, result.Role)
		assert.Equal(t, "Admin", result.Role.Name)
	})

	t.Run("error fetching non-existent email", func(t *testing.T) {
		result, err := repo.Login(ctx, "nonexistent@example.com")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find User by email")
		assert.Nil(t, result)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupUserTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	userID, _ := uuid.NewV7()
	expectedUser := &entity.User{
		ID:    userID,
		Email: "findbyid@example.com",
	}

	err := db.Create(expectedUser).Error
	require.NoError(t, err)

	t.Run("success finding user", func(t *testing.T) {
		result, err := repo.FindByID(ctx, expectedUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
	})

	t.Run("error fetching non-existent ID", func(t *testing.T) {
		randomID, _ := uuid.NewV7()
		result, err := repo.FindByID(ctx, randomID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find User by id")
		assert.Nil(t, result)
	})
}
