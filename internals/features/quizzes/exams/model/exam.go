package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ExamModel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NameExams     string         `gorm:"size:50;not null" json:"name_exams" validate:"required,max=50"`
	Status        string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status" validate:"required,oneof=active pending archived"`
	TotalQuestion pq.Int64Array  `gorm:"type:integer[];default:'{}'" json:"total_question"`
	IconURL       *string        `gorm:"size:100" json:"icon_url,omitempty" validate:"omitempty,url"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UnitID        uint           `json:"unit_id" validate:"required"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

// TableName mengatur nama tabel agar sesuai dengan skema database
func (ExamModel) TableName() string {
	return "exams"
}
