// 📁 controller/donation_controller.go
package controller

import (
	"fmt"
	"log"
	donationQuestionModel "quizku/internals/features/donations/donation_questions/model"
	"quizku/internals/features/donations/donations/model"
	donationModel "quizku/internals/features/donations/donations/model"
	donationService "quizku/internals/features/donations/donations/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DonationController struct {
	DB *gorm.DB
}

func NewDonationController(db *gorm.DB) *DonationController {
	return &DonationController{DB: db}
}

func (ctrl *DonationController) CreateDonation(c *fiber.Ctx) error {
	var body struct {
		UserID  string `json:"user_id"`
		Amount  int    `json:"amount"`
		Message string `json:"message"`
		Name    string `json:"name"`
		Email   string `json:"email"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userUUID, err := uuid.Parse(body.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id tidak valid"})
	}

	orderID := fmt.Sprintf("DONATION-%d", time.Now().UnixNano())
	donation := donationModel.Donation{
		UserID:  userUUID,
		Amount:  body.Amount,
		Message: body.Message,
		Status:  0,
		OrderID: orderID,
	}

	if err := ctrl.DB.Create(&donation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan donasi"})
	}

	token, err := donationService.GenerateSnapToken(donation, body.Name, body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token Midtrans"})
	}

	donation.PaymentToken = token
	ctrl.DB.Save(&donation)

	if donation.Status == donationModel.StatusPaid {
		soalCount := donation.Amount / 5000
		for i := 0; i < soalCount; i++ {
			entry := donationQuestionModel.DonationQuestionModel{
				DonationID:  donation.ID,
				QuestionID:  0,
				UserMessage: donation.Message,
			}
			if err := ctrl.DB.Create(&entry).Error; err != nil {
				log.Printf("[ERROR] Gagal buat slot kosong donation_question: %v", err)
			}
		}
		log.Printf("✅ Donasi paid langsung buat %d slot soal untuk donasi ID %d", soalCount, donation.ID)
	}

	return c.JSON(fiber.Map{
		"message":    "Donasi berhasil dibuat",
		"order_id":   donation.OrderID,
		"snap_token": token,
	})
}

func (ctrl *DonationController) HandleMidtransNotification(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid webhook"})
	}
	db := c.Locals("db").(*gorm.DB)
	if err := donationService.HandleDonationStatusWebhook(db, body); err != nil {
		log.Println("[ERROR] Webhook gagal:", err)
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

func (ctrl *DonationController) GetAllDonations(c *fiber.Ctx) error {
	var donations []model.Donation
	if err := ctrl.DB.Order("created_at desc").Find(&donations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data donasi"})
	}
	return c.JSON(donations)
}

func (ctrl *DonationController) GetDonationsByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id tidak valid"})
	}

	var donations []model.Donation
	if err := ctrl.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&donations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data donasi user"})
	}
	return c.JSON(donations)
}
