package entity_test

import (
	"auth/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPermission(t *testing.T) {
	action := "read"
	path := "/users"
	permission, err := entity.NewPermission(action, path)

	assert.Nil(t, err)
	assert.NotNil(t, permission)
	assert.NotNil(t, permission.ID)
	assert.Equal(t, action, permission.Action)
	assert.Equal(t, path, permission.Path)
}

func TestNewRole(t *testing.T) {
	t.Run("should create a new role successfully", func(t *testing.T) {
		name := "Admin"
		description := "Administrator role"
		permissions := []entity.Permission{
			{Action: "create", Path: "/admin"},
			{Action: "read", Path: "/reports"},
		}
		role, err := entity.NewRole(name, description, permissions)

		assert.Nil(t, err)
		assert.NotNil(t, role)
		assert.NotNil(t, role.ID)
		assert.Equal(t, name, role.Name)
		assert.Equal(t, description, *role.Description)
		assert.Len(t, role.Permissions, 2)
		assert.Equal(t, role.ID, role.Permissions[0].RoleID)
		assert.Equal(t, role.ID, role.Permissions[1].RoleID)
	})

	t.Run("should create a new role with nil description if empty", func(t *testing.T) {
		name := "User"
		description := ""
		permissions := []entity.Permission{
			{Action: "read", Path: "/profile"},
		}
		role, err := entity.NewRole(name, description, permissions)

		assert.Nil(t, err)
		assert.NotNil(t, role)
		assert.Nil(t, role.Description)
	})

	t.Run("should trim name and description", func(t *testing.T) {
		name := "  Editor  "
		description := "  Editor role with extra spaces  "
		permissions := []entity.Permission{
			{Action: "edit", Path: "/documents"},
		}
		role, err := entity.NewRole(name, description, permissions)

		assert.Nil(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "Editor", role.Name)
		assert.Equal(t, "Editor role with extra spaces", *role.Description)
	})

	t.Run("should return error if name is empty", func(t *testing.T) {
		name := ""
		description := "Some description"
		permissions := []entity.Permission{
			{Action: "read", Path: "/anything"},
		}
		role, err := entity.NewRole(name, description, permissions)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "name is invalid")
		assert.Nil(t, role)
	})

	t.Run("should return error if permissions are empty", func(t *testing.T) {
		name := "Viewer"
		description := "Viewer role"
		permissions := []entity.Permission{}
		role, err := entity.NewRole(name, description, permissions)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "role must have one permission")
		assert.Nil(t, role)
	})
}
