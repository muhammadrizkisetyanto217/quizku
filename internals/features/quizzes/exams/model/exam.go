package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ExamModel struct {
	ExamID               uint          `gorm:"column:exam_id;primaryKey" json:"exam_id"`
	ExamName             string        `gorm:"column:exam_name;size:100;not null" json:"exam_name" validate:"required,max=100"`
	ExamStatus           string        `gorm:"column:exam_status;type:varchar(20);default:'pending';check:exam_status IN ('active','pending','archived')" json:"exam_status" validate:"required,oneof=active pending archived"`
	ExamTotalQuestionIDs pq.Int64Array `gorm:"column:exam_total_question_ids;type:integer[];default:'{}'" json:"exam_total_question_ids"`
	ExamIconURL          *string       `gorm:"column:exam_icon_url;size:255" json:"exam_icon_url,omitempty" validate:"omitempty,url"`

	ExamUnitID    uint      `gorm:"column:exam_unit_id;not null" json:"exam_unit_id" validate:"required"`
	ExamCreatedBy uuid.UUID `gorm:"column:exam_created_by;type:uuid;not null;constraint:OnDelete:CASCADE" json:"exam_created_by"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// TableName overrides the default table name
func (ExamModel) TableName() string {
	return "exams"
}
