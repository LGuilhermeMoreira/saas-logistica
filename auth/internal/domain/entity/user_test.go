package entity_test

import (
	"auth/internal/domain/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	validRoleID := uuid.New().String()

	t.Run("creates user successfully", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "john@example.com", "password123", validRoleID)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "John Smith", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, validRoleID, user.RoleID.String())

		// password must be hashed, never stored in plain text
		assert.NotEqual(t, "password123", user.Password)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err)
	})

	t.Run("trims whitespace from name, email and password", func(t *testing.T) {
		user, err := entity.NewUser("  John Smith  ", "  john@example.com  ", "  password123  ", validRoleID)

		require.NoError(t, err)
		assert.Equal(t, "John Smith", user.Name)
		assert.Equal(t, "john@example.com", user.Email)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err)
	})

	t.Run("generates different IDs for different users", func(t *testing.T) {
		user1, err := entity.NewUser("User 1", "user1@example.com", "password123", validRoleID)
		require.NoError(t, err)

		user2, err := entity.NewUser("User 2", "user2@example.com", "password123", validRoleID)
		require.NoError(t, err)

		assert.NotEqual(t, user1.ID, user2.ID)
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		user, err := entity.NewUser("", "john@example.com", "password123", validRoleID)

		assert.Nil(t, user)
		assert.EqualError(t, err, "client must have a name")
	})

	t.Run("returns error when name is only whitespace", func(t *testing.T) {
		user, err := entity.NewUser("   ", "john@example.com", "password123", validRoleID)

		assert.Nil(t, user)
		assert.EqualError(t, err, "client must have a name")
	})

	t.Run("returns error when email is empty", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "", "password123", validRoleID)

		assert.Nil(t, user)
		assert.EqualError(t, err, "email is invalid")
	})

	t.Run("returns error when email is invalid", func(t *testing.T) {
		invalidEmails := []string{
			"not-an-email",
			"missing-at-sign.com",
			"@no-user.com",
			"user@",
		}

		for _, email := range invalidEmails {
			t.Run(email, func(t *testing.T) {
				user, err := entity.NewUser("John Smith", email, "password123", validRoleID)

				assert.Nil(t, user)
				assert.EqualError(t, err, "email is invalid")
			})
		}
	})

	t.Run("returns error when password is empty", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "john@example.com", "", validRoleID)

		assert.Nil(t, user)
		assert.EqualError(t, err, "password does not meet the security requirements")
	})

	t.Run("returns error when password has fewer than 6 characters", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "john@example.com", "12345", validRoleID)

		assert.Nil(t, user)
		assert.EqualError(t, err, "password does not meet the security requirements")
	})

	t.Run("accepts password with exactly 6 characters", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "john@example.com", "123456", validRoleID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("returns error when roleID is not a valid uuid", func(t *testing.T) {
		user, err := entity.NewUser("John Smith", "john@example.com", "password123", "invalid-id")

		assert.Nil(t, user)
		assert.Error(t, err)
	})
}

func TestUser_ToMap(t *testing.T) {
	roleID := uuid.New()
	role := entity.Role{Name: "admin"}
	role.ID = roleID

	user := entity.User{
		ID:       uuid.New(),
		Name:     "John Smith",
		Email:    "john@example.com",
		RoleID:   roleID,
		Role:     role,
		Password: "hash-should-not-appear",
	}

	result := user.ToMap()

	assert.Equal(t, user.ID, result["id"])
	assert.Equal(t, "John Smith", result["name"])
	assert.Equal(t, "john@example.com", result["email"])
	assert.Equal(t, roleID, result["role_id"])
	assert.Equal(t, "admin", result["role_name"])
	assert.Equal(t, user.CreatedAt, result["created_at"])
	assert.Equal(t, user.UpdatedAt, result["updated_at"])
	assert.Equal(t, user.DeletedAt, result["deleted_at"])

	// password must never be exposed in the map
	_, exists := result["password"]
	assert.False(t, exists)
}
