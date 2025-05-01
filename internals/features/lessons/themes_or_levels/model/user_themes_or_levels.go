package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserThemesOrLevelsModel struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	ThemesOrLevelsID uint              `gorm:"not null" json:"themes_or_levels_id"`
	CompleteUnit     datatypes.JSONMap `gorm:"type:jsonb" json:"complete_unit"`
	TotalUnit        pq.Int64Array     `gorm:"type:integer[];default:'{}'" json:"total_unit"`
	GradeResult      int               `gorm:"default:0" json:"grade_result"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

func (UserThemesOrLevelsModel) TableName() string {
	return "user_themes_or_levels"
}
