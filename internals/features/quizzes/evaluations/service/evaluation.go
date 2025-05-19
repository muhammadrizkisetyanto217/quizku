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

// ✅ Semantik JSON progress evaluasi
type EvaluationAttemptSummary struct {
	EvaluationAttemptCount int `json:"evaluation_attempt_count"` // Jumlah percobaan evaluasi
	EvaluationFinalGrade   int `json:"evaluation_final_grade"`   // Nilai evaluasi tertinggi
}

// ✅ Update progress evaluasi user pada user_unit (tanpa refactor kolom DB)
func UpdateUserUnitFromEvaluation(db *gorm.DB, userID uuid.UUID, evaluationUnitID uint, gradePercentage int) error {
	var userUnit userUnitModel.UserUnitModel

	err := db.Select("attempt_evaluation").
		Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
		First(&userUnit).Error
	if err != nil {
		log.Printf("[WARNING] Gagal ambil user_unit: user_id=%s unit_id=%d err=%v", userID, evaluationUnitID, err)
		return err
	}

	var evalData EvaluationAttemptSummary
	if len(userUnit.UserUnitAttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.UserUnitAttemptEvaluation, &evalData); err != nil {
			log.Printf("[ERROR] Gagal decode JSON attempt_evaluation: %v", err)
			return err
		}
	}

	// ✅ Update jumlah attempt dan nilai tertinggi jika lebih baik
	evalData.EvaluationAttemptCount++
	if gradePercentage > evalData.EvaluationFinalGrade {
		evalData.EvaluationFinalGrade = gradePercentage
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

// ✅ Reset attempt_evaluation jika user tidak punya evaluasi untuk unit tersebut
func CheckAndUnsetEvaluationStatus(db *gorm.DB, userID uuid.UUID, evaluationUnitID uint) error {
	var count int64
	err := db.Table("user_evaluations").
		Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		log.Printf("[INFO] Reset attempt_evaluation karena tidak ada user_evaluations: user_id=%s unit_id=%d", userID, evaluationUnitID)
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, evaluationUnitID).
			Update("attempt_evaluation", datatypes.JSON([]byte("null"))).Error
	}

	return nil
}
