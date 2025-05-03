package calculation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc_ValidExpressions(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"2+2*2", 6.0},
		{"(2+3)*4", 20.0},
		{"10*-2", -20.0},
		{"3.5+2.5*2", 8.5},
		{"2e3+500", 2500.0},
	}

	for _, tt := range tests {
		result, err := Calc(tt.expr)
		assert.NoError(t, err, "Для выражения %s ожидалось успешное вычисление", tt.expr)
		assert.Equal(t, tt.expected, result, "Выражение: %s", tt.expr)
	}
}

func TestCalc_InvalidExpressions(t *testing.T) {
	tests := []string{
		"2++2",
		"abc+def",
		"2+abc",
		"((2+3)",
		"2+",
		"2*",
		"2/",
		"2)",
	}

	for _, expr := range tests {
		_, err := Calc(expr)
		assert.Error(t, err, "Для выражения %s ожидалась ошибка", expr)
	}
}

func TestParseNegativeNumbers(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"-5+10", 5.0},
		{"10*-2", -20.0},
		{"-2*-3", 6.0},
		{"-(-5)", 5.0},
	}

	for _, tt := range tests {
		result, err := Calc(tt.expr)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, result, "Выражение: %s", tt.expr)
	}
}

func TestParseScientificNotation(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"2e3", 2000.0},
		{"5E2+2", 502.0},
		{"1e1*2e1", 200.0},
	}

	for _, tt := range tests {
		result, err := Calc(tt.expr)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, result)
	}
}
func TestCalc_DivisionByZero(t *testing.T) {
	result, err := Calc("10/0")
	assert.Error(t, err)
	assert.Equal(t, errors.New("division by zero"), err)
	assert.Equal(t, 0.0, result)
}

func TestCalc_UnbalancedParentheses(t *testing.T) {
	result, err := Calc("(2+2")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidExpression, err)
	assert.Equal(t, 0.0, result)
}

func TestParseExpression(t *testing.T) {
	tokens := []string{"2", "+", "3", "-", "1"}
	result, err := parseExpression(&tokens)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, result)
}

func TestParseTerm(t *testing.T) {
	tokens := []string{"2", "*", "3", "/", "6"}
	result, err := parseTerm(&tokens)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, result)
}

func TestParseTerm_DivideByZero(t *testing.T) {
	tokens := []string{"2", "/", "0"}
	result, err := parseTerm(&tokens)
	assert.Error(t, err)
	assert.EqualError(t, err, "division by zero")
	assert.Equal(t, 0.0, result)
}
