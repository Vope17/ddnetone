package model

import (
	"gorm.io/gorm"
	"time"
)

type MapRecord struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Difficulty string     `json:"difficulty"`
	MapName    string     `json:"map_name"`
	Runner     string     `json:"runner"`
	Score      int        `json:"score"`
	Points     int        `json:"points"`
	Stars      int        `json:"stars"`
	Note       string     `json:"note"`
	Status     int        `json:"status"` // 0:未完成, 1:進行中, 2:已完成
	FinishTime *time.Time `gorm:"column:finish_time" json:"finish_time"`

	HasDummy bool `gorm:"column:has_dummy" json:"has_dummy"`
}

// BeforeSave Hook
func (m *MapRecord) BeforeSave(tx *gorm.DB) error {
	if m.Runner == "" || m.Runner == "-" || m.Runner == "nan" {
		m.Status = 0
	} else if m.Status == 1 {
		// 明確指定為 WIP (1) 時保持原樣
	} else if m.Score > 0 {

		m.Status = 2
	} else {

		m.Status = 1
	}

	if m.Status == 2 && m.Score == 0 && m.Points > 0 {
		m.Score = m.Points
	}

	if m.Status == 2 && m.FinishTime == nil {
		now := time.Now()
		m.FinishTime = &now
	}

	return nil
}
