package model

import (
	"time"

	"gorm.io/gorm"
)

type ThemesOrLevelsNewsModel struct {
	ThemesNewsID              uint   `gorm:"column:themes_news_id;primaryKey;autoIncrement" json:"themes_news_id"`
	ThemesNewsTitle           string `gorm:"column:themes_news_title;type:varchar(255);not null" json:"themes_news_title"`
	ThemesNewsDescription     string `gorm:"column:themes_news_description;type:text;not null" json:"themes_news_description"`
	ThemesNewsIsPublic        bool   `gorm:"column:themes_news_is_public;default:true" json:"themes_news_is_public"`
	ThemesNewsThemesOrLevelID uint   `gorm:"column:themes_news_themes_or_level_id;not null" json:"themes_news_themes_or_level_id"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (ThemesOrLevelsNewsModel) TableName() string {
	return "themes_or_levels_news"
}
