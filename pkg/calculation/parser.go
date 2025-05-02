package calculation

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero    = errors.New("division by zero")
)

func logError(err error) {
	log.Printf("[ERROR] %v", err)
}

func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		logError(err)
		return 0, err
	}

	result, err := parseExpression(&tokens)
	if err != nil {
		logError(err)
		return 0, err
	}

	if len(tokens) != 0 {
		err := ErrInvalidExpression
		logError(err)
		return 0, err
	}

	return result, nil
}

func Tokenize(expr string) ([]string, error) {
	return tokenize(expr)
}

func tokenize(expr string) ([]string, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	var tokens []string
	var num strings.Builder

	for i, ch := range expr {
		if (ch >= '0' && ch <= '9') || ch == '.' || ch == 'e' || ch == 'E' {
			num.WriteRune(ch)
		} else {
			if num.Len() > 0 {
				tokens = append(tokens, num.String())
				num.Reset()
			}
			if ch == '-' && (i == 0 || expr[i-1] == '(') {
				tokens = append(tokens, "-1")
				tokens = append(tokens, "*")
			} else if strings.ContainsRune("+-*/()", ch) {
				tokens = append(tokens, string(ch))
			} else {
				return nil, ErrInvalidExpression
			}
		}
	}

	if num.Len() > 0 {
		tokens = append(tokens, num.String())
	}
	return tokens, nil
}

func parseExpression(tokens *[]string) (float64, error) {
	result, err := parseTerm(tokens)
	if err != nil {
		return 0, err
	}
	for len(*tokens) > 0 {
		op := (*tokens)[0]
		if op != "+" && op != "-" {
			break
		}
		*tokens = (*tokens)[1:]
		next, err := parseTerm(tokens)
		if err != nil {
			return 0, err
		}
		if op == "+" {
			result += next
		} else {
			result -= next
		}
	}
	return result, nil
}

func parseTerm(tokens *[]string) (float64, error) {
	result, err := parseFactor(tokens)
	if err != nil {
		return 0, err
	}
	for len(*tokens) > 0 {
		op := (*tokens)[0]
		if op != "*" && op != "/" {
			break
		}
		*tokens = (*tokens)[1:]
		next, err := parseFactor(tokens)
		if err != nil {
			return 0, err
		}
		if op == "*" {
			result *= next
		} else {
			if next == 0 {
				return 0, ErrDivisionByZero
			}
			result /= next
		}
	}
	return result, nil
}

func parseFactor(tokens *[]string) (float64, error) {
	if len(*tokens) == 0 {
		return 0, ErrInvalidExpression
	}
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	if token == "(" {
		result, err := parseExpression(tokens)
		if err != nil {
			return 0, err
		}
		if len(*tokens) == 0 || (*tokens)[0] != ")" {
			return 0, ErrInvalidExpression
		}
		*tokens = (*tokens)[1:]
		return result, nil
	}

	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, ErrInvalidExpression
	}
	return value, nil
}
