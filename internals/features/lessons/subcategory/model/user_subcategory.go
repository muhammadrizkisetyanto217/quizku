package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type UserSubcategoryModel struct {
	ID                     uint                 `gorm:"primaryKey" json:"id"`
	UserID                 uuid.UUID            `gorm:"type:uuid;not null;index:idx_user_subcategory_user_subcat,unique" json:"user_id"`
	SubcategoryID          int                  `gorm:"not null;index:idx_user_subcategory_user_subcat,unique" json:"subcategory_id"`
	CompleteThemesOrLevels datatypes.JSONMap    `gorm:"type:jsonb;default:'{}'" json:"complete_themes_or_levels"`
	TotalThemesOrLevels    pq.Int64Array        `gorm:"type:integer[];default:'{}'" json:"total_themes_or_levels"`
	GradeResult            int                  `gorm:"default:0" json:"grade_result"`
	CreatedAt              time.Time            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserSubcategoryModel) TableName() string {
	return "user_subcategory"
}