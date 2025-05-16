package testexam

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/users/test_exam/model"

	"gorm.io/gorm"
)

type TestExamSeed struct {
	TestExamName   string `json:"test_exam_name"`
	TestExamStatus string `json:"test_exam_status"`
}

func SeedTestExamsFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file test exams:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []TestExamSeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	var testExams []model.TestExam
	for _, s := range seeds {
		testExams = append(testExams, model.TestExam{
			TestExamName:   s.TestExamName,
			TestExamStatus: s.TestExamStatus,
		})
	}

	if len(testExams) > 0 {
		if err := db.Create(&testExams).Error; err != nil {
			log.Fatalf("❌ Gagal insert test_exams: %v", err)
		}
		log.Printf("✅ Berhasil insert %d test exam", len(testExams))
	} else {
		log.Println("ℹ️ Tidak ada data test exam untuk diinsert.")
	}
}
