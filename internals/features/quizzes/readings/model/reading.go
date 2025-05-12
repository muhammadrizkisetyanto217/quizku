package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"type:varchar(50);unique;not null" json:"title"`
	Status          string         `gorm:"type:varchar(10);default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	DescriptionLong string         `gorm:"type:text;not null" json:"description_long"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UnitID          uint           `json:"unit_id"`
	CreatedBy       uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"created_by"`
}

func (ReadingModel) TableName() string {
	return "readings"
}
