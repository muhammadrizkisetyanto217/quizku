package model

import (
	"time"

	"gorm.io/gorm"
)

type SubcategoryNewsModel struct {
	SubcategoryNewsID              uint           `gorm:"primaryKey;column:subcategory_news_id" json:"subcategory_news_id"`
	SubcategoryNewsTitle           string         `gorm:"type:varchar(255);not null;column:subcategory_news_title" json:"subcategory_news_title"`
	SubcategoryNewsDescription     string         `gorm:"type:text;not null;column:subcategory_news_description" json:"subcategory_news_description"`
	SubcategoryNewsIsPublic        bool           `gorm:"default:true;column:subcategory_news_is_public" json:"subcategory_news_is_public"`
	SubcategoryNewsSubcategoryID   uint           `gorm:"not null;column:subcategory_news_subcategory_id" json:"subcategory_news_subcategory_id"`
	CreatedAt                      time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt                      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (SubcategoryNewsModel) TableName() string {
	return "subcategory_news"
}
