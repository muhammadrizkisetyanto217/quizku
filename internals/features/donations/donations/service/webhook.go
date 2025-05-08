package service

import (
	"fmt"
	"log"
	"time"

	"quizku/internals/features/donations/donations/model"

	"gorm.io/gorm"
)

func HandleDonationStatusWebhook(db *gorm.DB, body map[string]interface{}) error {
	orderID := body["order_id"].(string)
	transactionStatus := body["transaction_status"].(string)

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
		donation.Status = "paid"
		donation.PaidAt = &now
	case "expire":
		donation.Status = "expired"
	case "cancel":
		donation.Status = "canceled"
	default:
		log.Println("Status tidak diproses:", transactionStatus)
	}

	return db.Save(&donation).Error
}
