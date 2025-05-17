package model

import (
	"time"

	"github.com/google/uuid"
)

type UserEvaluationModel struct {
	UserEvaluationID              uint      `gorm:"column:user_evaluation_id;primaryKey;autoIncrement" json:"user_evaluation_id"`
	UserEvaluationUserID          uuid.UUID `gorm:"column:user_evaluation_user_id;type:uuid;not null;index:idx_user_eval_eval,priority:1;index:idx_user_eval_unit,priority:1" json:"user_evaluation_user_id"`
	UserEvaluationEvaluationID    uint      `gorm:"column:user_evaluation_evaluation_id;not null;index:idx_user_eval_eval,priority:2" json:"user_evaluation_evaluation_id"`
	UserEvaluationUnitID          uint      `gorm:"column:user_evaluation_unit_id;not null;index:idx_user_eval_unit,priority:2" json:"user_evaluation_unit_id"`
	UserEvaluationAttempt         int       `gorm:"column:user_evaluation_attempt;default:1;not null" json:"user_evaluation_attempt"`
	UserEvaluationPercentageGrade int       `gorm:"column:user_evaluation_percentage_grade;default:0;not null" json:"user_evaluation_percentage_grade"`
	UserEvaluationTimeDuration    int       `gorm:"column:user_evaluation_time_duration;default:0;not null" json:"user_evaluation_time_duration"`
	UserEvaluationPoint           int       `gorm:"column:user_evaluation_point;default:0;not null" json:"user_evaluation_point"`
	CreatedAt                     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt                     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ✅ Mapping nama tabel
func (UserEvaluationModel) TableName() string {
	return "user_evaluations"
}
