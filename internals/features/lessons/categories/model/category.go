package model

import (
	"log"
	"time"

	subcategoriesModel "quizku/internals/features/lessons/subcategories/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CategoryModel struct {
	CategoryID                 uint           `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName               string         `gorm:"size:255;not null;column:category_name" json:"category_name"`
	CategoryStatus             string         `gorm:"type:varchar(10);default:'pending';check:category_status IN ('active', 'pending', 'archived');column:category_status" json:"category_status"`
	CategoryDescriptionShort   string         `gorm:"size:100;column:category_description_short" json:"category_description_short"`
	CategoryDescriptionLong    string         `gorm:"size:2000;column:category_description_long" json:"category_description_long"`
	CategoryTotalSubcategories pq.Int64Array  `gorm:"type:integer[];default:'{}';column:category_total_subcategories" json:"category_total_subcategories"`
	CategoryImageURL           string         `gorm:"size:100;column:category_image_url" json:"category_image_url"`
	CategoryDifficultyID       uint           `gorm:"column:category_difficulty_id" json:"category_difficulty_id"`
	CreatedAt                  time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	Subcategories []subcategoriesModel.SubcategoryModel `gorm:"foreignKey:SubcategoryCategoryID;references:CategoryID" json:"subcategories"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

func (c *CategoryModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalCategories(tx, c.CategoryDifficultyID)
}

func (c *CategoryModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalCategories(tx, c.CategoryDifficultyID)
}

func SyncTotalCategories(db *gorm.DB, difficultyID uint) error {
	log.Println("[SERVICE] SyncTotalCategories - difficultyID:", difficultyID)

	err := db.Exec(`
		UPDATE difficulties
		SET difficulty_total_categories = (
			SELECT ARRAY_AGG(category_id ORDER BY category_id)
			FROM categories
			WHERE category_difficulty_id = ? AND deleted_at IS NULL
		)
		WHERE difficulty_id = ?
	`, difficultyID, difficultyID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync difficulty_total_categories:", err)
	}
	return err
}
