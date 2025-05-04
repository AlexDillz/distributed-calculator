package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexDillz/distributed-calculator/internal/contextkeys"
	"github.com/AlexDillz/distributed-calculator/internal/storage"
)

func CalculateHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextkeys.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Expression string `json:"expression"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		result, calcErr := EvaluateExpression(req.Expression)
		if calcErr != nil {
			store.SaveExpression(userID, req.Expression, 0, calcErr.Error())
			http.Error(w, calcErr.Error(), http.StatusInternalServerError)
			return
		}

		store.SaveExpression(userID, req.Expression, result, "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]float64{"result": result})
	}
}

func ListExpressionsHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(contextkeys.UserIDKey).(int)

		exprs, err := store.ListExpressions(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var out struct {
			Expressions []map[string]interface{} `json:"expressions"`
		}
		for _, e := range exprs {
			status := "done"
			if e.Error != "" {
				status = "error"
			}
			out.Expressions = append(out.Expressions, map[string]interface{}{
				"id":     e.ID,
				"status": status,
				"result": e.Result,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

func GetExpressionHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(contextkeys.UserIDKey).(int)

		parts := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		e, err := store.GetExpression(userID, id)
		if err != nil {
			if err == storage.ErrNotFound {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		status := "done"
		if e.Error != "" {
			status = "error"
		}
		resp := map[string]map[string]interface{}{
			"expression": {
				"id":     e.ID,
				"status": status,
				"result": e.Result,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
