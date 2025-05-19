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

	issuedcertificateservice "quizku/internals/features/certificates/user_certificates/service"
	userSubcategoryModel "quizku/internals/features/lessons/subcategories/model"
	userThemeModel "quizku/internals/features/lessons/themes_or_levels/model"
	userUnitModel "quizku/internals/features/lessons/units/model"
)

func UpdateUserUnitFromExam(db *gorm.DB, userID uuid.UUID, examID uint, submittedExamGrade int) error {
	type EvaluationAttemptSummary struct {
		EvaluationAttemptCount int `json:"evaluation_attempt_count"`
		EvaluationFinalGrade   int `json:"evaluation_final_grade"`
	}

	log.Println("[SERVICE] UpdateUserUnitFromExam - userID:", userID, "examID:", examID, "submittedExamGrade:", submittedExamGrade)
	if submittedExamGrade < 0 || submittedExamGrade > 100 {
		return fmt.Errorf("nilai submittedExamGrade tidak valid: %d", submittedExamGrade)
	}

	// Ambil unit_id berdasarkan exam
	var unitID uint
	if err := db.Table("exams").Select("exam_unit_id").Where("exam_id = ?", examID).Scan(&unitID).Error; err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil exam_unit_id:", err)
		return fmt.Errorf("unit_id tidak ditemukan untuk examID %d", examID)
	}

	// Ambil record user_unit
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_unit_user_id = ? AND user_unit_unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		return err
	}

	// Hitung bonus berdasarkan aktivitas
	activityBonus := 0
	if userUnit.UserUnitAttemptReading > 0 {
		activityBonus += 5
	}
	var evalData EvaluationAttemptSummary
	if len(userUnit.UserUnitAttemptEvaluation) > 0 {
		if err := json.Unmarshal(userUnit.UserUnitAttemptEvaluation, &evalData); err == nil && evalData.EvaluationAttemptCount > 0 {
			activityBonus += 15
		}
	}

	// Cek penyelesaian section_quizzes
	var totalSections, completedSections int64
	db.Table("section_quizzes").Where("section_quizzes_unit_id = ?", unitID).Count(&totalSections)
	db.Table("user_section_quizzes").
		Joins("JOIN section_quizzes ON user_section_quizzes.user_section_quizzes_section_quizzes_id = section_quizzes.section_quizzes_id").
		Where("user_section_quizzes_user_id = ? AND section_quizzes.section_quizzes_unit_id = ?", userID, unitID).
		Count(&completedSections)
	if totalSections > 0 && totalSections == completedSections {
		activityBonus += 30
	}

	// Hitung grade_result untuk LEVEL UNIT
	unitFinalGrade := (submittedExamGrade / 2) + activityBonus
	updateUnit := map[string]interface{}{
		"user_unit_grade_result": unitFinalGrade,
		"user_unit_is_passed":    unitFinalGrade > 65,
		"updated_at":             time.Now(),
	}
	if submittedExamGrade > userUnit.UserUnitGradeExam {
		updateUnit["user_unit_grade_exam"] = submittedExamGrade
	}
	if err := db.Model(&userUnit).Updates(updateUnit).Error; err != nil {
		return err
	}
	if unitFinalGrade <= 65 {
		return nil
	}

	// Ambil theme ID dari unit
	var themesID uint
	if err := db.Table("units").Select("themes_or_level_id").Where("unit_id = ?", unitID).Scan(&themesID).Error; err != nil || themesID == 0 {
		return fmt.Errorf("themes_id tidak ditemukan dari unit")
	}

	// Update progress user_themes_or_levels
	var userTheme userThemeModel.UserThemesOrLevelsModel
	if err := db.Where("user_theme_user_id = ? AND user_theme_themes_or_level_id = ?", userID, themesID).First(&userTheme).Error; err != nil {
		return err
	}
	if userTheme.UserThemeCompleteUnit == nil {
		userTheme.UserThemeCompleteUnit = datatypes.JSONMap{}
	}
	userTheme.UserThemeCompleteUnit[fmt.Sprintf("%d", unitID)] = fmt.Sprintf("%d", unitFinalGrade)

	var expectedUnitIDs []int64
	if err := db.Table("units").Where("themes_or_level_id = ?", themesID).Pluck("unit_id", &expectedUnitIDs).Error; err != nil {
		return err
	}
	matchCount, total := 0, 0
	for _, id := range expectedUnitIDs {
		if val, ok := userTheme.UserThemeCompleteUnit[fmt.Sprintf("%d", id)]; ok {
			matchCount++
			if g, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
				total += g
			}
		}
	}
	userThemeFinalGradeResult := 0
	if len(expectedUnitIDs) > 0 {
		userThemeFinalGradeResult = total / len(expectedUnitIDs)
	}

	themeUpdateFields := map[string]interface{}{
		"user_theme_complete_unit": userTheme.UserThemeCompleteUnit,
	}
	if matchCount == len(expectedUnitIDs) && len(expectedUnitIDs) > 0 {
		themeUpdateFields["user_theme_grade_result"] = userThemeFinalGradeResult
	}
	if err := db.Model(&userTheme).Updates(themeUpdateFields).Error; err != nil {
		return err
	}

	// Ambil subcategory dari theme
	var subcategoryID int
	if err := db.Table("themes_or_levels").Select("themes_or_level_subcategory_id").Where("themes_or_level_id = ?", themesID).Scan(&subcategoryID).Error; err != nil {
		return err
	}

	// Update progress user_subcategory
	var userSub userSubcategoryModel.UserSubcategoryModel
	if err := db.Where("user_subcategory_user_id = ? AND user_subcategory_subcategory_id = ?", userID, subcategoryID).First(&userSub).Error; err != nil {
		return err
	}
	if userSub.UserSubcategoryCompleteThemesOrLevels == nil {
		userSub.UserSubcategoryCompleteThemesOrLevels = datatypes.JSONMap{}
	}
	userSub.UserSubcategoryCompleteThemesOrLevels[fmt.Sprintf("%d", themesID)] = fmt.Sprintf("%d", userThemeFinalGradeResult)

	var raw string
	if err := db.Table("subcategories").Select("subcategory_total_themes_or_levels").Where("subcategory_id = ?", subcategoryID).Scan(&raw).Error; err != nil {
		return err
	}
	var totalThemeIDs pq.Int64Array
	if err := totalThemeIDs.Scan(raw); err != nil {
		log.Println("[ERROR] Gagal parsing total_themes_or_levels:", err)
		return err
	}

	matchTheme, totalSub := 0, 0
	for _, id := range totalThemeIDs {
		if val, ok := userSub.UserSubcategoryCompleteThemesOrLevels[fmt.Sprintf("%d", id)]; ok {
			matchTheme++
			if g, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
				totalSub += g
			}
		}
	}
	userSubcategoryFinalGradeResult := 0
	if len(totalThemeIDs) > 0 {
		userSubcategoryFinalGradeResult = totalSub / len(totalThemeIDs)
	}

	// Ambil versi sertifikat terbaru
	var issuedVersion int
	row := db.Table("certificate_versions").
		Where("certificate_subcategory_id = ?", subcategoryID).
		Select("certificate_version_number").
		Order("certificate_version_number DESC").
		Limit(1).
		Row()
	if err := row.Scan(&issuedVersion); err != nil {
		log.Printf("[INFO] Tidak ditemukan versi sertifikat untuk subkategori ID %d", subcategoryID)
		issuedVersion = 0
	} else {
		log.Printf("[DEBUG] Versi sertifikat ditemukan: %d untuk subkategori ID %d", issuedVersion, subcategoryID)
	}

	subUpdateFields := map[string]interface{}{
		"user_subcategory_complete_themes_or_levels": userSub.UserSubcategoryCompleteThemesOrLevels,
	}
	if issuedVersion > 0 {
		if matchTheme == len(totalThemeIDs) && len(totalThemeIDs) > 0 {
			subUpdateFields["user_subcategory_grade_result"] = userSubcategoryFinalGradeResult
			subUpdateFields["user_subcategory_current_version"] = issuedVersion
			if err := issuedcertificateservice.CreateOrUpdateIssuedCertificate(db, userID, subcategoryID, issuedVersion); err != nil {
				log.Println("[WARNING] Gagal membuat/memperbarui sertifikat:", err)
			}
		} else if userSub.UserSubcategoryGradeResult > 0 && userSub.UserSubcategoryCurrentVersion < issuedVersion {
			log.Printf("[INFO] Update current_version karena sudah lulus dan ada versi baru: %d -> %d", userSub.UserSubcategoryCurrentVersion, issuedVersion)
			subUpdateFields["user_subcategory_current_version"] = issuedVersion
		}
	}

	if err := db.Model(&userSub).Updates(subUpdateFields).Error; err != nil {
		return err
	}

	return nil
}


func CheckAndUnsetExamStatus(db *gorm.DB, userID uuid.UUID, examID uint) error {
	// 🔍 Ambil unit ID dari exam
	var unitID uint
	err := db.Table("exams").
		Select("exam_unit_id").
		Where("exam_id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil exam_unit_id dari exam_id untuk reset status:", examID)
		return err
	}

	// 🔍 Cek apakah masih ada exam lain untuk unit tersebut
	var remainingExamCount int64
	err = db.Table("user_exams").
		Joins("JOIN exams ON exams.exam_id = user_exams.exam_id").
		Where("user_exams.user_id = ? AND exams.exam_unit_id = ?", userID, unitID).
		Count(&remainingExamCount).Error
	if err != nil {
		return err
	}

	// ❌ Jika tidak ada exam tersisa, reset nilai progress
	if remainingExamCount == 0 {
		log.Println("[INFO] Reset nilai exam dan result karena tidak ada user_exams tersisa, user_id:", userID, "unit_id:", unitID)

		// ✅ Reset nilai exam & hasil final di level unit
		if err := db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_unit_user_id = ? AND user_unit_unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"user_unit_grade_exam":   0,
				"user_unit_grade_result": 0,
				"updated_at":             time.Now(),
			}).Error; err != nil {
			return err
		}

		// 🔍 Ambil themes_or_level_id dari unit
		var themeID uint
		err = db.Table("units").
			Select("themes_or_level_id").
			Where("unit_id = ?", unitID).
			Scan(&themeID).Error
		if err != nil || themeID == 0 {
			log.Println("[ERROR] Gagal ambil themes_or_level_id dari unit:", unitID)
			return err
		}

		// 🔍 Cari record progress user untuk theme
		var userTheme userThemeModel.UserThemesOrLevelsModel
		err = db.Where("user_theme_user_id = ? AND user_theme_themes_or_level_id = ?", userID, themeID).
			First(&userTheme).Error
		if err != nil {
			log.Println("[WARNING] Tidak menemukan user_theme untuk reset complete_unit")
			return nil
		}

		if userTheme.UserThemeCompleteUnit != nil {
			unitKey := fmt.Sprintf("%d", unitID)
			delete(userTheme.UserThemeCompleteUnit, unitKey) // ✅ Hapus unit dari progress

			shouldResetThemeGrade := len(userTheme.UserThemeCompleteUnit) == 0

			updateThemeFields := map[string]interface{}{
				"complete_unit": userTheme.UserThemeCompleteUnit,
			}
			if shouldResetThemeGrade {
				updateThemeFields["grade_result"] = 0
			}

			if err := db.Model(&userTheme).Updates(updateThemeFields).Error; err != nil {
				log.Println("[ERROR] Gagal update user_theme saat reset:", err)
				return err
			}
		}
	}

	return nil
}
