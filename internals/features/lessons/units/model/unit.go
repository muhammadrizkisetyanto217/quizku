package model

import (
	"log"
	"time"

	"quizku/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UnitModel struct {
	UnitID                  uint          `gorm:"column:unit_id;primaryKey;autoIncrement" json:"unit_id"`
	UnitName                string        `gorm:"column:unit_name;type:varchar(50);unique;not null" json:"unit_name"`
	UnitStatus              string        `gorm:"column:unit_status;type:varchar(10);default:'pending';check:unit_status IN ('active','pending','archived')" json:"unit_status"`
	UnitDescriptionShort    string        `gorm:"column:unit_description_short;type:varchar(200);not null" json:"unit_description_short"`
	UnitDescriptionOverview string        `gorm:"column:unit_description_overview;type:text;not null" json:"unit_description_overview"`
	UnitImageURL            string        `gorm:"column:unit_image_url;type:varchar(100)" json:"unit_image_url"`
	UnitTotalSectionQuizzes pq.Int64Array `gorm:"column:unit_total_section_quizzes;type:integer[];default:'{}'" json:"unit_total_section_quizzes"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	UnitThemesOrLevelID uint      `gorm:"column:unit_themes_or_level_id;not null" json:"unit_themes_or_level_id"`
	UnitCreatedBy       uuid.UUID `gorm:"column:unit_created_by;type:uuid;not null;constraint:OnDelete:CASCADE" json:"unit_created_by"`

	SectionQuizzes []model.SectionQuizzesModel `gorm:"foreignKey:SectionQuizzesUnitID;references:UnitID" json:"section_quizzes"`
}

func (UnitModel) TableName() string {
	return "units"
}

// ✅ Sync ketika ada perubahan
func (u *UnitModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalUnits(tx, u.UnitThemesOrLevelID)
}

func (u *UnitModel) AfterDelete(tx *gorm.DB) error {
	log.Printf("[HOOK] AfterDelete triggered for UnitID: %d", u.UnitID)

	var themesOrLevelID uint
	if err := tx.Unscoped().
		Model(&UnitModel{}).
		Select("unit_themes_or_level_id").
		Where("unit_id = ?", u.UnitID).
		Take(&themesOrLevelID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch unit_themes_or_level_id after delete:", err)
		return err
	}

	log.Printf("[HOOK] Fetched unit_themes_or_level_id: %d for deleted UnitID: %d", themesOrLevelID, u.UnitID)
	return SyncTotalUnits(tx, themesOrLevelID)
}

// ✅ Update ARRAY total_unit pada themes_or_levels
func SyncTotalUnits(db *gorm.DB, themesOrLevelID uint) error {
	log.Println("[SERVICE] SyncTotalUnits - themesOrLevelID:", themesOrLevelID)

	err := db.Exec(`
		UPDATE themes_or_levels
		SET themes_or_level_total_unit = (
			SELECT ARRAY_AGG(unit_id ORDER BY unit_id)
			FROM units
			WHERE unit_themes_or_level_id = ? AND deleted_at IS NULL
		)
		WHERE themes_or_level_id = ?
	`, themesOrLevelID, themesOrLevelID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_unit:", err)
	}

	return err
}
