package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingModel struct {
	ReadingID              uint   `gorm:"column:reading_id;primaryKey;autoIncrement" json:"reading_id"`
	ReadingTitle           string `gorm:"column:reading_title;type:varchar(50);unique;not null" json:"reading_title"`
	ReadingStatus          string `gorm:"column:reading_status;type:varchar(10);default:'pending';check:reading_status IN ('active', 'pending', 'archived')" json:"reading_status"`
	ReadingDescriptionLong string `gorm:"column:reading_description_long;type:text;not null" json:"reading_description_long"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`

	ReadingUnitID    uint      `gorm:"column:reading_unit_id" json:"reading_unit_id"`
	ReadingCreatedBy uuid.UUID `gorm:"column:reading_created_by;type:uuid;not null;constraint:OnDelete:CASCADE" json:"reading_created_by"`
}

// ✅ Mapping nama tabel
func (ReadingModel) TableName() string {
	return "readings"
}
