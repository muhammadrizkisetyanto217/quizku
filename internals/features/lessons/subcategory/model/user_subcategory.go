package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserSubcategoryModel struct {
	ID                     uint              `gorm:"primaryKey" json:"id"`
	UserID                 uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	SubcategoryID          int               `gorm:"not null" json:"subcategories_id"`
	CompleteThemesOrLevels datatypes.JSONMap `gorm:"type:jsonb" json:"complete_themes_or_levels"`
	TotalThemesOrLevels    pq.Int64Array     `gorm:"type:integer[];default:'{}'" json:"total_themes_or_levels"`
	GradeResult            int               `gorm:"default:0" json:"grade_result"`
	CreatedAt              time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

func (UserSubcategoryModel) TableName() string {
	return "user_subcategory"
}
