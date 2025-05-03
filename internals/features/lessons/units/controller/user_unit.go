package controller

import (
	"log"
	"net/http"
	"strconv"

	themesOrLevelsModel "quizku/internals/features/lessons/themes_or_levels/model"
	userModel "quizku/internals/features/lessons/units/model"
	userSectionQuizzesModel "quizku/internals/features/quizzes/quizzes/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUnitController struct {
	DB *gorm.DB
}

func NewUserUnitController(db *gorm.DB) *UserUnitController {
	return &UserUnitController{DB: db}
}

// GET /api/user-units/:user_id
func (ctrl *UserUnitController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []userModel.UserUnitModel
	if err := ctrl.DB.
		Preload("SectionProgress").
		Where("user_id = ?", userID).
		Find(&data).Error; err != nil {
		log.Println("[ERROR] Gagal ambil data user_unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevelsAndUserID(c *fiber.Ctx) error {
	// 🔐 Ambil user_id dari token JWT
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user_id tidak ditemukan dalam token",
		})
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - format user_id tidak valid",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user_id bukan UUID valid",
		})
	}

	// 🎯 Ambil themes_or_levels_id dari path
	themesIDParam := c.Params("themes_or_levels_id")
	themesID, err := strconv.Atoi(themesIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "themes_or_levels_id tidak valid",
		})
	}

	// Step 1: Ambil data user_theme
	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Data user_theme tidak ditemukan",
		})
	}

	// Step 2: Ambil list unit_id dari TotalUnit
	var unitIDs []int64
	for _, id := range userTheme.TotalUnit {
		unitIDs = append(unitIDs, id)
	}

	// Step 3: Ambil unit + section_quizzes + quizzes
	var units []userModel.UnitModel
	if err := ctrl.DB.
		Preload("SectionQuizzes.Quizzes"). // 👈 ini tambahan penting
		Where("id IN ? AND themes_or_level_id = ?", unitIDs, themesID).
		Find(&units).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data unit, section_quizzes dan quizzes",
		})
	}

	// Step 4: Ambil user_units dengan preload SectionProgress
	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.
		Preload("SectionProgress", "user_id = ?", userID).
		Where("user_id = ? AND unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data progress unit",
		})
	}

	// Step 5: Ambil SectionProgress manual per unit
	for i := range userUnits {
		// Log: tampilkan TotalSectionQuizzes untuk debug
		log.Printf("[INFO] TotalSectionQuizzes untuk unit_id %d: %v", userUnits[i].UnitID, userUnits[i].TotalSectionQuizzes)

		var sectionProgress []userSectionQuizzesModel.UserSectionQuizzesModel

		// Cek apakah TotalSectionQuizzes ada isinya sebelum query
		if len(userUnits[i].TotalSectionQuizzes) > 0 {
			// Konversi pq.Int64Array ke []int64
			sectionIDs := make([]int64, len(userUnits[i].TotalSectionQuizzes))
			copy(sectionIDs, userUnits[i].TotalSectionQuizzes)

			if err := ctrl.DB.
				Where("user_id = ?", userUnits[i].UserID).
				Where("section_quizzes_id IN ?", sectionIDs).
				Find(&sectionProgress).Error; err != nil {
				log.Printf("[WARNING] Gagal ambil section_progress untuk unit_id %d: %v", userUnits[i].UnitID, err)
				continue
			}
			userUnits[i].SectionProgress = sectionProgress
			log.Printf("[SUCCESS] Berhasil ambil section_progress untuk unit_id %d: %d items", userUnits[i].UnitID, len(sectionProgress))
		} else {
			log.Printf("[INFO] Melewati section_progress karena TotalSectionQuizzes kosong untuk unit_id %d", userUnits[i].UnitID)
		}

	}

	// Step 6: Mapping unit_id -> user_unit (dengan SectionProgress yang sudah diisi)
	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		progressMap[u.UnitID] = u
	}

	// Step 7: Build response
	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}

	var result []ResponseUnit
	for _, unit := range units {
		progress := progressMap[unit.ID]
		result = append(result, ResponseUnit{
			UnitModel:    unit,
			UserProgress: progress,
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}