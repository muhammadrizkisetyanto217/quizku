package controller

import (
	"fmt"
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

func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevels(c *fiber.Ctx) error {
	// 🔐 Ambil user_id dari JWT
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

	// Step 1: Ambil user_theme
	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Data user_theme tidak ditemukan",
		})
	}

	// Step 2: Ambil semua units berdasarkan themes_or_levels_id
	var units []userModel.UnitModel
	if err := ctrl.DB.
		Preload("SectionQuizzes").
		Preload("SectionQuizzes.Quizzes").
		Where("themes_or_level_id = ?", themesID).
		Find(&units).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data unit, section_quizzes dan quizzes",
		})
	}

	// Step 3: Ambil user_unit
	var unitIDs []uint
	for _, unit := range units {
		unitIDs = append(unitIDs, unit.ID)
	}

	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.
		Preload("SectionProgress", "user_id = ?", userID).
		Where("user_id = ? AND unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data progress unit",
		})
	}

	// Step 4: Ambil SectionProgress per unit
	for i := range userUnits {
		// Ambil `TotalSectionQuizzes` langsung dari `unit.SectionQuizzes`
		var sectionQuizIDs []uint // Mengubah tipe menjadi uint, karena sectionQuiz.ID adalah uint
		for _, sectionQuiz := range units[i].SectionQuizzes {
			sectionQuizIDs = append(sectionQuizIDs, sectionQuiz.ID)
		}

		// Ambil SectionProgress berdasarkan section_quizzes_id dari `sectionQuizIDs`
		var sectionProgress []userSectionQuizzesModel.UserSectionQuizzesModel
		if len(sectionQuizIDs) > 0 {
			if err := ctrl.DB.
				Where("user_id = ?", userUnits[i].UserID).
				Where("section_quizzes_id IN ?", sectionQuizIDs).
				Find(&sectionProgress).Error; err != nil {
				log.Printf("[ERROR] Gagal ambil section_quizzes untuk user_id %v: %v", userUnits[i].UserID, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Gagal mengambil section_quizzes untuk user_id %v", userUnits[i].UserID),
				})
			}
			userUnits[i].SectionProgress = sectionProgress
		}
	}

	// Step 5: Map unit_id → user_unit
	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		progressMap[u.UnitID] = u
	}

	// Step 6: Build response
	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}
	var result []ResponseUnit
	for _, unit := range units {
		result = append(result, ResponseUnit{
			UnitModel:    unit,
			UserProgress: progressMap[unit.ID],
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}
