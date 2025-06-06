package game

import (
	"sync"
	"tic-tac-toe/backend/models"
)

type OfflineGame struct {
	Board      [3][3]string
	Current    string // "X" or "O"
	GameOver   bool
	Winner     string
	Stats      *models.OfflineStats
	StatsMutex sync.Mutex
}

func NewOfflineGame() *OfflineGame {
	return &OfflineGame{
		Board:    [3][3]string{},
		Current:  "X", // X always starts in offline mode
		GameOver: false,
		Stats:    &models.OfflineStats{},
	}
}

func (g *OfflineGame) Start() {
	g.Board = [3][3]string{}
	g.Current = "X"
	g.GameOver = false
	g.Winner = ""
}

func (g *OfflineGame) MakeMove(row, col int) bool {
	if g.GameOver || g.Board[row][col] != "" {
		return false
	}

	g.Board[row][col] = g.Current

	if g.checkWin(row, col) {
		g.GameOver = true
		g.Winner = g.Current
		g.updateStats(g.Winner)
		return true
	}

	if g.checkDraw() {
		g.GameOver = true
		g.updateStats("")
		return true
	}

	if g.Current == "X" {
		g.Current = "O"
	} else {
		g.Current = "X"
	}

	return true
}

func (g *OfflineGame) checkWin(row, col int) bool {
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

func (g *OfflineGame) checkDraw() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func (g *OfflineGame) updateStats(winner string) {
	g.StatsMutex.Lock()
	defer g.StatsMutex.Unlock()

	g.Stats.TotalGames++

	switch winner {
	case "X":
		g.Stats.Wins++
		g.Stats.CurrentStreak++
	case "O":
		g.Stats.Losses++
		g.Stats.CurrentStreak = 0
	default: // draw
		g.Stats.Draws++
		g.Stats.CurrentStreak = 0
	}
}
