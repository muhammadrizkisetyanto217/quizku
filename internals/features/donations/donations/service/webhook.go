package service

import (
	"fmt"
	"log"
	"time"

	"quizku/internals/features/donations/donations/model"

	"gorm.io/gorm"
)

func HandleDonationStatusWebhook(db *gorm.DB, body map[string]interface{}) error {
	orderID, ok := body["order_id"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing order_id in webhook body")
	}

	transactionStatus, ok := body["transaction_status"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing transaction_status in webhook body")
	}

	log.Println("📄 Order ID:", orderID)
	log.Println("📌 Transaction Status:", transactionStatus)

	var donation model.Donation
	if err := db.Where("order_id = ?", orderID).First(&donation).Error; err != nil {
		log.Println("[ERROR] Order tidak ditemukan:", err)
		return fmt.Errorf("donation with order_id %s not found", orderID)
	}

	switch transactionStatus {
	case "capture", "settlement":
		now := time.Now()
		donation.Status = model.StatusPaid
		donation.PaidAt = &now
	case "expire":
		donation.Status = model.StatusExpired
	case "cancel":
		donation.Status = model.StatusCanceled
	default:
		log.Println("Status tidak diproses:", transactionStatus)
	}

	if err := db.Save(&donation).Error; err != nil {
		log.Println("[ERROR] Gagal menyimpan status donasi:", err)
		return fmt.Errorf("gagal menyimpan status donasi")
	}

	log.Println("✅ Status donasi berhasil diperbarui")
	return nil
}
