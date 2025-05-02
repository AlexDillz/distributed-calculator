package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateExpression(t *testing.T) {
	result, err := EvaluateExpression("2+2*2")
	assert.NoError(t, err)
	assert.Equal(t, 6.0, result)
}
