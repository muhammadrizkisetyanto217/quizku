package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	userUnitModel "quizku/internals/features/lessons/units/model"
	quizzesModel "quizku/internals/features/quizzes/quizzes/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SectionProgress struct {
	ID      uint `json:"id"`
	Score   int  `json:"score"`
	Attempt int  `json:"attempt"`
}
type UserQuizProgress struct {
	QuizID        uint `json:"quiz_id"`
	QuizAttempt   int  `json:"quiz_attempt"`
	QuizBestScore int  `json:"quiz_best_score"`
}

func UpdateUserSectionIfQuizCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	sectionQuizzesID uint,
	userQuizID uint,
	userQuizAttempt int,
	userQuizPercentageGrade int,
) error {
	log.Println("[SERVICE] UpdateUserSectionIfQuizCompleted - user:", userID, "section:", sectionQuizzesID, "quiz:", userQuizID)

	// 1. Ambil semua quiz_id di section
	var quizzesInSection []quizzesModel.QuizModel
	if err := db.
		Where("quiz_section_quizzes_id = ? AND deleted_at IS NULL", sectionQuizzesID).
		Find(&quizzesInSection).Error; err != nil {
		log.Println("[ERROR] Gagal mengambil daftar kuis dalam section:", err)
		return err
	}
	var allQuizIDsInSection []uint
	for _, quiz := range quizzesInSection {
		allQuizIDsInSection = append(allQuizIDsInSection, uint(quiz.QuizID))
	}

	// 2. Ambil atau buat user_section_quizzes
	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.
		Where("user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?", userID, sectionQuizzesID).
		First(&userSection).Error

	var progressList []UserQuizProgress
	newProgress := UserQuizProgress{
		QuizID:        userQuizID,
		QuizAttempt:   userQuizAttempt,
		QuizBestScore: userQuizPercentageGrade,
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 🔹 Belum ada progress section → buat baru
		progressList = []UserQuizProgress{newProgress}
		jsonData, _ := json.Marshal(progressList)

		userSection = quizzesModel.UserSectionQuizzesModel{
			UserSectionQuizzesUserID:           userID,
			UserSectionQuizzesSectionQuizzesID: sectionQuizzesID,
			UserSectionQuizzesCompleteQuiz:     datatypes.JSON(jsonData),
			UserSectionQuizzesGradeResult:      0,
		}

		log.Println("[SERVICE] Membuat UserSectionQuizzes baru")
		return db.Create(&userSection).Error
	}

	// 3. Jika sudah ada → decode JSON progress
	if len(userSection.UserSectionQuizzesCompleteQuiz) > 0 {
		if err := json.Unmarshal(userSection.UserSectionQuizzesCompleteQuiz, &progressList); err != nil {
			log.Println("[ERROR] Gagal decode progress section sebelumnya:", err)
			progressList = []UserQuizProgress{}
		}
	}

	// 4. Update atau tambahkan entry progress quiz yang sedang dikerjakan
	updated := false
	for i, p := range progressList {
		if p.QuizID == userQuizID {
			if userQuizAttempt > p.QuizAttempt {
				progressList[i].QuizAttempt = userQuizAttempt
				updated = true
			}
			if userQuizPercentageGrade > p.QuizBestScore {
				progressList[i].QuizBestScore = userQuizPercentageGrade
				updated = true
			}
			break
		}
	}
	if !updated {
		progressList = append(progressList, newProgress)
	}

	// 5. Cek apakah semua kuis di section sudah dikerjakan
	completedQuizIDs := make(map[uint]bool)
	totalScore := 0
	for _, p := range progressList {
		completedQuizIDs[p.QuizID] = true
		totalScore += p.QuizBestScore
	}

	if len(completedQuizIDs) == len(allQuizIDsInSection) && len(allQuizIDsInSection) > 0 {
		userSection.UserSectionQuizzesGradeResult = totalScore / len(progressList)
		log.Println("[SERVICE] Semua kuis dalam section telah dikerjakan. GradeResult =", userSection.UserSectionQuizzesGradeResult)
	} else {
		userSection.UserSectionQuizzesGradeResult = 0
		log.Println("[SERVICE] Kuis belum lengkap. GradeResult direset ke 0")
	}

	// 6. Simpan progres terbaru ke database
	newJSON, _ := json.Marshal(progressList)
	userSection.UserSectionQuizzesCompleteQuiz = datatypes.JSON(newJSON)

	log.Println("[SERVICE] Menyimpan update UserSectionQuizzes")
	return db.Save(&userSection).Error
}

func UpdateUserUnitIfSectionCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	unitID uint,
	completedSectionID uint,
) error {
	log.Printf("[SERVICE] UpdateUserUnitIfSectionCompleted - userID: %s, unitID: %d, completedSectionID: %d",
		userID.String(), unitID, completedSectionID)

	// 1. Cek apakah section memiliki progres di user_section_quizzes
	var userSection quizzesModel.UserSectionQuizzesModel
	if err := db.Where(
		"user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?",
		userID, completedSectionID,
	).First(&userSection).Error; err != nil {
		log.Printf("[INFO] Section %d belum ada progress oleh user", completedSectionID)
		return nil
	}

	// 2. Ambil semua quiz dalam section tersebut
	var section quizzesModel.SectionQuizzesModel
	if err := db.Preload("Quizzes").
		Where("section_quizzes_id = ?", completedSectionID).
		First(&section).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil section ID %d: %v", completedSectionID, err)
		return err
	}
	totalQuizIDs := map[int]struct{}{}
	for _, quiz := range section.Quizzes {
		totalQuizIDs[int(quiz.QuizID)] = struct{}{}
	}

	// 3. Decode quiz yang telah diselesaikan dari user_section
	type UserQuizProgress struct {
		UserQuizQuizID    int `json:"quiz_id"`
		UserQuizAttempt   int `json:"quiz_attempt"`
		UserQuizBestScore int `json:"quiz_best_score"`
	}

	var completedQuizData []UserQuizProgress
	if err := json.Unmarshal(userSection.UserSectionQuizzesCompleteQuiz, &completedQuizData); err != nil {
		log.Printf("[ERROR] Gagal decode complete_quiz: %v", err)
		return err
	}
	completedQuizIDs := map[int]bool{}
	for _, quiz := range completedQuizData {
		completedQuizIDs[quiz.UserQuizQuizID] = true
	}

	// 4. Cek apakah semua quiz dari section telah dikerjakan
	for id := range totalQuizIDs {
		if !completedQuizIDs[id] {
			log.Printf("[INFO] Section %d belum lengkap, quiz ID %d belum dikerjakan", completedSectionID, id)
			return nil
		}
	}

	// 5. Ambil data user_unit
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_unit_user_id = ? AND user_unit_unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil user_unit: %v", err)
		return err
	}

	// 6. Update complete_section_quizzes jika belum tercatat
	var completedSectionIDs []int64
	if len(userUnit.UserUnitCompleteSectionQuizzes) > 0 {
		_ = json.Unmarshal(userUnit.UserUnitCompleteSectionQuizzes, &completedSectionIDs)
	}
	alreadyIncluded := false
	for _, sid := range completedSectionIDs {
		if uint(sid) == completedSectionID {
			alreadyIncluded = true
			break
		}
	}
	if !alreadyIncluded {
		completedSectionIDs = append(completedSectionIDs, int64(completedSectionID))
		if encoded, err := json.Marshal(completedSectionIDs); err == nil {
			userUnit.UserUnitCompleteSectionQuizzes = encoded
			userUnit.UpdatedAt = time.Now()
		}
	}

	// 7. Hitung nilai rata-rata kuis (GradeQuiz) dan total kelulusan (GradeResult) jika semua section sudah selesai
	var unit userUnitModel.UnitModel
	if err := db.Where("unit_id = ?", unitID).First(&unit).Error; err != nil {
		log.Printf("[ERROR] Gagal mengambil data unit (unit_id=%d): %v", unitID, err)
		return err
	}

	if len(unit.UnitTotalSectionQuizzes) > 0 && len(completedSectionIDs) == len(unit.UnitTotalSectionQuizzes) {
		totalQuizScore := 0
		sectionCount := 0

		for _, sectionID := range completedSectionIDs {
			var userSection quizzesModel.UserSectionQuizzesModel
			if err := db.Where(
				"user_section_quizzes_user_id = ? AND user_section_quizzes_section_quizzes_id = ?",
				userID, sectionID,
			).First(&userSection).Error; err != nil {
				log.Printf("[WARNING] Gagal mengambil progres section (section_id=%d): %v", sectionID, err)
				continue
			}
			totalQuizScore += userSection.UserSectionQuizzesGradeResult
			sectionCount++
		}

		if sectionCount > 0 {
			userUnit.UserUnitGradeQuiz = totalQuizScore / sectionCount
			userUnit.UserUnitGradeResult = (userUnit.UserUnitGradeQuiz + userUnit.UserUnitGradeExam + getGradeEvaluation(userUnit)) / 3
			userUnit.UserUnitIsPassed = userUnit.UserUnitGradeResult >= 70

			log.Printf("[SERVICE] ✅ Semua section selesai. GradeQuiz: %d, GradeResult: %d, IsPassed: %v",
				userUnit.UserUnitGradeQuiz, userUnit.UserUnitGradeResult, userUnit.UserUnitIsPassed)
		}
	}

	return db.Save(&userUnit).Error
}


func getGradeEvaluation(userUnit userUnitModel.UserUnitModel) int {
	type EvaluationAttempt struct {
		EvaluationAttemptCount int `json:"attempt"`
		EvaluationScore        int `json:"grade_evaluation"`
	}

	var evalData EvaluationAttempt
	if err := json.Unmarshal(userUnit.UserUnitAttemptEvaluation, &evalData); err != nil {
		log.Printf("[ERROR] Gagal mengurai JSON AttemptEvaluation: %v", err)
		return 0
	}
	return evalData.EvaluationScore
}

