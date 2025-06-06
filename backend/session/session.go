package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
	"tic-tac-toe/randomname/services"
	"time"

	"golang.org/x/net/websocket"
)

var (
	sessions   = make(map[string]*UserSession)
	sessionMux sync.Mutex
)

type UserSession struct {
	SessionID string
	User      *User
	Conn      *websocket.Conn
	CreatedAt time.Time
}

type User struct {
	ID       string
	Nickname string
	Online   bool
	InGame   bool
}

func GetOrCreateSession(w http.ResponseWriter, r *http.Request, ws *websocket.Conn) *UserSession {
	sessionMux.Lock()
	defer sessionMux.Unlock()

	// Try to get session from cookie
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		if sess, exists := sessions[cookie.Value]; exists {
			// Update connection if this is a WebSocket upgrade
			if ws != nil {
				sess.Conn = ws
			}
			sess.User.Online = true
			return sess
		}
	}

	// Create new session
	sessionID := generateSessionID()
	userID := "user-" + generateRandomString(8)
	nickname := services.RandomNickName()

	sess := &UserSession{
		SessionID: sessionID,
		User: &User{
			ID:       userID,
			Nickname: nickname,
			Online:   true,
			InGame:   false,
		},
		Conn:      ws,
		CreatedAt: time.Now(),
	}

	sessions[sessionID] = sess

	// Set cookie if this is an HTTP request (not WebSocket)
	if ws == nil && w != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400, // 1 day
			HttpOnly: true,
			Secure:   strings.HasPrefix(r.URL.Scheme, "https"),
			SameSite: http.SameSiteLaxMode,
		})
	}

	return sess
}

func GetSession(sessionID string) *UserSession {
	sessionMux.Lock()
	defer sessionMux.Unlock()
	return sessions[sessionID]
}

func RemoveSession(sessionID string) {
	sessionMux.Lock()
	defer sessionMux.Unlock()
	delete(sessions, sessionID)
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
