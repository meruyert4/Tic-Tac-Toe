package models

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	WSGameStart    = "game_start"
	WSMoveMade     = "move_made"
	WSGameOver     = "game_over"
	WSOpponentLeft = "opponent_left"
	WSError        = "error"
)

type WSGameStartPayload struct {
	GameID    string `json:"gameId"`
	PlayerID  string `json:"playerId"`
	Symbol    string `json:"symbol"` // "X" or "O"
	FirstTurn string `json:"firstTurn"`
}

type WSMoveMadePayload struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Symbol string `json:"symbol"`
}

type WSGameOverPayload struct {
	Winner   string       `json:"winner"` // "X", "O", "draw"
	Board3x3 [3][3]string `json:"board3x3"`
	Board4x4 [4][4]string `json:"board4x4"`
}
