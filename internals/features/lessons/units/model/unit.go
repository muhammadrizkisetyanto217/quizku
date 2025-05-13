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
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                string         `gorm:"type:varchar(50);unique;not null" json:"name"`
	Status              string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active','pending','archived')" json:"status"`
	DescriptionShort    string         `gorm:"type:varchar(200);not null" json:"description_short"`
	DescriptionOverview string         `gorm:"type:text;not null" json:"description_overview"`
	ImageURL            string         `gorm:"type:varchar(100)" json:"image_url"`
	TotalSectionQuizzes pq.Int64Array  `gorm:"type:integer[];default:'{}'" json:"total_section_quizzes"`
	CreatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ThemesOrLevelID     uint           `gorm:"not null" json:"themes_or_level_id"`
	CreatedBy           uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`

	SectionQuizzes []model.SectionQuizzesModel `gorm:"foreignKey:UnitID" json:"section_quizzes"`
}

func (UnitModel) TableName() string {
	return "units"
}

func (u *UnitModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalUnits(tx, u.ThemesOrLevelID)
}

func (u *UnitModel) AfterDelete(tx *gorm.DB) error {
	log.Printf("[HOOK] AfterDelete triggered for UnitID: %d", u.ID)

	var themesOrLevelID uint
	if err := tx.Unscoped().
		Model(&UnitModel{}).
		Select("themes_or_level_id").
		Where("id = ?", u.ID).
		Take(&themesOrLevelID).Error; err != nil {
		log.Println("[ERROR] Failed to fetch themes_or_level_id after delete:", err)
		return err
	}

	log.Printf("[HOOK] Fetched themes_or_level_id: %d for deleted UnitID: %d", themesOrLevelID, u.ID)
	return SyncTotalUnits(tx, themesOrLevelID)
}

func SyncTotalUnits(db *gorm.DB, themesOrLevelID uint) error {
	log.Println("[SERVICE] SyncTotalUnits - themesOrLevelID:", themesOrLevelID)

	err := db.Exec(`
		UPDATE themes_or_levels
		SET total_unit = (
			SELECT ARRAY_AGG(id ORDER BY id)
			FROM units
			WHERE themes_or_level_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, themesOrLevelID, themesOrLevelID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_unit:", err)
	}

	return err
}
