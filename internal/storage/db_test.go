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

func TestListAndGetExpression(t *testing.T) {
	os.Remove("test.db")
	db, err := NewStorage("test.db")
	assert.NoError(t, err)
	defer os.Remove("test.db")

	assert.NoError(t, db.RegisterUser("u", "p"))
	user, err := db.GetUserByLogin("u")
	assert.NoError(t, err)

	assert.NoError(t, db.SaveExpression(user.ID, "2+2", 4, ""))
	assert.NoError(t, db.SaveExpression(user.ID, "10/0", 0, "division by zero"))

	list, err := db.ListExpressions(user.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	assert.Equal(t, "2+2", list[0].Expression)
	assert.Equal(t, 4.0, list[0].Result)
	assert.Empty(t, list[0].Error)

	assert.Equal(t, "10/0", list[1].Expression)
	assert.Equal(t, 0.0, list[1].Result)
	assert.Equal(t, "division by zero", list[1].Error)

	rec, err := db.GetExpression(user.ID, list[1].ID)
	assert.NoError(t, err)
	assert.Equal(t, list[1].Expression, rec.Expression)

	_, err = db.GetExpression(user.ID, 9999)
	assert.Equal(t, ErrNotFound, err)
}
