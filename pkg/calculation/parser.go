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
		if isDigitOrDot(ch) || ch == 'e' || ch == 'E' {
			num.WriteRune(ch)
			continue
		}

		if num.Len() > 0 {
			tokens = append(tokens, num.String())
			num.Reset()
		}

		if ch == '-' && (i == 0 || isStartOfNewTerm(rune(expr[i-1]))) {
			tokens = append(tokens, "-1")
			tokens = append(tokens, "*")
			continue
		}

		if isOperator(ch) {
			tokens = append(tokens, string(ch))
		} else if ch == '(' || ch == ')' {
			tokens = append(tokens, string(ch))
		} else {
			return nil, ErrInvalidExpression
		}
	}

	if num.Len() > 0 {
		tokens = append(tokens, num.String())
	}
	return tokens, nil
}

func isStartOfNewTerm(ch rune) bool {
	return ch == '(' || ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func isDigitOrDot(ch rune) bool {
	return (ch >= '0' && ch <= '9') || ch == '.'
}

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
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

	if strings.HasPrefix(token, "-") {
		if len(token) == 1 {
			return 0, ErrInvalidExpression
		}
		value, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return 0, ErrInvalidExpression
		}
		return value, nil
	}

	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, ErrInvalidExpression
	}
	return value, nil
}
