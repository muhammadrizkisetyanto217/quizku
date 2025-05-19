package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserThemesOrLevelsModel struct {
	UserThemeID              uint              `gorm:"column:user_theme_id;primaryKey;autoIncrement" json:"user_theme_id"`
	UserThemeUserID          uuid.UUID         `gorm:"column:user_theme_user_id;type:uuid;not null;index:idx_user_theme_unique,unique" json:"user_theme_user_id"`
	UserThemeThemesOrLevelID uint              `gorm:"column:user_theme_themes_or_level_id;not null;index:idx_user_theme_unique,unique" json:"user_theme_themes_or_level_id"`
	UserThemeCompleteUnit    datatypes.JSONMap `gorm:"column:user_theme_complete_unit;type:jsonb;default:'{}'" json:"user_theme_complete_unit"`
	UserThemeGradeResult     int               `gorm:"column:user_theme_grade_result;default:0" json:"user_theme_grade_result"`
	CreatedAt                time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (UserThemesOrLevelsModel) TableName() string {
	return "user_themes_or_levels"
}
