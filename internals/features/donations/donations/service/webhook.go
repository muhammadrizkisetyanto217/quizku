package service

import (
	"fmt"
	"log"
	"time"

	donationService "quizku/internals/features/donations/donation_questions/service"
	donationModel "quizku/internals/features/donations/donations/model"

	"gorm.io/gorm"
)

// HandleDonationStatusWebhook dipanggil saat menerima notifikasi dari Midtrans
func HandleDonationStatusWebhook(db *gorm.DB, body map[string]interface{}) error {

	log.Println("🔥🔥🔥 WEBHOOK MASUK 🔥🔥🔥")
	log.Printf("📩 Payload dari Midtrans: %+v", body)

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

	// Ambil donasi berdasarkan order ID
	var donation donationModel.Donation
	if err := db.Where("order_id = ?", orderID).First(&donation).Error; err != nil {
		log.Println("[ERROR] Order tidak ditemukan:", err)
		return fmt.Errorf("donation with order_id %s not found", orderID)
	}

	// Proses status
	switch transactionStatus {
	case "capture", "settlement":
		now := time.Now()
		donation.Status = donationModel.StatusPaid
		donation.PaidAt = &now

		// ✅ Simpan status Paid ke DB
		if err := db.Save(&donation).Error; err != nil {
			log.Println("[ERROR] Gagal menyimpan status donasi:", err)
			return fmt.Errorf("gagal menyimpan status donasi")
		}

		// ✅ Setelah tersimpan, baru buat soal
		if err := donationService.CreateDonationQuestionsFromDonation(&donation, db); err != nil {
			log.Printf("[ERROR] Gagal generate slot soal: %v", err)
			// tidak return error agar Midtrans tetap dapat response 200
		}

	case "expire":
		donation.Status = donationModel.StatusExpired
		if err := db.Save(&donation).Error; err != nil {
			log.Println("[ERROR] Gagal menyimpan status expired:", err)
			return fmt.Errorf("gagal simpan status expire")
		}

	case "cancel":
		donation.Status = donationModel.StatusCanceled
		if err := db.Save(&donation).Error; err != nil {
			log.Println("[ERROR] Gagal menyimpan status cancel:", err)
			return fmt.Errorf("gagal simpan status cancel")
		}

	default:
		log.Println("📌 Status transaksi tidak diproses:", transactionStatus)
	}

	log.Println("✅ Status donasi berhasil diperbarui:", donation.Status)
	return nil
}
