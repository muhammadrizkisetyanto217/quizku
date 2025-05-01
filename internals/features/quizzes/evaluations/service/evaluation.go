package service

import (
	"encoding/json"
	"log"
	userUnitModel "quizku/internals/features/lessons/units/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AttemptEvaluationData struct {
	Attempt         int `json:"attempt"`
	GradeEvaluation int `json:"grade_evaluation"`
}

func UpdateUserUnitFromEvaluation(db *gorm.DB, userID uuid.UUID, unitID uint, gradePercentage int) error {
	var userUnit userUnitModel.UserUnitModel

	// Ambil hanya kolom yang diperlukan
	err := db.Select("attempt_evaluation").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		First(&userUnit).Error
	if err != nil {
		log.Printf("[WARNING] Gagal ambil user_unit: user_id=%s unit_id=%d err=%v", userID, unitID, err)
		return err
	}

	// Decode JSON
	var evalData AttemptEvaluationData
	if len(userUnit.AttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.AttemptEvaluation, &evalData); err != nil {
			log.Printf("[ERROR] Gagal decode JSON attempt_evaluation: %v", err)
			return err
		}
	}

	// Update data attempt dan nilai
	evalData.Attempt++
	if gradePercentage > evalData.GradeEvaluation {
		evalData.GradeEvaluation = gradePercentage
	}

	// Encode JSON kembali
	jsonEval, err := json.Marshal(evalData)
	if err != nil {
		log.Printf("[ERROR] Gagal encode JSON: %v", err)
		return err
	}

	// Simpan hanya attempt_evaluation
	updateData := map[string]interface{}{
		"attempt_evaluation": datatypes.JSON(jsonEval),
		"updated_at":         time.Now(),
	}

	return db.Model(&userUnitModel.UserUnitModel{}).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Updates(updateData).Error
}

func CheckAndUnsetEvaluationStatus(db *gorm.DB, userID uuid.UUID, unitID uint) error {
	var count int64
	err := db.Table("user_evaluations").
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Update("attempt_evaluation", 0).Error
	}

	return nil
}
