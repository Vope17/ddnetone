package model

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User      string    `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
