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

// AfterSave akan sinkronisasi total_themes_or_levels setelah create/update
func (t *ThemesOrLevelsModel) AfterSave(tx *gorm.DB) (err error) {
	log.Printf("[HOOK] AfterSave triggered for ThemeID: %d", t.ID)
	return SyncTotalThemesOrLevels(tx, t.SubcategoriesID)
}

// AfterDelete akan sinkronisasi ulang subkategori meskipun data sudah soft deleted
func (t *ThemesOrLevelsModel) AfterDelete(tx *gorm.DB) (err error) {
	log.Printf("[HOOK] AfterDelete triggered for ThemeID: %d", t.ID)

	var subcategoryID int
	if err := tx.Unscoped().
		Model(&ThemesOrLevelsModel{}).
		Select("subcategories_id").
		Where("id = ?", t.ID).
		Take(&subcategoryID).Error; err != nil {
		log.Println("[ERROR] Gagal ambil subcategories_id setelah delete:", err)
		return err
	}

	log.Printf("[HOOK] Ditemukan subcategoryID: %d untuk ThemeID: %d", subcategoryID, t.ID)
	return SyncTotalThemesOrLevels(tx, subcategoryID)
}

// Sinkronisasi ulang total_themes_or_levels (array of theme IDs) ke tabel subcategories
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
	} else {
		log.Println("[SUCCESS] Synced total_themes_or_levels for subcategoryID:", subcategoryID)
	}

	return err
}
