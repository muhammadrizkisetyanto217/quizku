package model

import (
	"time"
)

type RankRequirement struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Rank      int        `gorm:"unique;not null" json:"rank"`
	Name      string     `gorm:"type:varchar(100)" json:"name"`
	MinLevel  int        `gorm:"not null" json:"min_level"`
	MaxLevel  *int       `json:"max_level"` // nullable
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (RankRequirement) TableName() string {
	return "rank_requirements"
}
