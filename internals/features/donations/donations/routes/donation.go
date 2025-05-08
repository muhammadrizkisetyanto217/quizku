package route

import (
	"quizku/internals/features/donations/donations/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DonationRoutes(api fiber.Router, db *gorm.DB) {
	donationCtrl := controller.NewDonationController(db)

	donationRoutes := api.Group("/donations")
	donationRoutes.Post("/", donationCtrl.CreateDonation)                   // Buat donasi + Snap token
	donationRoutes.Get("/", donationCtrl.GetAllDonations)                   // Semua donasi
	donationRoutes.Get("/user/:user_id", donationCtrl.GetDonationsByUserID) // Donasi per user

	// Webhook (dibiarkan tetap pakai app karena di luar protected routes)
	// Tapi kamu bisa pisahkan di main.go langsung jika ingin murni konsisten
	// atau buat route khusus public di main
}
