package model

import (
	"time"

	"gorm.io/gorm"
)

type CategoryNewsModel struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"type:varchar(255);not null"`
	Description  string         `gorm:"type:text;not null"`
	IsPublic     bool           `gorm:"default:true"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CategoryID int            `json:"category_id"`
}

func (CategoryNewsModel) TableName() string {
	return "categories_news"
}