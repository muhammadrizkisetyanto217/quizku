package model

import (
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ThemesOrLevelsModel struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"type:varchar(255)" json:"name"`
	Status           string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')" json:"status"`
	DescriptionShort string         `gorm:"type:varchar(100)" json:"description_short"`
	DescriptionLong  string         `gorm:"type:varchar(2000)" json:"description_long"`
	TotalUnit        pq.Int64Array  `gorm:"type:integer[];default:'{}'" json:"total_unit"` 
	ImageURL         string         `gorm:"type:varchar(100)" json:"image_url"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	SubcategoriesID  int            `gorm:"column:subcategories_id" json:"subcategories_id"`
}

func (ThemesOrLevelsModel) TableName() string {
	return "themes_or_levels"
}


func (t *ThemesOrLevelsModel) AfterSave(tx *gorm.DB) (err error) {
	return SyncTotalThemesOrLevels(tx, t.SubcategoriesID)
}

func (t *ThemesOrLevelsModel) AfterDelete(tx *gorm.DB) (err error) {
	return SyncTotalThemesOrLevels(tx, t.SubcategoriesID)
}

func SyncTotalThemesOrLevels(db *gorm.DB, subcategoryID int) error {
	log.Println("[SERVICE] SyncTotalThemesOrLevels - subcategoryID:", subcategoryID)

	err := db.Exec(`
		UPDATE subcategories
		SET total_themes_or_levels = (
			SELECT ARRAY_AGG(id)
			FROM themes_or_levels
			WHERE subcategories_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, subcategoryID, subcategoryID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_themes_or_levels:", err)
	}
	return err
}

