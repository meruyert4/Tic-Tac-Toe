package game

import (
	"math/rand"
	"sync"
	"tic-tac-toe/backend/models"
	"tic-tac-toe/backend/session"
)

type QuickGame struct {
	Players    [2]*models.User
	Board      [3][3]string
	Current    int
	GameOver   bool
	Winner     string
	Symbols    map[string]string
	Stats      *models.QuickGameStats
	StatsMutex sync.Mutex
}

var (
	quickGamePool   = make(map[string]*QuickGame)
	searchingPlayer *models.User
	poolMutex       sync.Mutex
)

func NewQuickGame() *QuickGame {
	return &QuickGame{
		Board:    [3][3]string{},
		Symbols:  make(map[string]string),
		GameOver: false,
	}
}

func (g *QuickGame) GetOpponent(playerID string) *models.User {
	for _, player := range g.Players {
		if player != nil && player.ID != playerID {
			return player
		}
	}
	return nil
}

// CurrentPlayer returns the player whose turn it is
func (g *QuickGame) CurrentPlayer() *models.User {
	if g.Current < 0 || g.Current >= len(g.Players) {
		return nil
	}
	return g.Players[g.Current]
}

func (g *QuickGame) Start() {
	// Assign random symbols
	symbols := []string{"X", "O"}
	rand.Shuffle(len(symbols), func(i, j int) {
		symbols[i], symbols[j] = symbols[j], symbols[i]
	})

	g.Symbols[g.Players[0].ID] = symbols[0]
	g.Symbols[g.Players[1].ID] = symbols[1]

	// Randomize who goes first
	g.Current = rand.Intn(2)

	g.updateStats(true)
}

func (g *QuickGame) MakeMove(playerID string, row, col int) bool {
	if g.GameOver || g.Players[g.Current].ID != playerID || g.Board[row][col] != "" {
		return false
	}

	g.Board[row][col] = g.Symbols[playerID]

	if g.checkWin(row, col) {
		g.GameOver = true
		g.Winner = playerID
		g.updateStats(false)
		return true
	}

	if g.checkDraw() {
		g.GameOver = true
		g.updateStats(false)
		return true
	}

	g.Current = (g.Current + 1) % 2
	return true
}

func (g *QuickGame) checkWin(row, col int) bool {
	symbol := g.Board[row][col]

	// Check row
	if g.Board[row][0] == symbol && g.Board[row][1] == symbol && g.Board[row][2] == symbol {
		return true
	}

	// Check column
	if g.Board[0][col] == symbol && g.Board[1][col] == symbol && g.Board[2][col] == symbol {
		return true
	}

	// Check diagonals
	if row == col && g.Board[0][0] == symbol && g.Board[1][1] == symbol && g.Board[2][2] == symbol {
		return true
	}

	if row+col == 2 && g.Board[0][2] == symbol && g.Board[1][1] == symbol && g.Board[2][0] == symbol {
		return true
	}

	return false
}

func (g *QuickGame) checkDraw() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func (g *QuickGame) updateStats(isNewGame bool) {
	g.StatsMutex.Lock()
	defer g.StatsMutex.Unlock()

	if isNewGame {
		g.Stats.TotalGames++
		g.Stats.ActiveGames++
	} else {
		g.Stats.ActiveGames--
	}
}

func FindOpponent(sess *session.UserSession) *QuickGame {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if searchingPlayer != nil {
		game := NewQuickGame()
		game.Players[0] = searchingPlayer
		game.Players[1] = (*models.User)(sess.User)
		game.Stats = &models.QuickGameStats{}
		searchingPlayer = nil
		return game
	}

	searchingPlayer = (*models.User)(sess.User)
	return nil
}

func CancelSearch(player *models.User) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if searchingPlayer != nil && searchingPlayer.ID == player.ID {
		searchingPlayer = nil
	}
}
