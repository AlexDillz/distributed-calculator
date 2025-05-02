package server

import (
	"github.com/AlexDillz/distributed-calculator/internal/storage"
)

func SaveUserExpression(store *storage.Storage, userID int, expr string, result float64, errStr string) {
	store.SaveExpression(userID, expr, result, errStr)
}
