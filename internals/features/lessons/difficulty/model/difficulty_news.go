package model

import (
	"time"

	"gorm.io/gorm"
)

type DifficultyNewsModel struct {
	DifficultyNewsID           uint           `gorm:"column:difficulty_news_id;primaryKey;autoIncrement" json:"difficulty_news_id"`
	DifficultyNewsTitle        string         `gorm:"column:difficulty_news_title;type:varchar(255);not null" json:"difficulty_news_title"`
	DifficultyNewsDescription  string         `gorm:"column:difficulty_news_description;type:text;not null" json:"difficulty_news_description"`
	DifficultyNewsIsPublic     bool           `gorm:"column:difficulty_news_is_public;default:true" json:"difficulty_news_is_public"`
	DifficultyNewsDifficultyID uint           `gorm:"column:difficulty_news_difficulty_id;not null" json:"difficulty_news_difficulty_id"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (DifficultyNewsModel) TableName() string {
	return "difficulties_news"
}
