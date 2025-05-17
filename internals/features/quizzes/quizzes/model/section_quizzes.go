package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SectionQuizzesModel struct {
	SectionQuizzesID           uint           `gorm:"column:section_quizzes_id;primaryKey;autoIncrement" json:"section_quizzes_id"`
	SectionQuizzesName         string         `gorm:"column:section_quizzes_name;size:50;not null" json:"section_quizzes_name"`
	SectionQuizzesStatus       string         `gorm:"column:section_quizzes_status;type:varchar(10);default:'pending';check:section_quizzes_status IN ('active', 'pending', 'archived')" json:"section_quizzes_status"`
	SectionQuizzesMaterials    string         `gorm:"column:section_quizzes_materials;type:text;not null" json:"section_quizzes_materials"`
	SectionQuizzesIconURL      string         `gorm:"column:section_quizzes_icon_url;size:100" json:"section_quizzes_icon_url"`
	SectionQuizzesTotalQuizzes pq.Int64Array  `gorm:"column:section_quizzes_total_quizzes;type:integer[];default:'{}'" json:"section_quizzes_total_quizzes"`
	CreatedAt                  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt                  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
	SectionQuizzesUnitID       uint           `gorm:"column:section_quizzes_unit_id;not null" json:"section_quizzes_unit_id"`
	SectionQuizzesCreatedBy    uuid.UUID      `gorm:"column:section_quizzes_created_by;type:uuid;not null" json:"section_quizzes_created_by"`

	Quizzes []QuizModel `gorm:"foreignKey:QuizSectionQuizzesID" json:"quizzes"`
}

func (SectionQuizzesModel) TableName() string {
	return "section_quizzes"
}

// ✅ AfterSave: Sinkronisasi array ID section_quizzes setelah simpan
func (s *SectionQuizzesModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalSectionQuizzes(tx, s.SectionQuizzesUnitID)
}

// ✅ AfterDelete: Sinkronisasi array ID section_quizzes setelah dihapus
func (s *SectionQuizzesModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalSectionQuizzes(tx, s.SectionQuizzesUnitID)
}

// ✅ Sinkronisasi field total_section_quizzes di tabel units
func SyncTotalSectionQuizzes(db *gorm.DB, unitID uint) error {
	log.Println("[SERVICE] SyncTotalSectionQuizzes - unitID:", unitID)

	err := db.Exec(`
		UPDATE units
		SET total_section_quizzes = (
			SELECT ARRAY_AGG(section_quizzes_id ORDER BY section_quizzes_id)
			FROM section_quizzes
			WHERE section_quizzes_unit_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, unitID, unitID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_section_quizzes:", err)
	}
	return err
}
