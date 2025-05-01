package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

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
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id:", examID)
		return err
	}

	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("[WARNING] user_unit tidak ditemukan saat UpdateUserUnitFromExam, user_id:", userID, "unit_id:", unitID)
		}
		return err
	}

	// Hitung bonus
	activityBonus := 0
	if userUnit.AttemptReading > 0 {
		activityBonus += 5
	}
	var evalData AttemptEvaluationData
	if len(userUnit.AttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.AttemptEvaluation, &evalData); err != nil {
			log.Println("[ERROR] Gagal decode AttemptEvaluation JSON:", err)
		} else if evalData.Attempt > 0 {
			activityBonus += 15
		}
	}

	var totalSections, completedSections int64
	_ = db.Table("section_quizzes").Where("unit_id = ?", unitID).Count(&totalSections).Error
	_ = db.Table("user_section_quizzes").
		Joins("JOIN section_quizzes ON user_section_quizzes.section_quizzes_id = section_quizzes.id").
		Where("user_section_quizzes.user_id = ? AND section_quizzes.unit_id = ?", userID, unitID).
		Count(&completedSections).Error

	if totalSections > 0 && totalSections == completedSections {
		activityBonus += 30
	}

	// Final grade_result
	var gradeResult int
	if activityBonus == 0 {
		gradeResult = grade / 2
	} else {
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

	// ✅ Tambahkan ke complete_unit jika lulus
	if gradeResult > 65 {
		var themesID uint
		err := db.Table("units").
			Select("themes_or_level_id").
			Where("id = ?", unitID).
			Scan(&themesID).Error
		if err != nil || themesID == 0 {
			log.Println("[ERROR] Gagal ambil themes_or_level_id dari unit:", unitID)
			return err
		}

		var userTheme userThemeModel.UserThemesOrLevelsModel
		if err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
			First(&userTheme).Error; err != nil {
			log.Println("[ERROR] Tidak menemukan user_theme record")
			return err
		}

		if userTheme.CompleteUnit == nil {
			userTheme.CompleteUnit = datatypes.JSONMap{}
		}

		unitKey := fmt.Sprintf("%d", unitID)
		gradeStr := fmt.Sprintf("%d", gradeResult)
		userTheme.CompleteUnit[unitKey] = gradeStr

		matchCount := 0
		for _, expectedID := range userTheme.TotalUnit {
			if _, ok := userTheme.CompleteUnit[fmt.Sprintf("%d", expectedID)]; ok {
				matchCount++
			}
		}

		if matchCount == len(userTheme.TotalUnit) && len(userTheme.TotalUnit) > 0 {
			var totalGrade int
			for _, expectedID := range userTheme.TotalUnit {
				if str, ok := userTheme.CompleteUnit[fmt.Sprintf("%d", expectedID)].(string); ok {
					if g, err := strconv.Atoi(str); err == nil {
						totalGrade += g
					}
				}
			}
			avgGrade := totalGrade / len(userTheme.TotalUnit)
			userTheme.GradeResult = avgGrade

			if err := db.Model(&userTheme).Updates(map[string]interface{}{
				"complete_unit": userTheme.CompleteUnit,
				"grade_result":  avgGrade,
			}).Error; err != nil {
				log.Println("[ERROR] Gagal update complete_unit dan grade_result:", err)
				return err
			}

			// ✅ Lanjut update user_subcategory
			var subcategoryID int
			err = db.Table("themes_or_levels").
				Select("subcategories_id").
				Where("id = ?", themesID).
				Scan(&subcategoryID).Error
			if err != nil || subcategoryID == 0 {
				log.Println("[ERROR] Gagal ambil subcategory_id dari themes_id:", themesID)
				return err
			}

			var userSub userSubcategoryModel.UserSubcategoryModel
			err = db.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).
				First(&userSub).Error
			if err != nil {
				log.Println("[ERROR] Tidak menemukan user_subcategory record")
				return err
			}

			if userSub.CompleteThemesOrLevels == nil {
				userSub.CompleteThemesOrLevels = datatypes.JSONMap{}
			}
			userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", themesID)] = fmt.Sprintf("%d", avgGrade)

			matchCount := 0
			for _, themeID := range userSub.TotalThemesOrLevels {
				if _, ok := userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", themeID)]; ok {
					matchCount++
				}
			}

			if matchCount == len(userSub.TotalThemesOrLevels) && len(userSub.TotalThemesOrLevels) > 0 {
				var totalThemeGrade int
				for _, themeID := range userSub.TotalThemesOrLevels {
					if str, ok := userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", themeID)].(string); ok {
						if g, err := strconv.Atoi(str); err == nil {
							totalThemeGrade += g
						}
					}
				}
				avgThemeGrade := totalThemeGrade / len(userSub.TotalThemesOrLevels)
				userSub.GradeResult = avgThemeGrade

				if err := db.Model(&userSub).Updates(map[string]interface{}{
					"complete_themes_or_levels": userSub.CompleteThemesOrLevels,
					"grade_result":              avgThemeGrade,
				}).Error; err != nil {
					log.Println("[ERROR] Gagal update complete_themes_or_levels dan grade_result:", err)
					return err
				}
			} else {
				if err := db.Model(&userSub).
					Update("complete_themes_or_levels", userSub.CompleteThemesOrLevels).Error; err != nil {
					log.Println("[ERROR] Gagal update complete_themes_or_levels:", err)
					return err
				}
			}
		} else {
			// Belum semua unit selesai → update complete_unit saja
			if err := db.Model(&userTheme).
				Update("complete_unit", userTheme.CompleteUnit).Error; err != nil {
				log.Println("[ERROR] Gagal update complete_unit:", err)
				return err
			}
		}
	}

	return nil
}

// ✅ Final: CheckAndUnsetExamStatus
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

		// ✅ Reset nilai di user_unit
		if err := db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"grade_exam":   0,
				"grade_result": 0,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			return err
		}

		// ✅ Hapus dari complete_unit di user_themes_or_levels
		var themesID uint
		err := db.Table("units").
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
			return nil // aman, tidak perlu error kalau tidak ketemu
		}

		if userTheme.CompleteUnit != nil {
			unitKey := fmt.Sprintf("%d", unitID)
			delete(userTheme.CompleteUnit, unitKey)

			if err := db.Model(&userTheme).
				Update("complete_unit", userTheme.CompleteUnit).Error; err != nil {
				log.Println("[ERROR] Gagal hapus unit dari complete_unit:", err)
				return err
			}
		}
	}

	return nil
}
