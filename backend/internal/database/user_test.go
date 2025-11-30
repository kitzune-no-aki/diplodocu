package database

import (
	"testing"

	"github.com/kitzune-no-aki/diplodocu/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Webuser{})
	require.NoError(t, err)

	return db
}

func TestSyncUser_CreatesNewUser(t *testing.T) {
	db := setupTestDB(t)

	user, err := SyncUser(db, "keycloak-user-123", "John Doe")

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "keycloak-user-123", user.ID)
	assert.Equal(t, "John Doe", *user.Name)

	// Verify user exists in DB
	var count int64
	db.Model(&models.Webuser{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSyncUser_UpdatesExistingUser(t *testing.T) {
	db := setupTestDB(t)

	// Create initial user
	_, err := SyncUser(db, "keycloak-user-123", "Old Name")
	require.NoError(t, err)

	// Sync again with new name
	user, err := SyncUser(db, "keycloak-user-123", "New Name")

	require.NoError(t, err)
	assert.Equal(t, "New Name", *user.Name)

	// Verify still only one user
	var count int64
	db.Model(&models.Webuser{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSyncUser_HandlesEmptyName(t *testing.T) {
	db := setupTestDB(t)

	user, err := SyncUser(db, "keycloak-user-123", "")

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "keycloak-user-123", user.ID)
	assert.Nil(t, user.Name) // Empty string becomes nil
}

func TestSyncUser_IdempotentSync(t *testing.T) {
	db := setupTestDB(t)

	// Sync multiple times with same data
	_, err := SyncUser(db, "keycloak-user-123", "Same Name")
	require.NoError(t, err)

	_, err = SyncUser(db, "keycloak-user-123", "Same Name")
	require.NoError(t, err)

	_, err = SyncUser(db, "keycloak-user-123", "Same Name")
	require.NoError(t, err)

	// Should still be only one user
	var count int64
	db.Model(&models.Webuser{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSyncUser_MultipleUsers(t *testing.T) {
	db := setupTestDB(t)

	// Create multiple users
	user1, err := SyncUser(db, "user-1", "User One")
	require.NoError(t, err)

	user2, err := SyncUser(db, "user-2", "User Two")
	require.NoError(t, err)

	user3, err := SyncUser(db, "user-3", "User Three")
	require.NoError(t, err)

	// Verify all users created
	assert.Equal(t, "user-1", user1.ID)
	assert.Equal(t, "user-2", user2.ID)
	assert.Equal(t, "user-3", user3.ID)

	var count int64
	db.Model(&models.Webuser{}).Count(&count)
	assert.Equal(t, int64(3), count)
}

func TestSyncUser_PreservesIDOnUpdate(t *testing.T) {
	db := setupTestDB(t)

	// Create user
	originalUser, err := SyncUser(db, "keycloak-user-123", "Original")
	require.NoError(t, err)

	// Update user
	updatedUser, err := SyncUser(db, "keycloak-user-123", "Updated")
	require.NoError(t, err)

	// ID should remain the same
	assert.Equal(t, originalUser.ID, updatedUser.ID)
	assert.Equal(t, "Updated", *updatedUser.Name)
}
