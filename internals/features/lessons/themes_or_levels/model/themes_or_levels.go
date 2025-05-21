package model

import (
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ThemesOrLevelsModel struct {
	ThemesOrLevelID               uint          `gorm:"column:themes_or_level_id;primaryKey;autoIncrement" json:"themes_or_level_id"`
	ThemesOrLevelName             string        `gorm:"column:themes_or_level_name;type:varchar(255);not null" json:"themes_or_level_name"`
	ThemesOrLevelStatus           string        `gorm:"column:themes_or_level_status;type:varchar(10);default:'pending';check:themes_or_level_status IN ('active','pending','archived')" json:"themes_or_level_status"`
	ThemesOrLevelDescriptionShort string        `gorm:"column:themes_or_level_description_short;type:varchar(100)" json:"themes_or_level_description_short"`
	ThemesOrLevelDescriptionLong  string        `gorm:"column:themes_or_level_description_long;type:varchar(2000)" json:"themes_or_level_description_long"`
	ThemesOrLevelTotalUnit        pq.Int64Array `gorm:"column:themes_or_level_total_unit;type:integer[];default:'{}'" json:"themes_or_level_total_unit"`
	ThemesOrLevelImageURL         string        `gorm:"column:themes_or_level_image_url;type:varchar(100)" json:"themes_or_level_image_url"`
	ThemesOrLevelSubcategoryID    int           `gorm:"column:themes_or_level_subcategory_id" json:"themes_or_level_subcategory_id"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

func (ThemesOrLevelsModel) TableName() string {
	return "themes_or_levels"
}

func (t *ThemesOrLevelsModel) AfterSave(tx *gorm.DB) error {
	log.Printf("[HOOK] AfterSave triggered for ThemeID: %d", t.ThemesOrLevelID)
	return SyncTotalThemesOrLevels(tx, t.ThemesOrLevelSubcategoryID)
}

func (t *ThemesOrLevelsModel) AfterDelete(tx *gorm.DB) error {
	log.Printf("[HOOK] AfterDelete triggered for ThemeID: %d", t.ThemesOrLevelID)

	var subcategoryID int
	if err := tx.Unscoped().
		Model(&ThemesOrLevelsModel{}).
		Select("themes_or_level_subcategory_id").
		Where("themes_or_level_id = ?", t.ThemesOrLevelID).
		Take(&subcategoryID).Error; err != nil {
		log.Println("[ERROR] Gagal ambil subcategory_id setelah delete:", err)
		return err
	}

	log.Printf("[HOOK] Ditemukan subcategoryID: %d untuk ThemeID: %d", subcategoryID, t.ThemesOrLevelID)
	return SyncTotalThemesOrLevels(tx, subcategoryID)
}

// Sync ulang ke subcategories.total_themes_or_levels
func SyncTotalThemesOrLevels(db *gorm.DB, subcategoryID int) error {
	log.Println("[SERVICE] SyncTotalThemesOrLevels - subcategoryID:", subcategoryID)

	err := db.Exec(`
		UPDATE subcategories
		SET subcategory_total_themes_or_levels = (
			SELECT ARRAY_AGG(themes_or_level_id ORDER BY themes_or_level_id)
			FROM themes_or_levels
			WHERE themes_or_level_subcategory_id = ? AND deleted_at IS NULL
		)
		WHERE subcategory_id = ?
	`, subcategoryID, subcategoryID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync subcategory_total_themes_or_levels:", err)
	} else {
		log.Println("[SUCCESS] Sync berhasil untuk subcategoryID:", subcategoryID)
	}

	return err
}
