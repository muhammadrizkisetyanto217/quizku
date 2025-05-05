package controller

import (
	"log"
	UserReadingModel "quizku/internals/features/quizzes/readings/model"
	"quizku/internals/features/quizzes/readings/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	activityService "quizku/internals/features/progress/daily_activities/service"
)

type UserReadingController struct {
	DB *gorm.DB
}

func NewUserReadingController(db *gorm.DB) *UserReadingController {
	return &UserReadingController{DB: db}
}

// POST /user-readings
func (ctrl *UserReadingController) CreateUserReading(c *fiber.Ctx) error {
	// ✅ Ambil user_id dari JWT (string UUID)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// ✅ Parse body tanpa user_id
	type InputBody struct {
		ReadingID uint `json:"reading_id" validate:"required"`
		UnitID    uint `json:"unit_id" validate:"required"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// ✅ Validasi
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// ✅ Siapkan data
	input := UserReadingModel.UserReading{
		UserID:    userUUID,
		ReadingID: body.ReadingID,
		UnitID:    body.UnitID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ✅ Hitung attempt ke-n
	var latestAttempt int
	err = ctrl.DB.Table("user_readings").
		Select("COALESCE(MAX(attempt), 0)").
		Where("user_id = ? AND reading_id = ?", input.UserID, input.ReadingID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	input.Attempt = latestAttempt + 1

	// ✅ Simpan ke DB
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user reading"})
	}

	// ✅ Update progres & poin
	if err := service.UpdateUserUnitFromReading(ctrl.DB, input.UserID, input.UnitID); err != nil {
		log.Println("[ERROR] Gagal update user_unit:", err)
	}
	if err := service.AddPointFromReading(ctrl.DB, input.UserID, input.ReadingID, input.Attempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin:", err)
	}

	// ✅ Tambahkan aktivitas harian
	if err := activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserID); err != nil {
		log.Println("[ERROR] Gagal mencatat aktivitas harian:", err)
	}

	log.Printf("[SUCCESS] UserReading created: user_id=%s, reading_id=%d, attempt=%d\n",
		input.UserID.String(), input.ReadingID, input.Attempt)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User reading created successfully",
		"data":    input,
	})
}

// GET /user-readings
func (ctrl *UserReadingController) GetAllUserReading(c *fiber.Ctx) error {
	var readings []UserReadingModel.UserReading

	if err := ctrl.DB.Find(&readings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user readings",
		})
	}

	return c.JSON(readings)
}

// GET /api/user-readings/user/:user_id
func (ctrl *UserReadingController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var readings []UserReadingModel.UserReading
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&readings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil user_readings",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User readings fetched successfully",
		"data":    readings,
	})
}
