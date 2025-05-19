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
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserUnitController struct {
	DB *gorm.DB
}

func NewUserUnitController(db *gorm.DB) *UserUnitController {
	return &UserUnitController{DB: db}
}

// 🟢 GET /api/user-units/:user_id
// Mengambil semua data progres unit milik user berdasarkan user_unit_user_id.
// Data yang dikembalikan termasuk relasi SectionProgress per unit.
func (ctrl *UserUnitController) GetByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	// Validasi UUID
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	var data []userModel.UserUnitModel

	// Ambil semua user_unit berdasarkan user_unit_user_id
	if err := ctrl.DB.
		Preload("SectionProgress").
		Where("user_unit_user_id = ?", userID).
		Find(&data).Error; err != nil {

		log.Println("[ERROR] Gagal ambil data user_unit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"total": len(data),
		"data":  data,
	})
}

// 🟢 GET /api/user-units/themes/:themes_or_levels_id
// Mengambil seluruh unit dalam sebuah theme (themes_or_levels_id) beserta progres user di tiap unit.
// Progress meliputi section_progress, complete_section_quizzes, dan field lain dari user_unit.
func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevels(c *fiber.Ctx) error {
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

	// Ambil themes_or_levels_id dari URL
	themesIDParam := c.Params("themes_or_levels_id")
	themesID, err := strconv.Atoi(themesIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "themes_or_levels_id tidak valid",
		})
	}

	// Pastikan user punya entri user_themes_or_levels
	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_themes_or_levels_user_id = ? AND user_themes_or_levels_themes_or_levels_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Data user_theme tidak ditemukan",
		})
	}

	// Ambil semua unit dalam theme tersebut
	var units []userModel.UnitModel
	if err := ctrl.DB.
		Preload("SectionQuizzes").
		Preload("SectionQuizzes.Quizzes").
		Where("unit_themes_or_level_id = ?", themesID).
		Find(&units).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data unit, section_quizzes dan quizzes",
		})
	}

	// Buat mapping unit_id & section_quiz_id
	var unitIDs []uint
	sectionQuizToUnit := make(map[uint]uint)
	for _, unit := range units {
		unitIDs = append(unitIDs, unit.UnitID)
		for _, section := range unit.SectionQuizzes {
			sectionQuizToUnit[section.SectionQuizzesID] = unit.UnitID
		}
	}

	// Ambil user_unit berdasarkan user_id dan unit_id
	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.
		Where("user_unit_user_id = ? AND user_unit_unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data progress unit",
		})
	}

	// Ambil seluruh section_progress
	var allSectionProgress []userSectionQuizzesModel.UserSectionQuizzesModel
	if err := ctrl.DB.
		Where("user_section_quizzes_user_id = ?", userID).
		Where("user_section_quizzes_section_quizzes_id IN ?", keys(sectionQuizToUnit)).
		Find(&allSectionProgress).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil seluruh section_progress user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil seluruh section_progress user",
		})
	}

	// Mapping progress dan section selesai
	progressPerUnit := make(map[uint][]userSectionQuizzesModel.UserSectionQuizzesModel)
	completedMap := make(map[uint][]int64)
	for _, sp := range allSectionProgress {
		sectionID := sp.UserSectionQuizzesSectionQuizzesID
		unitID := sectionQuizToUnit[sectionID]
		progressPerUnit[unitID] = append(progressPerUnit[unitID], sp)
		if len(sp.UserSectionQuizzesCompleteQuiz) > 0 {
			completedMap[unitID] = append(completedMap[unitID], int64(sectionID))
		}
	}

	// Gabungkan progress ke user_unit
	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		u.SectionProgress = progressPerUnit[u.UserUnitUnitID]
		progressMap[u.UserUnitUnitID] = u

		if completed, ok := completedMap[u.UserUnitUnitID]; ok && len(completed) > 0 {
			_ = ctrl.DB.Model(&userModel.UserUnitModel{}).
				Where("user_unit_id = ?", u.UserUnitID).
				Update("user_unit_complete_section_quizzes", pq.Int64Array(completed)).Error
		}
	}

	// Gabungkan ke response akhir
	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}
	var result []ResponseUnit
	for _, unit := range units {
		userProgress := progressMap[unit.UnitID]
		result = append(result, ResponseUnit{
			UnitModel:    unit,
			UserProgress: userProgress,
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func keys(m map[uint]uint) []uint {
	k := make([]uint, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
