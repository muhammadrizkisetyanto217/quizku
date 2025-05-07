package model

import (
	"time"
)

type TestExam struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Status    string    `gorm:"type:varchar(10);default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

func (TestExam) TableName() string {
	return "test_exam"
}
