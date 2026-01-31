package model

type Player struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	Role              string  `json:"role"`
	ScoreContribution float64 `json:"score_contrib"`
	MapCount          int     `json:"map_count"`
	ContributionRate  float64 `json:"contribution_rate"`
}
