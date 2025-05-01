package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SectionQuizzesModel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	NameQuizzes      string         `gorm:"size:50;not null" json:"name_quizzes"`
	Status           string         `gorm:"size:10;default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	MaterialsQuizzes string         `gorm:"type:text;not null" json:"materials_quizzes"`
	IconURL          string         `gorm:"size:100" json:"icon_url"`
	TotalQuizzes     pq.Int64Array  `gorm:"type:integer[];default:'{}'" json:"total_quizzes"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy        uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
	UnitID           uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"unit_id"`
	Quizzes          []QuizModel    `gorm:"foreignKey:SectionQuizID" json:"quizzes"`
}

func (SectionQuizzesModel) TableName() string {
	return "section_quizzes"
}

// ✅ AfterSave: Sinkronisasi array ID section_quizzes setelah simpan
func (s *SectionQuizzesModel) AfterSave(tx *gorm.DB) error {
	return SyncTotalSectionQuizzes(tx, s.UnitID)
}

// ✅ AfterDelete: Sinkronisasi array ID section_quizzes setelah dihapus
func (s *SectionQuizzesModel) AfterDelete(tx *gorm.DB) error {
	return SyncTotalSectionQuizzes(tx, s.UnitID)
}

// ✅ Sinkronisasi field total_section_quizzes di tabel units
func SyncTotalSectionQuizzes(db *gorm.DB, unitID uint) error {
	log.Println("[SERVICE] SyncTotalSectionQuizzes - unitID:", unitID)

	err := db.Exec(`
		UPDATE units
		SET total_section_quizzes = (
			SELECT ARRAY_AGG(id ORDER BY id)
			FROM section_quizzes
			WHERE unit_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, unitID, unitID).Error

	if err != nil {
		log.Println("[ERROR] Failed to sync total_section_quizzes:", err)
	}
	return err
}
