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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to the database.")

	err = db.AutoMigrate(&entity.Role{}, &entity.Permission{})
	require.NoError(t, err, "failed to migrate entities")

	return db
}

func TestRoleRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewRoleRepository(db)
	ctx := context.Background()

	perm, _ := entity.NewPermission("read", "/users")
	role, err := entity.NewRole("Admin", "Administrator role", []entity.Permission{*perm})
	require.NoError(t, err)

	err = repo.Create(ctx, role)
	assert.NoError(t, err)

	var savedRole entity.Role
	err = db.Preload("Permissions").First(&savedRole, "id = ?", role.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, role.ID, savedRole.ID)
	assert.Equal(t, role.Name, savedRole.Name)
	assert.Len(t, savedRole.Permissions, 1)
}

func TestRoleRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewRoleRepository(db)
	ctx := context.Background()

	perm, _ := entity.NewPermission("write", "/posts")
	expectedRole, _ := entity.NewRole("Editor", "Editor role", []entity.Permission{*perm})

	err := db.Create(expectedRole).Error
	require.NoError(t, err)

	t.Run("success finding existing role", func(t *testing.T) {
		result, err := repo.FindByID(ctx, expectedRole.ID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedRole.ID, result.ID)
	})

	t.Run("Error retrieving non-existent ID", func(t *testing.T) {
		randomID := uuid.New()
		result, err := repo.FindByID(ctx, randomID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, result)
	})
}

func TestRoleRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewRoleRepository(db)
	ctx := context.Background()

	db.Exec("DELETE FROM roles")

	perm1, _ := entity.NewPermission("read", "/a")
	role1, _ := entity.NewRole("Role A", "", []entity.Permission{*perm1})

	perm2, _ := entity.NewPermission("read", "/b")
	role2, _ := entity.NewRole("Role B", "", []entity.Permission{*perm2})

	err := db.Create([]*entity.Role{role1, role2}).Error
	require.NoError(t, err)

	result, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRoleRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewRoleRepository(db)
	ctx := context.Background()

	perm, _ := entity.NewPermission("delete", "/all")
	role, _ := entity.NewRole("SuperAdmin", "", []entity.Permission{*perm})

	err := db.Create(role).Error
	require.NoError(t, err)

	err = repo.Delete(ctx, role)
	assert.NoError(t, err)

	var count int64
	db.Model(&entity.Role{}).Where("id = ?", role.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
