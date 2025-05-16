package user_test_exam

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/users/test_exam/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTestExamSeed struct {
	UserID       uuid.UUID `json:"user_test_exam_user_id"`
	TestExamID   uint      `json:"user_test_exam_test_exam_id"`
	Grade        int       `json:"user_test_exam_percentage_grade"`
	TimeDuration int       `json:"user_test_exam_time_duration"`
}

func SeedUserTestExamsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file user_test_exams:", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []UserTestExamSeed
	if err := json.Unmarshal(content, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	var records []model.UserTestExam
	for _, item := range seeds {
		records = append(records, model.UserTestExam{
			UserTestExamUserID:          item.UserID,
			UserTestExamTestExamID:      item.TestExamID,
			UserTestExamPercentageGrade: item.Grade,
			UserTestExamTimeDuration:    item.TimeDuration,
		})
	}

	if len(records) > 0 {
		if err := db.Create(&records).Error; err != nil {
			log.Fatalf("❌ Gagal insert user_test_exams: %v", err)
		}
		log.Printf("✅ Berhasil insert %d data user_test_exam", len(records))
	} else {
		log.Println("ℹ️ Tidak ada data user_test_exam untuk diinsert.")
	}
}
