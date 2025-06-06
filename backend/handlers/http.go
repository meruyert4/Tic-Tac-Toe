package handlers

import (
	"encoding/json"
	"net/http"
	"tic-tac-toe/backend/logger"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health check requested")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
