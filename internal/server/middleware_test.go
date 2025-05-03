package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	authMiddleware := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	authMiddleware.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
