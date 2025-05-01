package model

import (
	"time"

	"github.com/google/uuid"
)

type UserPointLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Points     int       `json:"points"`
	SourceType int       `json:"source_type"`
	SourceID   int       `json:"source_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName untuk override nama default jika perlu
func (UserPointLog) TableName() string {
	return "user_point_logs"
}
