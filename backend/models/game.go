package models

type GameStatePayload struct {
	Board        [3][3]string `json:"board"`
	Current      string       `json:"current"`
	PlayerSymbol string       `json:"playerSymbol"`
	GameOver     bool         `json:"gameOver"`
	Winner       string       `json:"winner"`
}
