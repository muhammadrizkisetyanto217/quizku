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
	QuestionLinkID         int       `gorm:"column:question_link_id;primaryKey" json:"question_link_id"`
	QuestionLinkQuestionID int       `gorm:"column:question_link_question_id;not null" json:"question_link_question_id"`
	QuestionLinkTargetType int       `gorm:"column:question_link_target_type;not null;check:question_link_target_type IN (1,2,3,4)" json:"question_link_target_type"`
	QuestionLinkTargetID   int       `gorm:"column:question_link_target_id;not null" json:"question_link_target_id"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// Opsional: Helper untuk menampilkan nama tipe
func (q *QuestionLink) TargetTypeName() string {
	if name, ok := TargetTypeMap[q.QuestionLinkTargetType]; ok {
		return name
	}
	return "unknown"
}

func (QuestionLink) TableName() string {
	return "question_links"
}
