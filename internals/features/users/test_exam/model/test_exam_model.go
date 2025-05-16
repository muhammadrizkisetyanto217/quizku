package model

import "time"

type TestExam struct {
	TestExamID     uint      `gorm:"column:test_exam_id;primaryKey" json:"test_exam_id"`
	TestExamName   string    `gorm:"column:test_exam_name;type:varchar(50);not null" json:"test_exam_name"`
	TestExamStatus string    `gorm:"column:test_exam_status;type:varchar(10);default:'pending'" json:"test_exam_status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (TestExam) TableName() string {
	return "test_exams"
}
