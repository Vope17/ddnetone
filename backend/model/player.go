package model

type Player struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// PlayerStats 用於 leaderboard 回傳（從 map_records 即時計算）
type PlayerStats struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Role              string  `json:"role"`
	ScoreContribution float64 `json:"score_contrib"`
	MapCount          int     `json:"map_count"`
}
