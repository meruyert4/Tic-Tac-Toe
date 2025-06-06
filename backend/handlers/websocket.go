package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"tic-tac-toe/backend/game"
	"tic-tac-toe/backend/logger"
	"tic-tac-toe/backend/models"
	"tic-tac-toe/backend/session"
	"time"

	"golang.org/x/net/websocket"
)

var (
	connections   = make(map[string]*websocket.Conn)
	connMutex     sync.Mutex
	activeGames   = make(map[string]*game.QuickGame)
	activeGameMux sync.Mutex
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type MovePayload struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type GameStatePayload struct {
	Board        [3][3]string `json:"board"`
	Current      string       `json:"current"`      // player ID whose turn it is
	PlayerSymbol string       `json:"playerSymbol"` // X or O for the receiving player
	GameOver     bool         `json:"gameOver"`
	Winner       string       `json:"winner"` // player ID if there's a winner
}

func WebSocketHandler(ws *websocket.Conn) {
	logger.Info("New WebSocket connection")
	defer ws.Close()

	req := ws.Request()

	if req.Response == nil {
		req.Response = &http.Response{Header: make(http.Header)}
	}

	sess := session.GetOrCreateSession(nil, req, ws)
	if sess == nil {
		logger.Error("Failed to create session")
		return
	}

	// Store connection
	connMutex.Lock()
	connections[sess.User.ID] = ws
	connMutex.Unlock()

	defer func() {
		// Cleanup when connection closes
		connMutex.Lock()
		delete(connections, sess.User.ID)
		connMutex.Unlock()

		// Mark user as offline
		sess.User.Online = false

		// Handle player disconnection from any active game
		handlePlayerDisconnect(sess.User.ID)
	}()

	logger.Info("User connected", "nickname", sess.User.Nickname, "id", sess.User.ID)

	// Send welcome message with session info
	sendMessage(ws, WSMessage{
		Type: "welcome",
		Payload: map[string]interface{}{
			"userId":    sess.User.ID,
			"nickname":  sess.User.Nickname,
			"sessionId": sess.SessionID,
			"inGame":    sess.User.InGame,
		},
	})

	if sess.User.InGame {
		activeGameMux.Lock()
		for gameID, game := range activeGames {
			for _, player := range game.Players {
				if player.ID == sess.User.ID {
					sendMessage(ws, WSMessage{
						Type: "gameReconnect",
						Payload: map[string]interface{}{
							"gameId":     gameID,
							"opponent":   game.GetOpponent(sess.User.ID).Nickname,
							"yourSymbol": game.Symbols[sess.User.ID],
							"yourTurn":   game.CurrentPlayer().ID == sess.User.ID,
							"board":      game.Board,
						},
					})
					break
				}
			}
		}
		activeGameMux.Unlock()
	}

	for {
		var msgBytes []byte
		if err := websocket.Message.Receive(ws, &msgBytes); err != nil {
			logger.Error("WebSocket receive error", "error", err)
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			logger.Error("Failed to unmarshal message", "error", err)
			continue
		}

		logger.Info("Received message", "type", msg.Type, "payload", msg.Payload)

		switch msg.Type {
		case "findGame":
			handleFindGame(sess)
		case "makeMove":
			handleMakeMove(sess, msg.Payload)
		case "cancelSearch":
			handleCancelSearch(sess)
		case "playAgain":
			handlePlayAgain(sess)
		case "heartbeat":
			// Respond to heartbeat to keep connection alive
			sendMessage(ws, WSMessage{Type: "heartbeat", Payload: time.Now().Unix()})
		default:
			logger.Warn("Unknown message type", "type", msg.Type)
		}
	}
}

func handleFindGame(sess *session.UserSession) {
	conn := connections[sess.User.ID]
	if conn == nil {
		return
	}

	// Add player to matchmaking
	game := game.FindOpponent(sess)
	if game == nil {
		// No opponent found yet, waiting
		sendMessage(conn, WSMessage{
			Type:    "searching",
			Payload: nil,
		})
		return
	}

	// Game found, notify both players
	gameID := "game-" + strconv.Itoa(rand.Intn(1000000000))

	activeGameMux.Lock()
	activeGames[gameID] = game
	activeGameMux.Unlock()

	// Start the game
	game.Start()

	// Notify both players
	for i, player := range game.Players {
		if playerConn, ok := connections[player.ID]; ok {
			sendMessage(playerConn, WSMessage{
				Type: "gameStart",
				Payload: map[string]interface{}{
					"gameId":     gameID,
					"opponent":   game.Players[1-i].Nickname,
					"yourSymbol": game.Symbols[player.ID],
					"yourTurn":   game.Current == i,
				},
			})
		}
	}
}

func handleMakeMove(sess *session.UserSession, payload interface{}) {
	activeGameMux.Lock()
	defer activeGameMux.Unlock()

	// Find the game the player is in
	var game *game.QuickGame
	var gameID string
	for id, g := range activeGames {
		for _, p := range g.Players {
			if p.ID == sess.User.ID {
				game = g
				gameID = id
				break
			}
		}
		if game != nil {
			break
		}
	}

	if game == nil {
		return
	}

	// Parse move
	var move MovePayload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal move payload", "error", err)
		return
	}
	if err := json.Unmarshal(jsonPayload, &move); err != nil {
		logger.Error("Failed to unmarshal move payload", "error", err)
		return
	}

	// Make the move
	validMove := game.MakeMove(sess.User.ID, move.Row, move.Col)
	if !validMove {
		return
	}

	// Send game state to both players
	for _, player := range game.Players {
		if playerConn, ok := connections[player.ID]; ok {
			sendMessage(playerConn, WSMessage{
				Type: "gameState",
				Payload: GameStatePayload{
					Board:        game.Board,
					Current:      game.Players[game.Current].ID,
					PlayerSymbol: game.Symbols[player.ID],
					GameOver:     game.GameOver,
					Winner:       game.Winner,
				},
			})
		}
	}

	// If game over, send result and clean up
	if game.GameOver {
		for _, player := range game.Players {
			if playerConn, ok := connections[player.ID]; ok {
				result := "draw"
				if game.Winner != "" {
					if game.Winner == player.ID {
						result = "win"
					} else {
						result = "loss"
					}
				}

				sendMessage(playerConn, WSMessage{
					Type: "gameOver",
					Payload: map[string]interface{}{
						"result": result,
						"gameId": gameID,
					},
				})
			}
		}

		// Remove game from active games
		delete(activeGames, gameID)
	}
}

func handleCancelSearch(sess *session.UserSession) {
	game.CancelSearch((*models.User)(sess.User))
	if conn, ok := connections[sess.User.ID]; ok {
		sendMessage(conn, WSMessage{
			Type:    "searchCancelled",
			Payload: nil,
		})
	}
}

func handlePlayAgain(sess *session.UserSession) {
	handleFindGame(sess)
}

func handlePlayerDisconnect(userID string) {
	activeGameMux.Lock()
	defer activeGameMux.Unlock()

	// Find any game the player was in
	var game *game.QuickGame
	var gameID string
	for id, g := range activeGames {
		for _, p := range g.Players {
			if p.ID == userID {
				game = g
				gameID = id
				break
			}
		}
		if game != nil {
			break
		}
	}

	if game == nil {
		return
	}

	// Notify the other player
	for _, player := range game.Players {
		if player.ID != userID {
			if conn, ok := connections[player.ID]; ok {
				sendMessage(conn, WSMessage{
					Type:    "opponentDisconnected",
					Payload: nil,
				})
			}
		}
	}

	// Remove the game
	delete(activeGames, gameID)
}

func sendMessage(ws *websocket.Conn, msg WSMessage) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal message", "error", err)
		return
	}

	if err := websocket.Message.Send(ws, string(msgBytes)); err != nil {
		logger.Error("Failed to send WebSocket message", "error", err)
	}
}
