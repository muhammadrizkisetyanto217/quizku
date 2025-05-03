package model

import (
	"time"

	"gorm.io/gorm"
)

type ThemesOrLevelsNewsModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"type:varchar(255);not null"`
	Description     string         `gorm:"type:text;not null"`
	IsPublic        bool           `gorm:"default:true"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	ThemesOrLevelsID uint           `gorm:"column:themes_or_levels_id;not null" json:"themes_or_levels_id"`
}

func (ThemesOrLevelsNewsModel) TableName() string {
	return "themes_or_levels_news"
}