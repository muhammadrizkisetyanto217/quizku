package controller

import (
	"log"
	UserReadingModel "quizku/internals/features/quizzes/readings/model"
	readingModel "quizku/internals/features/quizzes/readings/model"
	"quizku/internals/features/quizzes/readings/service"
	"time"

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
// Fungsi ini menangani pencatatan aktivitas membaca oleh user.
// Endpoint ini memerlukan autentikasi JWT dan akan:
// - Menyimpan data pembacaan (user_id, reading_id, unit_id, attempt, timestamp)
// - Mengupdate progres user_unit terkait
// - Menambahkan poin sesuai attempt dan reading
// - Mencatat aktivitas harian (daily streak)

func (ctrl *UserReadingController) CreateUserReading(c *fiber.Ctx) error {
	// ✅ Ambil user_id dari JWT token
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("[ERROR] Invalid UUID format:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// ✅ Parse body dan validasi isi minimal
	type InputBody struct {
		ReadingID uint `json:"reading_id"`
	}
	var body InputBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("[ERROR] Failed to parse body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if body.ReadingID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "reading_id is required"})
	}

	// ✅ Ambil data reading terkait menggunakan kolom semantik
	var reading readingModel.ReadingModel
	if err := ctrl.DB.
		Select("reading_id, reading_unit_id").
		First(&reading, "reading_id = ?", body.ReadingID).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Reading not found"})
	}

	// ✅ Cek dan hitung attempt sebelumnya
	var latestAttempt int
	err = ctrl.DB.Table("user_readings").
		Select("COALESCE(MAX(user_reading_attempt), 0)").
		Where("user_reading_user_id = ? AND user_reading_reading_id = ?", userUUID, body.ReadingID).
		Scan(&latestAttempt).Error
	if err != nil {
		log.Println("[ERROR] Failed to count latest attempt:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Inisialisasi user reading
	input := UserReadingModel.UserReading{
		UserReadingUserID:    userUUID,
		UserReadingReadingID: body.ReadingID,
		UserReadingUnitID:    reading.ReadingUnitID,
		UserReadingAttempt:   latestAttempt + 1,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// ✅ Simpan ke DB
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user reading"})
	}

	// ✅ Update progress user_unit
	if err := service.UpdateUserUnitFromReading(ctrl.DB, input.UserReadingUserID, input.UserReadingUnitID); err != nil {
		log.Println("[ERROR] Gagal update user_unit:", err)
	}

	// ✅ Tambahkan poin
	if err := service.AddPointFromReading(ctrl.DB, input.UserReadingUserID, input.UserReadingReadingID, input.UserReadingAttempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin:", err)
	}

	// ✅ Update aktivitas harian
	if err := activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserReadingUserID); err != nil {
		log.Println("[ERROR] Gagal mencatat aktivitas harian:", err)
	}

	log.Printf("[SUCCESS] UserReading created: user_id=%s, reading_id=%d, attempt=%d\n",
		input.UserReadingUserID.String(), input.UserReadingReadingID, input.UserReadingAttempt)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User reading created successfully",
		"data":    input,
	})
}

// GET /user-readings
// 🔹 Ambil semua data pembacaan user dari tabel user_readings (tidak difilter).
// ⚠️ Umumnya hanya digunakan untuk keperluan admin atau debug.
func (ctrl *UserReadingController) GetAllUserReading(c *fiber.Ctx) error {
	var readings []UserReadingModel.UserReading

	if err := ctrl.DB.Find(&readings).Error; err != nil {
		log.Println("[ERROR] Failed to fetch all user_readings:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user readings",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d user_readings\n", len(readings))
	return c.JSON(fiber.Map{
		"message": "All user readings fetched successfully",
		"total":   len(readings),
		"data":    readings,
	})
}

// GET /api/user-readings/user/:user_id
// 🔹 Ambil seluruh data pembacaan (reading) untuk satu user tertentu berdasarkan UUID.
// Digunakan untuk menampilkan riwayat bacaan user di dashboard atau profil.

func (ctrl *UserReadingController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		log.Println("[ERROR] Invalid user_id:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var readings []UserReadingModel.UserReading
	if err := ctrl.DB.
		Where("user_reading_user_id = ?", userUUID).
		Find(&readings).Error; err != nil {
		log.Println("[ERROR] Failed to fetch user_readings for user_id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data bacaan user",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d user_readings for user_id=%s\n", len(readings), userUUID.String())
	return c.JSON(fiber.Map{
		"message": "User readings fetched successfully",
		"total":   len(readings),
		"data":    readings,
	})
}
