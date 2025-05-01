package model

import (
	"time"
)

type LevelRequirement struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Level      int       `json:"level" gorm:"unique;not null"`
	Name       string    `json:"name"`
	MinPoints  int       `json:"min_points" gorm:"not null"`
	MaxPoints  *int      `json:"max_points"` // nullable
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (LevelRequirement) TableName() string {
	return "level_requirements"
}
