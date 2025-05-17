package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	userUnitModel "quizku/internals/features/lessons/units/model"
)

// ✅ Gaya semantik untuk isi JSON progress evaluasi
type AttemptEvaluationData struct {
	UserUnitEvaluationAttempt         int `json:"user_unit_evaluation_attempt"`
	UserUnitEvaluationGradeEvaluation int `json:"user_unit_evaluation_grade_evaluation"`
}

// ✅ Update progress evaluasi pada tabel user_units (tanpa refactor kolom DB)
func UpdateUserUnitFromEvaluation(db *gorm.DB, userID uuid.UUID, evaluationUnitID uint, gradePercentage int) error {
	var userUnit userUnitModel.UserUnitModel

	err := db.Select("attempt_evaluation").
		Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
		First(&userUnit).Error
	if err != nil {
		log.Printf("[WARNING] Gagal ambil user_unit: user_id=%s unit_id=%d err=%v", userID, evaluationUnitID, err)
		return err
	}

	var evalData AttemptEvaluationData
	if len(userUnit.AttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.AttemptEvaluation, &evalData); err != nil {
			log.Printf("[ERROR] Gagal decode JSON attempt_evaluation: %v", err)
			return err
		}
	}

	// Update attempt dan nilai jika lebih baik
	evalData.UserUnitEvaluationAttempt++
	if gradePercentage > evalData.UserUnitEvaluationGradeEvaluation {
		evalData.UserUnitEvaluationGradeEvaluation = gradePercentage
	}

	encoded, err := json.Marshal(evalData)
	if err != nil {
		log.Printf("[ERROR] Gagal encode JSON attempt_evaluation: %v", err)
		return err
	}

	updateData := map[string]interface{}{
		"attempt_evaluation": datatypes.JSON(encoded),
		"updated_at":         time.Now(),
	}

	return db.Model(&userUnitModel.UserUnitModel{}).
		Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
		Updates(updateData).Error
}

// ✅ Reset field jika tidak ada user_evaluations untuk unit ini
func CheckAndUnsetEvaluationStatus(db *gorm.DB, userID uuid.UUID, evaluationUnitID uint) error {
	var count int64
	err := db.Table("user_evaluations").
		Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
			Update("attempt_evaluation", 0).Error
	}

	return nil
}
