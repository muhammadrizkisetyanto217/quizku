package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSubcategoryModel struct {
	ID                     uint              `gorm:"primaryKey" json:"id"`
	UserID                 uuid.UUID         `gorm:"type:uuid;not null;index:idx_user_subcategory_user_subcat,unique" json:"user_id"`
	SubcategoryID          int               `gorm:"not null;index:idx_user_subcategory_user_subcat,unique" json:"subcategory_id"`
	CompleteThemesOrLevels datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"complete_themes_or_levels"`
	GradeResult            int               `gorm:"default:0" json:"grade_result"`
	CurrentVersion         int               `gorm:"default:1;index:idx_user_subcategory_version" json:"current_version"`
	CreatedAt              time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}


func (UserSubcategoryModel) TableName() string {
	return "user_subcategory"
}
