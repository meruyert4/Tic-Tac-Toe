package handlers

import (
	"tic-tac-toe/backend/logger"

	"golang.org/x/net/websocket"
)

func WebSocketHandler(ws *websocket.Conn) {
	logger.Info("New WebSocket Connection")
	defer ws.Close()

	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			logger.Error("WebSocket recieve an error", "error", err)
			break
		}

		logger.Info("Received message", "message", msg)
		if err := websocket.Message.Send(ws, "Echo: "+msg); err != nil {
			logger.Error("WebSocket send error", "error", err)
			break
		}
	}
}
