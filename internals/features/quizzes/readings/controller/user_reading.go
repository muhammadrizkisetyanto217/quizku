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
//
// Request Body:
//
//	{
//	  "reading_id": 12
//	}
//
// Response (201):
//
//	{
//	  "message": "User reading created successfully",
//	  "data": { ...UserReading }
//	}
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

	// ✅ Parse body dan validasi isi minimal (reading_id wajib)
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

	// ✅ Ambil data reading terkait untuk mendapatkan unit_id-nya
	var reading readingModel.ReadingModel
	if err := ctrl.DB.Select("id, unit_id").First(&reading, body.ReadingID).Error; err != nil {
		log.Println("[ERROR] Reading not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Reading not found"})
	}

	// ✅ Inisialisasi data yang akan disimpan
	input := UserReadingModel.UserReading{
		UserID:    userUUID,
		ReadingID: body.ReadingID,
		UnitID:    reading.UnitID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ✅ Hitung attempt ke-n (increment berdasarkan attempt terakhir)
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

	// ✅ Simpan entri ke database
	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user reading:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user reading"})
	}

	// ✅ Update progress pada user_unit (jika seluruh bacaan di unit selesai)
	if err := service.UpdateUserUnitFromReading(ctrl.DB, input.UserID, input.UnitID); err != nil {
		log.Println("[ERROR] Gagal update user_unit:", err)
	}

	// ✅ Tambahkan poin dari reading berdasarkan attempt
	if err := service.AddPointFromReading(ctrl.DB, input.UserID, input.ReadingID, input.Attempt); err != nil {
		log.Println("[ERROR] Gagal menambahkan poin:", err)
	}

	// ✅ Update aktivitas harian user (daily streak)
	if err := activityService.UpdateOrInsertDailyActivity(ctrl.DB, input.UserID); err != nil {
		log.Println("[ERROR] Gagal mencatat aktivitas harian:", err)
	}

	// ✅ Kembalikan response sukses
	log.Printf("[SUCCESS] UserReading created: user_id=%s, reading_id=%d, attempt=%d\n",
		input.UserID.String(), input.ReadingID, input.Attempt)

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user readings",
		})
	}

	return c.JSON(readings)
}
// GET /api/user-readings/user/:user_id
// 🔹 Ambil seluruh data pembacaan (reading) untuk satu user tertentu berdasarkan UUID.
// Digunakan untuk menampilkan riwayat bacaan user di dashboard atau profil.

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
