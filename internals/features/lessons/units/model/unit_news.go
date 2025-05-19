package model

import (
	"time"

	"gorm.io/gorm"
)

type UnitNewsModel struct {
	UnitNewsID          uint           `gorm:"column:unit_news_id;primaryKey;autoIncrement" json:"unit_news_id"`
	UnitNewsUnitID      uint           `gorm:"column:unit_news_unit_id;not null" json:"unit_news_unit_id"`
	UnitNewsTitle       string         `gorm:"column:unit_news_title;type:varchar(255);not null" json:"unit_news_title"`
	UnitNewsDescription string         `gorm:"column:unit_news_description;type:text;not null" json:"unit_news_description"`
	UnitNewsIsPublic    bool           `gorm:"column:unit_news_is_public;default:true" json:"unit_news_is_public"`

	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (UnitNewsModel) TableName() string {
	return "unit_news"
}
