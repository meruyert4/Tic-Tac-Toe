package models

import (
	"time"

	"golang.org/x/net/websocket"
)

type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
	InGame   bool   `json:"ingame"`
}

type UserSession struct {
	SessionID string
	User      *User
	Conn      *websocket.Conn
	CreatedAt time.Time
}
