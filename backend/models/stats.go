package models

type QuickGameStats struct {
	TotalGames  int `json:"totalgames"`
	ActiveGames int `json:"activegames"`
	OnlineUsers int `json:"onlineusers"`
}

type PlayerStats struct {
	Wins       int `json:"wins"`
	Losses     int `json:"losses"`
	Draws      int `json:"draws"`
	TotalGames int `json:"totalgames"`
}

type OfflineStats struct {
	PlayerStats
	CurrentStreak int `json:"currentstreak"`
}

type ComputerStats struct {
	PlayerWin   int `json:"playerwin"`
	ComputerWin int `json:"computerwin"`
}
