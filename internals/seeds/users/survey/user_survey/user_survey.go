package survey

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/users/survey/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSurveySeed struct {
	UserID           uuid.UUID `json:"user_id"`
	SurveyQuestionID int       `json:"survey_question_id"`
	UserAnswer       string    `json:"user_answer"`
}

func SeedUserSurveysFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file user survey:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var seeds []UserSurveySeed
	if err := json.Unmarshal(file, &seeds); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	var userSurveys []model.UserSurvey
	for _, s := range seeds {
		userSurveys = append(userSurveys, model.UserSurvey{
			UserSurveyUserID:           s.UserID,
			UserSurveyQuestionID: s.SurveyQuestionID,
			UserSurveyAnswer:       s.UserAnswer,
		})
	}

	if len(userSurveys) > 0 {
		if err := db.Create(&userSurveys).Error; err != nil {
			log.Fatalf("❌ Gagal insert user_surveys: %v", err)
		}
		log.Printf("✅ Berhasil insert %d user survey", len(userSurveys))
	} else {
		log.Println("ℹ️ Tidak ada data user survey untuk diinsert.")
	}
}
