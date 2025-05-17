package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// EvaluationModel merepresentasikan tabel evaluations di database
type EvaluationModel struct {
	EvaluationID            uint          `gorm:"column:evaluation_id;primaryKey" json:"evaluation_id"`
	EvaluationName          string        `gorm:"column:evaluation_name;size:50;not null" json:"evaluation_name"`
	EvaluationStatus        string        `gorm:"column:evaluation_status;type:varchar(10);default:'pending';check:evaluation_status IN ('active', 'pending', 'archived')" json:"evaluation_status" validate:"required,oneof=active pending archived"`
	EvaluationTotalQuestion pq.Int64Array `gorm:"column:evaluation_total_question;type:integer[];default:'{}'" json:"evaluation_total_question"`
	EvaluationIconURL       *string       `gorm:"column:evaluation_icon_url;size:100" json:"evaluation_icon_url,omitempty" validate:"omitempty,url"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	EvaluationUnitID    uint      `gorm:"column:evaluation_unit_id" json:"evaluation_unit_id" validate:"required"`
	EvaluationCreatedBy uuid.UUID `gorm:"column:evaluation_created_by;type:uuid;not null;constraint:OnDelete:CASCADE" json:"evaluation_created_by"`
}

// TableName mengatur nama tabel agar sesuai dengan skema database
func (EvaluationModel) TableName() string {
	return "evaluations"
}
