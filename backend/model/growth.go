package model

type GrowthData struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Hours     float64 `json:"hours"`
	Points    int     `json:"points"`
	Runner    string  `json:"runner"`
	MapName   string  `json:"map_name"`
	MapPoints int     `json:"map_points"`
	Maps      int     `json:"maps"`
	Timestamp string  `json:"timestamp"`
}
