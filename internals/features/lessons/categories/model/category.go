package model

import (
	"log"
	"time"

	subcategoriesModel "quizku/internals/features/lessons/subcategory/model"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID                 uint                                  `json:"id" gorm:"primaryKey"`
	Name               string                                `json:"name" gorm:"size:255;not null"`
	Status             string                                `json:"status" gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')"`
	DescriptionShort   string                                `json:"description_short" gorm:"size:100"`
	DescriptionLong    string                                `json:"description_long" gorm:"size:2000"`
	TotalSubcategories pq.Int64Array                         `json:"total_subcategories" gorm:"type:integer[];default:'{}'"`
	ImageURL           string                                `json:"image_url" gorm:"size:100"`
	UpdateNews         datatypes.JSON                        `json:"update_news"`
	DifficultyID       uint                                  `json:"difficulty_id"`
	CreatedAt          time.Time                             `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time                             `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt                        `json:"deleted_at" gorm:"index"`
	Subcategories      []subcategoriesModel.SubcategoryModel `json:"subcategories" gorm:"foreignKey:CategoriesID"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

func (c *CategoryModel) AfterSave(tx *gorm.DB) (err error) {
	return SyncTotalCategories(tx, c.DifficultyID)
}

func (c *CategoryModel) AfterDelete(tx *gorm.DB) (err error) {
	return SyncTotalCategories(tx, c.DifficultyID)
}

func SyncTotalCategories(db *gorm.DB, difficultyID uint) error {
	log.Println("[SERVICE] SyncTotalCategories - difficultyID:", difficultyID)

	err := db.Exec(`
		UPDATE difficulties
		SET total_categories = (
			SELECT ARRAY_AGG(id ORDER BY id)
			FROM categories
			WHERE difficulty_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, difficultyID, difficultyID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_categories:", err)
	}
	return err
}