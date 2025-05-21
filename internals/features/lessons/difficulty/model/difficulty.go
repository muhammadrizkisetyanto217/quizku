package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DifficultyModel struct {
	DifficultyID               uint           `gorm:"column:difficulty_id;primaryKey;autoIncrement" json:"difficulty_id"`
	DifficultyName             string         `gorm:"column:difficulty_name;type:varchar(255);not null" json:"difficulty_name"`
	DifficultyStatus           string         `gorm:"column:difficulty_status;type:varchar(10);default:'pending';check:difficulty_status IN ('active', 'pending', 'archived')" json:"difficulty_status"`
	DifficultyDescriptionShort string         `gorm:"column:difficulty_description_short;type:varchar(200)" json:"difficulty_description_short"`
	DifficultyDescriptionLong  string         `gorm:"column:difficulty_description_long;type:varchar(3000)" json:"difficulty_description_long"`
	DifficultyTotalCategories  pq.Int64Array  `gorm:"column:difficulty_total_categories;type:integer[];default:'{}'" json:"difficulty_total_categories"`
	DifficultyImageURL         string         `gorm:"column:difficulty_image_url;type:varchar(255)" json:"difficulty_image_url"`
	DifficultyUpdateNews       []byte         `gorm:"column:difficulty_update_news;type:jsonb" json:"difficulty_update_news,omitempty"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (DifficultyModel) TableName() string {
	return "difficulties"
}
