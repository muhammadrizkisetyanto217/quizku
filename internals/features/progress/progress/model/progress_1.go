package model

import (
	"time"

	"github.com/google/uuid"
)

type UserProgress struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	TotalPoints int       `json:"total_points"`
	Level       int       `json:"level" gorm:"default:1"`
	Rank        int       `json:"rank" gorm:"not null;default:1"`
	LastUpdated time.Time `json:"last_updated"`
}

func (UserProgress) TableName() string {
	return "user_progress"
}
