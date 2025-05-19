package model

import (
	"log"
	"time"

	themesOrLevelsModel "quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SubcategoryModel struct {
	SubcategoryID                  uint           `json:"subcategory_id" gorm:"primaryKey;column:subcategory_id;autoIncrement"`
	SubcategoryName                string         `json:"subcategory_name" gorm:"type:varchar(255);column:subcategory_name"`
	SubcategoryStatus              string         `json:"subcategory_status" gorm:"type:varchar(10);default:'pending';check:subcategory_status IN ('active','pending','archived')"`
	SubcategoryDescriptionLong     string         `json:"subcategory_description_long" gorm:"type:text;column:subcategory_description_long"`
	SubcategoryTotalThemesOrLevels pq.Int64Array  `json:"subcategory_total_themes_or_levels" gorm:"type:integer[];default:'{}';column:subcategory_total_themes_or_levels"`
	SubcategoryImageURL            string         `json:"subcategory_image_url" gorm:"type:text;column:subcategory_image_url"`
	
	CreatedAt           time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt           *time.Time     `json:"updated_at" gorm:"column:subcategory_updated_at"`
	DeletedAt           gorm.DeletedAt `json:"subcategory_deleted_at" gorm:"index;column:subcategory_deleted_at"`
	SubcategoryCategoryID          uint           `json:"subcategory_category_id" gorm:"column:subcategory_category_id"`

	ThemesOrLevels []themesOrLevelsModel.ThemesOrLevelsModel `json:"themes_or_levels" gorm:"foreignKey:ThemesOrLevelSubcategoryID;references:SubcategoryID"`
}

func (SubcategoryModel) TableName() string {
	return "subcategories"
}

func (s *SubcategoryModel) AfterSave(tx *gorm.DB) (err error) {
	return SyncTotalSubcategories(tx, s.SubcategoryCategoryID)
}

func (s *SubcategoryModel) AfterDelete(tx *gorm.DB) (err error) {
	return SyncTotalSubcategories(tx, s.SubcategoryCategoryID)
}

func SyncTotalSubcategories(db *gorm.DB, categoryID uint) error {
	log.Println("[SERVICE] SyncTotalSubcategories - categoryID:", categoryID)

	err := db.Exec(`
		UPDATE categories
		SET category_total_subcategories = (
			SELECT ARRAY_AGG(subcategory_id)
			FROM subcategories
			WHERE subcategory_category_id = ? AND subcategory_deleted_at IS NULL
		)
		WHERE category_id = ?
	`, categoryID, categoryID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync category_total_subcategories:", err)
	}

	return err
}
