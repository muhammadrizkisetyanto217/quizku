package service

import (
	"fmt"
	"log"
	"quizku/internals/features/donations/donations/model"
	donationQuestionModel "quizku/internals/features/donations/donation_questions/model"
	"time"

	"gorm.io/gorm"
)

// HandleDonationStatusWebhook dipanggil saat menerima notifikasi dari Midtrans

func HandleDonationStatusWebhook(db *gorm.DB, body map[string]interface{}) error {
	orderID, ok1 := body["order_id"].(string)
	status, ok2 := body["transaction_status"].(string)

	if !ok1 || !ok2 {
		log.Println("[ERROR] Payload tidak lengkap:", body)
		return fmt.Errorf("invalid payload")
	}

	log.Println("📄 Order ID:", orderID)
	log.Println("📌 Transaction Status:", status)

	var donation model.Donation
	if err := db.Where("order_id = ?", orderID).First(&donation).Error; err != nil {
		log.Println("[ERROR] Donasi tidak ditemukan:", err)
		return fmt.Errorf("donation with order_id %s not found", orderID)
	}

	// Update status donasi
	switch status {
	case "capture", "settlement":
		now := time.Now()
		donation.Status = "paid"
		donation.PaidAt = &now

		// Tambahkan entri ke donation_questions
		totalSoal := int(donation.Amount) / 5000
		for i := 0; i < totalSoal; i++ {
			soal := donationQuestionModel.DonationQuestionModel{
				DonationID:  donation.ID,
				QuestionID:  0,                // default dulu, nanti bisa diisi real ID soal kalau ada
				UserMessage: donation.Message, // optional, bisa juga kosong
			}
			if err := db.Create(&soal).Error; err != nil {
				log.Println("[ERROR] Gagal buat donation_question:", err)
				// Lanjut ke soal berikutnya, jangan return
			}
		}

	case "expire":
		donation.Status = "expired"
	case "cancel":
		donation.Status = "canceled"
	default:
		log.Println("Status tidak diproses:", status)
	}

	return db.Save(&donation).Error
}
