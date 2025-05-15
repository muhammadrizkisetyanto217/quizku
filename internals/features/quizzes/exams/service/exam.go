package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	// issuedcertificateservice "quizku/internals/features/certificates/issued_certificates/service"
	userSubcategoryModel "quizku/internals/features/lessons/subcategory/model"
	userThemeModel "quizku/internals/features/lessons/themes_or_levels/model"
	userUnitModel "quizku/internals/features/lessons/units/model"
)

type AttemptEvaluationData struct {
	Attempt         int `json:"attempt"`
	GradeEvaluation int `json:"grade_evaluation"`
}

func UpdateUserUnitFromExam(db *gorm.DB, userID uuid.UUID, examID uint, grade int) error {
	log.Println("[SERVICE] UpdateUserUnitFromExam - userID:", userID, "examID:", examID, "grade:", grade)
	if grade < 0 || grade > 100 {
		return fmt.Errorf("nilai grade tidak valid: %d", grade)
	}

	var unitID uint
	if err := db.Table("exams").Select("unit_id").Where("id = ?", examID).Scan(&unitID).Error; err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id:", examID)
		return err
	}

	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		return err
	}

	activityBonus := 0
	if userUnit.AttemptReading > 0 {
		activityBonus += 5
	}
	var evalData AttemptEvaluationData
	if len(userUnit.AttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.AttemptEvaluation, &evalData); err == nil && evalData.Attempt > 0 {
			activityBonus += 15
		}
	}
	var totalSections, completedSections int64
	_ = db.Table("section_quizzes").Where("unit_id = ?", unitID).Count(&totalSections).Error
	_ = db.Table("user_section_quizzes").Joins("JOIN section_quizzes ON user_section_quizzes.section_quizzes_id = section_quizzes.id").Where("user_section_quizzes.user_id = ? AND section_quizzes.unit_id = ?", userID, unitID).Count(&completedSections).Error
	if totalSections > 0 && totalSections == completedSections {
		activityBonus += 30
	}

	gradeResult := grade / 2
	if activityBonus > 0 {
		gradeResult = activityBonus + (grade / 2)
	}

	updates := map[string]interface{}{
		"grade_result": gradeResult,
		"is_passed":    gradeResult > 65,
		"updated_at":   time.Now(),
	}
	if grade > userUnit.GradeExam {
		updates["grade_exam"] = grade
	}
	if err := db.Model(&userUnit).Updates(updates).Error; err != nil {
		return err
	}

	if gradeResult > 65 {
		var themesID uint
		if err := db.Table("units").Select("themes_or_level_id").Where("id = ?", unitID).Scan(&themesID).Error; err != nil || themesID == 0 {
			return err
		}

		var userTheme userThemeModel.UserThemesOrLevelsModel
		if err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).First(&userTheme).Error; err != nil {
			return err
		}

		if userTheme.CompleteUnit == nil {
			userTheme.CompleteUnit = datatypes.JSONMap{}
		}
		userTheme.CompleteUnit[fmt.Sprintf("%d", unitID)] = fmt.Sprintf("%d", gradeResult)

		var expectedUnitIDs []int64
		if err := db.Table("units").Where("themes_or_level_id = ?", themesID).Pluck("id", &expectedUnitIDs).Error; err != nil {
			return err
		}

		matchCount := 0
		for _, id := range expectedUnitIDs {
			if _, ok := userTheme.CompleteUnit[fmt.Sprintf("%d", id)]; ok {
				matchCount++
			}
		}

		total := 0
		for _, id := range expectedUnitIDs {
			rawVal := userTheme.CompleteUnit[fmt.Sprintf("%d", id)]
			strVal := fmt.Sprintf("%v", rawVal)
			if g, err := strconv.Atoi(strVal); err == nil {
				total += g
			}
		}
		avg := 0
		if len(expectedUnitIDs) > 0 {
			avg = total / len(expectedUnitIDs)
		}

		updateFields := map[string]interface{}{
			"complete_unit": userTheme.CompleteUnit,
		}
		if matchCount == len(expectedUnitIDs) && len(expectedUnitIDs) > 0 {
			updateFields["grade_result"] = avg
		}
		if err := db.Model(&userTheme).Updates(updateFields).Error; err != nil {
			return err
		}

		var subcategoryID int
		if err := db.Table("themes_or_levels").Select("subcategories_id").Where("id = ?", themesID).Scan(&subcategoryID).Error; err != nil {
			return err
		}

		var userSub userSubcategoryModel.UserSubcategoryModel
		if err := db.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).First(&userSub).Error; err != nil {
			return err
		}
		if userSub.CompleteThemesOrLevels == nil {
			userSub.CompleteThemesOrLevels = datatypes.JSONMap{}
		}
		userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", themesID)] = fmt.Sprintf("%d", avg)

		var raw string
		if err := db.Table("subcategories").Select("total_themes_or_levels").Where("id = ?", subcategoryID).Scan(&raw).Error; err != nil {
			return err
		}
		var totalThemeIDs pq.Int64Array
		if err := totalThemeIDs.Scan(raw); err != nil {
			log.Println("[ERROR] Gagal parsing total_themes_or_levels:", err)
			return err
		}

		matchTheme := 0
		for _, id := range totalThemeIDs {
			if _, ok := userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", id)]; ok {
				matchTheme++
			}
		}

		totalSub := 0
		for _, id := range totalThemeIDs {
			rawVal := userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", id)]
			strVal := fmt.Sprintf("%v", rawVal)
			if g, err := strconv.Atoi(strVal); err == nil {
				totalSub += g
			}
		}
		avgSub := 0
		if len(totalThemeIDs) > 0 {
			avgSub = totalSub / len(totalThemeIDs)
		}

		updateSubFields := map[string]interface{}{
			"complete_themes_or_levels": userSub.CompleteThemesOrLevels,
		}
		if matchTheme == len(totalThemeIDs) && len(totalThemeIDs) > 0 {
			updateSubFields["grade_result"] = avgSub
		}
		if err := db.Model(&userSub).Updates(updateSubFields).Error; err != nil {
			return err
		}

		// if matchTheme == len(totalThemeIDs) && len(totalThemeIDs) > 0 {
		// 	if err := issuedcertificateservice.CreateIssuedCertificateIfEligible(db, userID, subcategoryID); err != nil {
		// 		return err
		// 	}
		// }

	}

	return nil
}

func CheckAndUnsetExamStatus(db *gorm.DB, userID uuid.UUID, examID uint) error {
	var unitID uint
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id untuk reset status:", examID)
		return err
	}

	var count int64
	err = db.Table("user_exams").
		Joins("JOIN exams ON exams.id = user_exams.exam_id").
		Where("user_exams.user_id = ? AND exams.unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		log.Println("[INFO] Reset nilai exam dan result karena tidak ada user_exams tersisa, user_id:", userID, "unit_id:", unitID)
		if err := db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"grade_exam":   0,
				"grade_result": 0,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			return err
		}

		var themesID uint
		err = db.Table("units").
			Select("themes_or_level_id").
			Where("id = ?", unitID).
			Scan(&themesID).Error
		if err != nil || themesID == 0 {
			log.Println("[ERROR] Gagal ambil themes_or_level_id dari unit:", unitID)
			return err
		}

		var userTheme userThemeModel.UserThemesOrLevelsModel
		err = db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
			First(&userTheme).Error
		if err != nil {
			log.Println("[WARNING] Tidak menemukan user_theme untuk reset complete_unit")
			return nil
		}

		if userTheme.CompleteUnit != nil {
			unitKey := fmt.Sprintf("%d", unitID)
			delete(userTheme.CompleteUnit, unitKey)

			shouldResetGrade := len(userTheme.CompleteUnit) == 0

			updateFields := map[string]interface{}{
				"complete_unit": userTheme.CompleteUnit,
			}
			if shouldResetGrade {
				updateFields["grade_result"] = 0
			}

			if err := db.Model(&userTheme).Updates(updateFields).Error; err != nil {
				log.Println("[ERROR] Gagal update user_theme saat reset:", err)
				return err
			}
		}
	}

	return nil
}
