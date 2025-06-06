package models

import "golang.org/x/net/websocket"

type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
	InGame   bool   `json:"ingame"`
}

type UserSession struct {
	UserId    string          `json:"userid"`
	SessionId string          `json:"sessionid"`
	Conn      *websocket.Conn `json:"-"`
}
