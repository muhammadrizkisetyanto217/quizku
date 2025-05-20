package model

import (
	"time"

	"gorm.io/gorm"
)

type CategoryNewsModel struct {
	CategoryNewsID          uint           `gorm:"primaryKey;column:category_news_id" json:"category_news_id"`
	CategoryNewsTitle       string         `gorm:"type:varchar(255);not null;column:category_news_title" json:"category_news_title"`
	CategoryNewsDescription string         `gorm:"type:text;not null;column:category_news_description" json:"category_news_description"`
	CategoryNewsIsPublic    bool           `gorm:"default:true;column:category_news_is_public" json:"category_news_is_public"`
	CategoryNewsCategoryID  int            `gorm:"column:category_news_category_id" json:"category_news_category_id"`
	CreatedAt               time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (CategoryNewsModel) TableName() string {
	return "category_news"
}
