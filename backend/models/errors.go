package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

const (
	ErrUserNotFound = "user not found"
	ErrGameNotFound = "game not found"
	ErrInvalidMove  = "invalid move"
	ErrNotYourTurn  = "not your turn"
	ErrGameFinished = "game already finished"
)
