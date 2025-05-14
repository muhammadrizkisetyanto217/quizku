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

type QuizProgress struct {
	ID      uint `json:"id"`
	Attempt int  `json:"attempt"`
	Score   int  `json:"score"`
}

type SectionProgress struct {
	ID      uint `json:"id"`
	Score   int  `json:"score"`
	Attempt int  `json:"attempt"`
}

func UpdateUserSectionIfQuizCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	sectionID uint,
	quizID uint,
	attempt int,
	percentageGrade int,
) error {
	log.Println("[SERVICE] UpdateUserSectionIfQuizCompleted - userID:", userID, "sectionID:", sectionID, "quizID:", quizID, "attempt:", attempt, "score:", percentageGrade)

	// 1. Ambil semua quiz di section
	var allQuizzes []quizzesModel.QuizModel
	if err := db.Where("section_quizzes_id = ? AND deleted_at IS NULL", sectionID).Find(&allQuizzes).Error; err != nil {
		log.Println("[ERROR] Failed to fetch quizzes for section:", err)
		return err
	}
	var totalQuizIDs []uint
	for _, quiz := range allQuizzes {
		totalQuizIDs = append(totalQuizIDs, uint(quiz.ID))
	}

	// 2. Ambil atau buat user_section
	var userSection quizzesModel.UserSectionQuizzesModel
	err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).First(&userSection).Error

	var progressList []QuizProgress
	newProgress := QuizProgress{
		ID:      quizID,
		Attempt: attempt,
		Score:   percentageGrade,
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		progressList = []QuizProgress{newProgress}
		jsonData, _ := json.Marshal(progressList)
		userSection = quizzesModel.UserSectionQuizzesModel{
			UserID:           userID,
			SectionQuizzesID: sectionID,
			CompleteQuiz:     datatypes.JSON(jsonData),
			GradeResult:      0,
		}
		log.Println("[SERVICE] Creating new UserSectionQuizzesModel")
		return db.Create(&userSection).Error
	}

	// 3. Update progress existing
	if len(userSection.CompleteQuiz) > 0 {
		if err := json.Unmarshal(userSection.CompleteQuiz, &progressList); err != nil {
			log.Println("[ERROR] Failed to parse existing complete_quiz:", err)
			progressList = []QuizProgress{}
		}
	}

	updated := false
	for i, p := range progressList {
		if p.ID == quizID {
			if attempt > p.Attempt {
				progressList[i].Attempt = attempt
				updated = true
			}
			if percentageGrade > p.Score {
				progressList[i].Score = percentageGrade
				updated = true
			}
			break
		}
	}
	if !updated {
		progressList = append(progressList, newProgress)
	}

	// 4. Cek apakah semua quiz sudah dikerjakan
	completedQuizIDs := make(map[uint]bool)
	totalScore := 0
	for _, p := range progressList {
		completedQuizIDs[p.ID] = true
		totalScore += p.Score
	}

	// 5. Hitung nilai rata-rata jika semua quiz selesai
	if len(completedQuizIDs) == len(totalQuizIDs) && len(totalQuizIDs) > 0 {
		userSection.GradeResult = totalScore / len(progressList)
		log.Println("[SERVICE] Semua quiz selesai - GradeResult:", userSection.GradeResult)
	} else {
		userSection.GradeResult = 0
		log.Println("[SERVICE] Quiz belum selesai semua, GradeResult diset ke 0")
	}

	// 6. Simpan kembali
	newJSON, _ := json.Marshal(progressList)
	userSection.CompleteQuiz = datatypes.JSON(newJSON)

	log.Println("[SERVICE] Updating UserSectionQuizzesModel")
	return db.Save(&userSection).Error
}

func UpdateUserUnitIfSectionCompleted(
	db *gorm.DB,
	userID uuid.UUID,
	unitID uint,
	sectionID uint,
) error {
	log.Println("[SERVICE] UpdateUserUnitIfSectionCompleted - userID:", userID, "unitID:", unitID, "sectionID:", sectionID)

	// 1. Cek apakah section sudah memiliki progress
	var userSection quizzesModel.UserSectionQuizzesModel
	if err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sectionID).
		First(&userSection).Error; err != nil {
		log.Printf("[INFO] Section %d belum ada progress oleh user", sectionID)
		return nil
	}

	// 2. Ambil quiz ID hanya dari section yang sedang dicek
	var currentSection quizzesModel.SectionQuizzesModel
	if err := db.Preload("Quizzes").Where("id = ?", sectionID).First(&currentSection).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil section ID %d: %v", sectionID, err)
		return err
	}
	totalQuizIDs := map[int]struct{}{}
	for _, quiz := range currentSection.Quizzes {
		totalQuizIDs[int(quiz.ID)] = struct{}{}
	}

	// 3. Decode quiz yang telah diselesaikan dari user_section
	var completedQuizData []struct {
		ID      int `json:"id"`
		Score   int `json:"score"`
		Attempt int `json:"attempt"`
	}
	if err := json.Unmarshal(userSection.CompleteQuiz, &completedQuizData); err != nil {
		log.Printf("[ERROR] Gagal decode complete_quiz: %v", err)
		return err
	}
	completedIDs := map[int]bool{}
	for _, q := range completedQuizData {
		completedIDs[q.ID] = true
	}

	// 4. Cek apakah semua quiz dari section telah dikerjakan
	for id := range totalQuizIDs {
		if !completedIDs[id] {
			log.Printf("[INFO] Section %d belum lengkap, quiz ID %d belum dikerjakan", sectionID, id)
			return nil
		}
	}

	// 5. Ambil data user_unit
	var userUnit userUnitModel.UserUnitModel
	if err := db.Where("user_id = ? AND unit_id = ?", userID, unitID).First(&userUnit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil user_unit: %v", err)
		return err
	}

	// 6. Update complete_section_quizzes jika belum tercatat
	var completeSectionIDs []int64
	if len(userUnit.CompleteSectionQuizzes) > 0 {
		_ = json.Unmarshal(userUnit.CompleteSectionQuizzes, &completeSectionIDs)
	}
	alreadyExists := false
	for _, sid := range completeSectionIDs {
		if uint(sid) == sectionID {
			alreadyExists = true
			break
		}
	}
	if !alreadyExists {
		completeSectionIDs = append(completeSectionIDs, int64(sectionID))
		updatedJSON, err := json.Marshal(completeSectionIDs)
		if err != nil {
			log.Printf("[ERROR] Gagal encode complete_section_quizzes: %v", err)
			return err
		}
		userUnit.CompleteSectionQuizzes = updatedJSON
		userUnit.UpdatedAt = time.Now()
	}

	// 7. Hitung Grade jika semua section telah selesai
	var unit userUnitModel.UnitModel
	if err := db.Where("id = ?", unitID).First(&unit).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil unit: %v", err)
		return err
	}

	if len(unit.TotalSectionQuizzes) > 0 && len(completeSectionIDs) == len(unit.TotalSectionQuizzes) {
		total := 0
		count := 0
		for _, sid := range completeSectionIDs {
			var usq quizzesModel.UserSectionQuizzesModel
			if err := db.Where("user_id = ? AND section_quizzes_id = ?", userID, sid).
				First(&usq).Error; err != nil {
				log.Printf("[WARNING] Gagal ambil user_section_quizzes untuk section %d: %v", sid, err)
				continue
			}
			total += usq.GradeResult
			count++
		}
		if count > 0 {
			userUnit.GradeQuiz = total / count
			userUnit.GradeResult = (userUnit.GradeQuiz + userUnit.GradeExam + getGradeEvaluation(userUnit)) / 3
			userUnit.IsPassed = userUnit.GradeResult >= 70
			log.Printf("[SERVICE] Update GradeQuiz: %d, GradeResult: %d", userUnit.GradeQuiz, userUnit.GradeResult)
		}
	}

	return db.Save(&userUnit).Error
}

func getGradeEvaluation(u userUnitModel.UserUnitModel) int {
	type Eval struct {
		Attempt         int `json:"attempt"`
		GradeEvaluation int `json:"grade_evaluation"`
	}
	var e Eval
	if err := json.Unmarshal(u.AttemptEvaluation, &e); err != nil {
		return 0
	}
	return e.GradeEvaluation
}
