package model

type GrowthData struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Hours  float64 `json:"hours"`
	Points int     `json:"points"`
	Runner string  `json:"runner"`

	Maps      int    `json:"maps"`
	Timestamp string `json:"timestamp"`
}
