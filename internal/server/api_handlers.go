package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexDillz/distributed-calculator/internal/storage"
)

func CalculateHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		var req struct {
			Expression string `json:"expression"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		result, err := EvaluateExpression(req.Expression)
		if err != nil {
			store.SaveExpression(userID, req.Expression, 0, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
			return
		}

		store.SaveExpression(userID, req.Expression, result, "")
		json.NewEncoder(w).Encode(map[string]float64{"result": result})
	}
}
