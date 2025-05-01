package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Fungsi untuk memperbarui UserSectionQuizzesModel
// UserSectionQuizzesModel menyimpan daftar kuis yang telah diselesaikan dalam suatu section
type UserCategoryModel struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	UserID           uuid.UUID     `gorm:"type:uuid;not null;index:idx_user_category_user_cat,unique" json:"user_id"`
	CategoryID       int           `gorm:"not null;index:idx_user_category_user_cat,unique" json:"category_id"`
	CompleteCategory pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"complete_category"`
	TotalCategory    pq.Int64Array `gorm:"type:integer[];default:'{}'" json:"total_category"`
	GradeResult      int           `gorm:"default:0" json:"grade_result"`
	CreatedAt        time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserCategoryModel) TableName() string {
	return "user_category"
}
