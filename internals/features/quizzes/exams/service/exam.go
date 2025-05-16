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

	issuedcertificateservice "quizku/internals/features/certificates/issued_certificates/service"
	userSubcategoryModel "quizku/internals/features/lessons/subcategories/model"
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

	// Ambil unit_id berdasarkan exam
	var unitID uint
	if err := db.Table("exams").Select("unit_id").Where("id = ?", examID).Scan(&unitID).Error; err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id:", examID)
		return err
	}

	// Ambil record user_unit untuk user dan unit tersebut
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		return err
	}

	// Hitung bonus berdasarkan aktivitas lain
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

	// Cek apakah semua section_quizzes pada unit sudah diselesaikan
	var totalSections, completedSections int64
	_ = db.Table("section_quizzes").Where("unit_id = ?", unitID).Count(&totalSections).Error
	_ = db.Table("user_section_quizzes").
		Joins("JOIN section_quizzes ON user_section_quizzes.section_quizzes_id = section_quizzes.id").
		Where("user_section_quizzes.user_id = ? AND section_quizzes.unit_id = ?", userID, unitID).
		Count(&completedSections).Error
	if totalSections > 0 && totalSections == completedSections {
		activityBonus += 30
	}

	//* ✅ Hitung grade_result untuk LEVEL UNIT (user_unit.grade_result)
	gradeResult := grade / 2
	if activityBonus > 0 {
		gradeResult = activityBonus + (grade / 2)
	}

	// Update nilai di tabel user_unit
	updates := map[string]interface{}{
		"grade_result": gradeResult, // ✅ GRADE RESULT UNTUK UNIT
		"is_passed":    gradeResult > 65,
		"updated_at":   time.Now(),
	}
	if grade > userUnit.GradeExam {
		updates["grade_exam"] = grade
	}
	if err := db.Model(&userUnit).Updates(updates).Error; err != nil {
		return err
	}

	// Jika belum lulus, stop di sini
	if gradeResult <= 65 {
		return nil
	}

	// Ambil theme ID dari unit
	var themesID uint
	if err := db.Table("units").Select("themes_or_level_id").Where("id = ?", unitID).Scan(&themesID).Error; err != nil || themesID == 0 {
		return err
	}

	// Update progress untuk user_themes_or_levels
	var userTheme userThemeModel.UserThemesOrLevelsModel
	if err := db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).First(&userTheme).Error; err != nil {
		return err
	}
	if userTheme.CompleteUnit == nil {
		userTheme.CompleteUnit = datatypes.JSONMap{}
	}
	userTheme.CompleteUnit[fmt.Sprintf("%d", unitID)] = fmt.Sprintf("%d", gradeResult)

	// Hitung nilai rata-rata dari seluruh unit dalam theme tersebut
	var expectedUnitIDs []int64
	if err := db.Table("units").Where("themes_or_level_id = ?", themesID).Pluck("id", &expectedUnitIDs).Error; err != nil {
		return err
	}
	matchCount := 0
	total := 0
	for _, id := range expectedUnitIDs {
		if val, ok := userTheme.CompleteUnit[fmt.Sprintf("%d", id)]; ok {
			matchCount++
			if g, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
				total += g
			}
		}
	}
	avg := 0
	if len(expectedUnitIDs) > 0 {
		avg = total / len(expectedUnitIDs)
	}

	//* ✅ grade_result untuk LEVEL THEME (user_themes_or_levels.grade_result)
	updateFields := map[string]interface{}{
		"complete_unit": userTheme.CompleteUnit,
	}
	if matchCount == len(expectedUnitIDs) && len(expectedUnitIDs) > 0 {
		updateFields["grade_result"] = avg // ✅ GRADE RESULT UNTUK THEME
	}
	if err := db.Model(&userTheme).Updates(updateFields).Error; err != nil {
		return err
	}

	// Ambil subcategory dari theme
	var subcategoryID int
	if err := db.Table("themes_or_levels").Select("subcategories_id").Where("id = ?", themesID).Scan(&subcategoryID).Error; err != nil {
		return err
	}

	// Update progress untuk user_subcategory
	var userSub userSubcategoryModel.UserSubcategoryModel
	if err := db.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).First(&userSub).Error; err != nil {
		return err
	}
	if userSub.CompleteThemesOrLevels == nil {
		userSub.CompleteThemesOrLevels = datatypes.JSONMap{}
	}
	userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", themesID)] = fmt.Sprintf("%d", avg)

	// Ambil total themes dari subcategory (array themes_id)
	var raw string
	if err := db.Table("subcategories").Select("total_themes_or_levels").Where("id = ?", subcategoryID).Scan(&raw).Error; err != nil {
		return err
	}
	var totalThemeIDs pq.Int64Array
	if err := totalThemeIDs.Scan(raw); err != nil {
		log.Println("[ERROR] Gagal parsing total_themes_or_levels:", err)
		return err
	}

	// Hitung nilai rata-rata dari semua themes
	matchTheme := 0
	totalSub := 0
	for _, id := range totalThemeIDs {
		if val, ok := userSub.CompleteThemesOrLevels[fmt.Sprintf("%d", id)]; ok {
			matchTheme++
			if g, err := strconv.Atoi(fmt.Sprintf("%v", val)); err == nil {
				totalSub += g
			}
		}
	}
	avgSub := 0
	if len(totalThemeIDs) > 0 {
		avgSub = totalSub / len(totalThemeIDs)
	}

	// Ambil versi sertifikat terbaru
	var issuedVersion int
	row := db.Table("certificate_versions").
		Where("subcategory_id = ?", subcategoryID).
		Select("version_number").
		Order("version_number DESC").
		Limit(1).
		Row()
	if err := row.Scan(&issuedVersion); err != nil {
		log.Printf("[INFO] Tidak ditemukan versi sertifikat untuk subkategori ID %d, current_version tidak akan diupdate", subcategoryID)
		issuedVersion = 0
	} else {
		log.Printf("[DEBUG] Versi sertifikat ditemukan: %d untuk subkategori ID %d", issuedVersion, subcategoryID)
	}

	//* ✅ grade_result untuk LEVEL SUBCATEGORY (user_subcategory.grade_result)
	updateSubFields := map[string]interface{}{
		"complete_themes_or_levels": userSub.CompleteThemesOrLevels,
	}
	if issuedVersion > 0 {
		if matchTheme == len(totalThemeIDs) && len(totalThemeIDs) > 0 {
			updateSubFields["grade_result"] = avgSub // ✅ GRADE RESULT UNTUK SUBCATEGORY
			updateSubFields["current_version"] = issuedVersion
			if err := issuedcertificateservice.CreateOrUpdateIssuedCertificate(db, userID, subcategoryID, issuedVersion); err != nil {
				log.Println("[WARNING] Gagal membuat/memperbarui sertifikat:", err)
			}
		} else if userSub.GradeResult > 0 && userSub.CurrentVersion < issuedVersion {
			log.Printf("[INFO] Update current_version karena sudah lulus dan ada versi baru: %d -> %d", userSub.CurrentVersion, issuedVersion)
			updateSubFields["current_version"] = issuedVersion
		}
	}

	if err := db.Model(&userSub).Updates(updateSubFields).Error; err != nil {
		return err
	}

	return nil
}

func CheckAndUnsetExamStatus(db *gorm.DB, userID uuid.UUID, examID uint) error {
	// 🔍 Ambil unit_id dari exam
	var unitID uint
	err := db.Table("exams").
		Select("unit_id").
		Where("id = ?", examID).
		Scan(&unitID).Error
	if err != nil || unitID == 0 {
		log.Println("[ERROR] Gagal ambil unit_id dari exam_id untuk reset status:", examID)
		return err
	}

	// 🔍 Cek apakah masih ada user_exams lain untuk unit ini
	var count int64
	err = db.Table("user_exams").
		Joins("JOIN exams ON exams.id = user_exams.exam_id").
		Where("user_exams.user_id = ? AND exams.unit_id = ?", userID, unitID).
		Count(&count).Error
	if err != nil {
		return err
	}

	// ❌ Jika tidak ada exam tersisa untuk unit tersebut, reset nilai progress
	if count == 0 {
		log.Println("[INFO] Reset nilai exam dan result karena tidak ada user_exams tersisa, user_id:", userID, "unit_id:", unitID)

		// ✅ Reset nilai exam dan hasil akhir pada level UNIT
		// grade_exam    → nilai mentah dari ujian terakhir
		// grade_result  → nilai final gabungan dengan aktivitas, dihitung oleh UpdateUserUnitFromExam
		if err := db.Model(&userUnitModel.UserUnitModel{}).
			Where("user_id = ? AND unit_id = ?", userID, unitID).
			Updates(map[string]interface{}{
				"grade_exam":   0, // ✅ Reset nilai ujian
				"grade_result": 0, // ✅ Reset hasil final unit (kombinasi exam + aktivitas)
				"updated_at":   time.Now(),
			}).Error; err != nil {
			return err
		}

		// 🔍 Ambil themes_or_level_id dari unit
		var themesID uint
		err = db.Table("units").
			Select("themes_or_level_id").
			Where("id = ?", unitID).
			Scan(&themesID).Error
		if err != nil || themesID == 0 {
			log.Println("[ERROR] Gagal ambil themes_or_level_id dari unit:", unitID)
			return err
		}

		// 🔍 Cari record progress user untuk theme tersebut
		var userTheme userThemeModel.UserThemesOrLevelsModel
		err = db.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
			First(&userTheme).Error
		if err != nil {
			log.Println("[WARNING] Tidak menemukan user_theme untuk reset complete_unit")
			return nil // Tidak error fatal, bisa dilanjut
		}

		if userTheme.CompleteUnit != nil {
			unitKey := fmt.Sprintf("%d", unitID)
			delete(userTheme.CompleteUnit, unitKey) // ✅ Hapus unit dari progress theme

			// ❓ Apakah semua unit pada theme sudah dihapus?
			shouldResetGrade := len(userTheme.CompleteUnit) == 0

			updateFields := map[string]interface{}{
				"complete_unit": userTheme.CompleteUnit,
			}
			if shouldResetGrade {
				// ✅ Reset nilai theme jika tidak ada unit tersisa
				// grade_result di theme adalah hasil rata-rata semua unit di theme
				// Jika tidak ada satupun unit, maka hasilnya dianggap belum ada
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
