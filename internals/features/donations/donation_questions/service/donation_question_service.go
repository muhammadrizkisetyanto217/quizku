package service

// import (
// 	"fmt"
// 	"log"

// 	donationModel "quizku/internals/features/donations/donations/model"
// 	donationQuestionModel "quizku/internals/features/donations/donation_questions/model"
// 	questionModel "quizku/internals/features/quizzes/questions/model"

// 	"gorm.io/gorm"
// )

// func CreateDonationQuestionsFromDonation(donation *donationModel.Donation, db *gorm.DB) error {
// 	if donation.Status != donationModel.StatusPaid {
// 		return fmt.Errorf("donasi belum paid, tidak buat soal")
// 	}

// 	soalCount := donation.Amount / 5000
// 	if soalCount <= 0 {
// 		return fmt.Errorf("jumlah soal 0, tidak buat soal")
// 	}

// 	var questions []questionModel.QuestionModel
// 	if err := db.
// 		Where("status = ?", "active").
// 		Order("RANDOM()").
// 		Limit(soalCount).
// 		Find(&questions).Error; err != nil {
// 		return fmt.Errorf("gagal ambil soal: %v", err)
// 	}

// 	for _, q := range questions {
// 		entry := donationQuestionModel.DonationQuestionModel{
// 			DonationID:  donation.ID,
// 			QuestionID:  q.ID,
// 			UserMessage: donation.Message,
// 		}
// 		if err := db.Create(&entry).Error; err != nil {
// 			log.Printf("[ERROR] Gagal buat donation_question untuk soal %d: %v", q.ID, err)
// 		}
// 	}

// 	log.Printf("✅ Berhasil buat %d soal dari donasi ID %d", len(questions), donation.ID)
// 	return nil
// }
