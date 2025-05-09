package model

import (
	"time"
)

const (
	TargetTypeQuiz       = 1
	TargetTypeEvaluation = 2
	TargetTypeExam       = 3
	TargetTypeTest       = 4
)

var TargetTypeMap = map[int]string{
	TargetTypeQuiz:       "quiz",
	TargetTypeEvaluation: "evaluation",
	TargetTypeExam:       "exam",
	TargetTypeTest:       "test",
}

var TargetTypeNameToInt = map[string]int{
	"quiz":       TargetTypeQuiz,
	"evaluation": TargetTypeEvaluation,
	"exam":       TargetTypeExam,
	"test":       TargetTypeTest,
}

type QuestionLink struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	QuestionID int       `gorm:"not null" json:"question_id"`
	TargetType int       `gorm:"not null;check:target_type IN (1,2,3,4)" json:"target_type"`
	TargetID   int       `gorm:"not null" json:"target_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Method untuk mapping angka ke nama (opsional)
func (q *QuestionLink) TargetTypeName() string {
	if name, ok := TargetTypeMap[q.TargetType]; ok {
		return name
	}
	return "unknown"
}

// ⬅️ Ini yang kamu minta
func (QuestionLink) TableName() string {
	return "question_links"
}
