package models

import "golang.org/x/net/websocket"

type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
	InGame   bool   `json:"ingame"`
}

type UserSession struct {
	UserID    string          `json:"userId"`
	SessionID string          `json:"sessionId"`
	Conn      *websocket.Conn `json:"-"`
}
