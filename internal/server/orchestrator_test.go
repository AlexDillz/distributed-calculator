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

func TestEvaluateExpressionErrors(t *testing.T) {
	_, err := EvaluateExpression("2++2")
	assert.Error(t, err)

	_, err = EvaluateExpression("5/0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "division by zero")
}
