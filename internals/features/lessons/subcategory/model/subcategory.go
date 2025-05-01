package model

import (
	"log"
	"time"
	themesOrLevelsModel "quizku/internals/features/lessons/themes_or_levels/model"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubcategoryModel struct {
	ID                  uint                                      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string                                    `json:"name" gorm:"type:varchar(255)"`
	Status              string                                    `json:"status" gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')"`
	DescriptionLong     string                                    `json:"description_long" gorm:"type:varchar(2000)"`
	TotalThemesOrLevels pq.Int64Array                             `gorm:"type:integer[];default:'{}'" json:"total_themes_or_levels"`
	ImageURL            string                                    `json:"image_url" gorm:"type:varchar(100)"`
	UpdateNews          datatypes.JSON                            `json:"update_news"` // pakai JSONB di PostgreSQL
	CreatedAt           time.Time                                 `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           *time.Time                                `json:"updated_at"`
	DeletedAt           gorm.DeletedAt                            `json:"deleted_at" gorm:"index"`
	CategoriesID        uint                                      `json:"categories_id"`
	ThemesOrLevels      []themesOrLevelsModel.ThemesOrLevelsModel `json:"themes_or_levels" gorm:"foreignKey:SubcategoriesID"`
}

func (SubcategoryModel) TableName() string {
	return "subcategories"
}

func (s *SubcategoryModel) AfterSave(tx *gorm.DB) (err error) {
	return SyncTotalSubcategories(tx, s.CategoriesID)
}

func (s *SubcategoryModel) AfterDelete(tx *gorm.DB) (err error) {
	return SyncTotalSubcategories(tx, s.CategoriesID)
}

func SyncTotalSubcategories(db *gorm.DB, categoryID uint) error {
	log.Println("[SERVICE] SyncTotalSubcategories - categoryID:", categoryID)

	err := db.Exec(`
		UPDATE categories
		SET total_subcategories = (
			SELECT ARRAY_AGG(id)
			FROM subcategories
			WHERE categories_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, categoryID, categoryID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_subcategories:", err)
	}

	return err
}
