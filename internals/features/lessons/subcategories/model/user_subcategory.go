package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSubcategoryModel struct {
	UserSubcategoryID                     uint              `gorm:"column:user_subcategory_id;primaryKey;autoIncrement" json:"user_subcategory_id"`
	UserSubcategoryUserID                 uuid.UUID         `gorm:"column:user_subcategory_user_id;type:uuid;not null;index:idx_user_subcategory_user_subcategory,unique" json:"user_subcategory_user_id"`
	UserSubcategorySubcategoryID          int               `gorm:"column:user_subcategory_subcategory_id;not null;index:idx_user_subcategory_user_subcategory,unique" json:"user_subcategory_subcategory_id"`
	UserSubcategoryCompleteThemesOrLevels datatypes.JSONMap `gorm:"column:user_subcategory_complete_themes_or_levels;type:jsonb;default:'{}'" json:"user_subcategory_complete_themes_or_levels"`
	UserSubcategoryGradeResult            int               `gorm:"column:user_subcategory_grade_result;default:0" json:"user_subcategory_grade_result"`
	UserSubcategoryCurrentVersion		 int               `gorm:"column:user_subcategory_current_version;default:1" json:"user_subcategory_current_version"`
	CreatedAt                             time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt                             time.Time         `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (UserSubcategoryModel) TableName() string {
	return "user_subcategory"
}
