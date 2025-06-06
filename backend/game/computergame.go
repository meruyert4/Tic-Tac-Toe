package game

import (
	"math/rand"
	"sync"
	"time"

	"tic-tac-toe/backend/models"
)

type ComputerGame struct {
	Board      [3][3]string
	Player     string // "X" or "O"
	Computer   string // opposite of Player
	Current    string // "X" or "O"
	GameOver   bool
	Winner     string
	Stats      *models.ComputerStats
	StatsMutex sync.Mutex
}

func NewComputerGame(playerSymbol string) *ComputerGame {
	computerSymbol := "O"
	if playerSymbol == "O" {
		computerSymbol = "X"
	}

	return &ComputerGame{
		Board:    [3][3]string{},
		Player:   playerSymbol,
		Computer: computerSymbol,
		Current:  "X", // X always starts
		GameOver: false,
		Stats:    &models.ComputerStats{},
	}
}

func (g *ComputerGame) Start() {
	g.Board = [3][3]string{}
	g.GameOver = false
	g.Winner = ""
	g.Current = "X"

	if g.Current == g.Computer {
		g.makeComputerMove()
	}
}

func (g *ComputerGame) MakePlayerMove(row, col int) bool {
	if g.GameOver || g.Current != g.Player || g.Board[row][col] != "" {
		return false
	}

	g.Board[row][col] = g.Player

	if g.checkWin(row, col) {
		g.GameOver = true
		g.Winner = g.Player
		g.updateStats(true)
		return true
	}

	if g.checkDraw() {
		g.GameOver = true
		g.updateStats(false)
		return true
	}

	g.Current = g.Computer
	g.makeComputerMove()
	return true
}

func (g *ComputerGame) makeComputerMove() {
	if g.GameOver {
		return
	}

	// Simple AI: first try to win, then block player, then random move
	row, col := g.findWinningMove(g.Computer)
	if row == -1 {
		row, col = g.findWinningMove(g.Player) // block player
	}
	if row == -1 {
		emptyCells := [][2]int{}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if g.Board[i][j] == "" {
					emptyCells = append(emptyCells, [2]int{i, j})
				}
			}
		}

		if len(emptyCells) > 0 {
			rand.Seed(time.Now().UnixNano())
			choice := emptyCells[rand.Intn(len(emptyCells))]
			row, col = choice[0], choice[1]
		} else {
			return
		}
	}

	g.Board[row][col] = g.Computer

	if g.checkWin(row, col) {
		g.GameOver = true
		g.Winner = g.Computer
		g.updateStats(false)
		return
	}

	if g.checkDraw() {
		g.GameOver = true
		g.updateStats(false)
		return
	}

	g.Current = g.Player
}

func (g *ComputerGame) findWinningMove(symbol string) (int, int) {
	// Check rows
	for i := 0; i < 3; i++ {
		count := 0
		emptyCol := -1
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == symbol {
				count++
			} else if g.Board[i][j] == "" {
				emptyCol = j
			}
		}
		if count == 2 && emptyCol != -1 {
			return i, emptyCol
		}
	}

	// Check columns
	for j := 0; j < 3; j++ {
		count := 0
		emptyRow := -1
		for i := 0; i < 3; i++ {
			if g.Board[i][j] == symbol {
				count++
			} else if g.Board[i][j] == "" {
				emptyRow = i
			}
		}
		if count == 2 && emptyRow != -1 {
			return emptyRow, j
		}
	}

	// Check diagonals
	count := 0
	emptyPos := -1
	for i := 0; i < 3; i++ {
		if g.Board[i][i] == symbol {
			count++
		} else if g.Board[i][i] == "" {
			emptyPos = i
		}
	}
	if count == 2 && emptyPos != -1 {
		return emptyPos, emptyPos
	}

	count = 0
	emptyPos = -1
	for i := 0; i < 3; i++ {
		if g.Board[i][2-i] == symbol {
			count++
		} else if g.Board[i][2-i] == "" {
			emptyPos = i
		}
	}
	if count == 2 && emptyPos != -1 {
		return emptyPos, 2 - emptyPos
	}

	return -1, -1
}

func (g *ComputerGame) checkWin(row, col int) bool {
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

func (g *ComputerGame) checkDraw() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func (g *ComputerGame) updateStats(playerWon bool) {
	g.StatsMutex.Lock()
	defer g.StatsMutex.Unlock()

	if playerWon {
		g.Stats.PlayerWin++
	} else {
		g.Stats.ComputerWin++
	}
}
