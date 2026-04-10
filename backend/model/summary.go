package model

import "time"

type Summary struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CurrentScore  int       `json:"current_score"`
	TargetScore   int       `json:"target_score"`
	CompletedMaps int       `json:"completed_maps"`
	LoadedMaps    int       `json:"loaded_maps"`
	LastUpdate    time.Time `json:"last_update"`
}
