package server

import (
	"encoding/json"
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

		result, calcErr := EvaluateExpression(req.Expression)
		if calcErr != nil {
			// Сохраняем ошибку в БД
			store.SaveExpression(userID, req.Expression, 0, calcErr.Error())
			http.Error(w, calcErr.Error(), http.StatusInternalServerError)
			return
		}

		// Сохраняем успешное вычисление в БД
		store.SaveExpression(userID, req.Expression, result, "")
		json.NewEncoder(w).Encode(map[string]float64{"result": result})
	}
}
