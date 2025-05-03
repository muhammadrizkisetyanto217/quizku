package user

import (
	"encoding/json"
	"log"
	"os"
	"quizku/internals/features/users/user/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersProfileSeed struct {
	UserID       uuid.UUID        `json:"user_id"`
	DonationName string           `json:"donation_name"`
	FullName     string           `json:"full_name"`
	Gender       *models.Gender   `json:"gender"`
	PhoneNumber  string           `json:"phone_number"`
	Bio          string           `json:"bio"`
	Location     string           `json:"location"`
	Occupation   string           `json:"occupation"`
}

func SeedUsersProfileFromJSON(db *gorm.DB, filePath string) {
	log.Println("📥 Membaca file:", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Gagal membaca file JSON: %v", err)
	}

	var profiles []UsersProfileSeed
	if err := json.Unmarshal(file, &profiles); err != nil {
		log.Fatalf("❌ Gagal decode JSON: %v", err)
	}

	for _, p := range profiles {
		var existing models.UsersProfileModel
		err := db.Where("user_id = ?", p.UserID).First(&existing).Error
		if err == nil {
			log.Printf("ℹ️ Profil user dengan ID '%s' sudah ada, dilewati.", p.UserID)
			continue
		}

		newProfile := models.UsersProfileModel{
			UserID:       p.UserID,
			DonationName: p.DonationName,
			FullName:     p.FullName,
			Gender:       p.Gender,
			PhoneNumber:  p.PhoneNumber,
			Bio:          p.Bio,
			Location:     p.Location,
			Occupation:   p.Occupation,
		}

		if err := db.Create(&newProfile).Error; err != nil {
			log.Printf("❌ Gagal insert profil user ID %s: %v", p.UserID, err)
		} else {
			log.Printf("✅ Berhasil insert profil user ID %s", p.UserID)
		}
	}
}
