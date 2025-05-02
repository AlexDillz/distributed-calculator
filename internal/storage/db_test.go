package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *Storage {
	db, err := NewStorage("test.db")
	assert.NoError(t, err)
	return db
}

func teardownDB() {
	os.Remove("test.db")
}

func TestRegisterUser(t *testing.T) {
	db := setupDB(t)
	defer teardownDB()

	err := db.RegisterUser("testuser", "password123")
	assert.NoError(t, err)

	user, err := db.GetUserByLogin("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Login)
}
