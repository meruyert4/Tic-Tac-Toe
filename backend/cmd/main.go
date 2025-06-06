package main

import (
	"net/http"
	"os"
	"tic-tac-toe/backend/handlers"
	"tic-tac-toe/backend/logger"
	"tic-tac-toe/backend/utils"

	"golang.org/x/net/websocket"
)

func main() {
	// logger
	logger, logFile := logger.SetupLogger()
	defer logFile.Close()

	// Redis
	utils.InitRedis()
	logger.Info("Redis connected successfully")

	// Setup HTTP routes
	http.HandleFunc("/api/health", handlers.HealthCheck)
	http.Handle("/ws", websocket.Handler(handlers.WebSocketHandler))

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server running http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
