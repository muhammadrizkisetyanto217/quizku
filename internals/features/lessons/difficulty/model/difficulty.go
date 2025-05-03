package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DifficultyModel struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"size:255;not null" json:"name"`
	Status           string         `gorm:"size:10;default:'pending';check:status IN ('active', 'pending', 'archived')" json:"status"`
	DescriptionShort string         `gorm:"size:200" json:"description_short"`
	DescriptionLong  string         `gorm:"size:3000" json:"description_long"`
	TotalCategories  pq.Int64Array  `json:"total_categories" gorm:"type:integer[];default:'{}'"`
	ImageURL         string         `gorm:"size:100" json:"image_url"`
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (DifficultyModel) TableName() string {
	return "difficulties"
}