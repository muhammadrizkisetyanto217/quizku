
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type UserThemesOrLevelsModel struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID         `gorm:"type:uuid;not null;index:idx_user_themes_user_theme,unique" json:"user_id"`
	ThemesOrLevelsID uint               `gorm:"not null;index:idx_user_themes_user_theme,unique" json:"themes_or_levels_id"`
	CompleteUnit     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"complete_unit"`
	TotalUnit        pq.Int64Array     `gorm:"type:integer[];default:'{}'" json:"total_unit"`
	GradeResult      int               `gorm:"default:0" json:"grade_result"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserThemesOrLevelsModel) TableName() string {
	return "user_themes_or_levels"
}