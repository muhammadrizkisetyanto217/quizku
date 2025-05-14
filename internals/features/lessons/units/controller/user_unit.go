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

// ✅ Refactored: grade_result & is_passed hanya diubah oleh service exam
// ✅ Refactored: grade_result & is_passed hanya diubah oleh service exam
// ✅ Refactored: grade_result & is_passed hanya diubah oleh service exam
func (ctrl *UserUnitController) GetUserUnitsByThemesOrLevels(c *fiber.Ctx) error {
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

	themesIDParam := c.Params("themes_or_levels_id")
	themesID, err := strconv.Atoi(themesIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "themes_or_levels_id tidak valid",
		})
	}

	var userTheme themesOrLevelsModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ? AND themes_or_levels_id = ?", userID, themesID).
		First(&userTheme).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Data user_theme tidak ditemukan",
		})
	}

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

	var unitIDs []uint
	sectionQuizToUnit := make(map[uint]uint)
	for _, unit := range units {
		unitIDs = append(unitIDs, unit.ID)
		for _, section := range unit.SectionQuizzes {
			sectionQuizToUnit[section.ID] = unit.ID
		}
	}

	var userUnits []userModel.UserUnitModel
	if err := ctrl.DB.
		Where("user_id = ? AND unit_id IN ?", userID, unitIDs).
		Find(&userUnits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal ambil data progress unit",
		})
	}

	// Ambil semua section_progress dalam 1 query saja
	var allSectionProgress []userSectionQuizzesModel.UserSectionQuizzesModel
	if err := ctrl.DB.
		Where("user_id = ?", userID).
		Where("section_quizzes_id IN ?", keys(sectionQuizToUnit)).
		Find(&allSectionProgress).Error; err != nil {
		log.Printf("[ERROR] Gagal ambil seluruh section_progress user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil seluruh section_progress user",
		})
	}

	// Kelompokkan section_progress berdasarkan unit_id
	progressPerUnit := make(map[uint][]userSectionQuizzesModel.UserSectionQuizzesModel)
	completedMap := make(map[uint][]int64)
	for _, sp := range allSectionProgress {
		unitID := sectionQuizToUnit[sp.SectionQuizzesID]
		progressPerUnit[unitID] = append(progressPerUnit[unitID], sp)
		if len(sp.CompleteQuiz) > 0 {
			completedMap[unitID] = append(completedMap[unitID], int64(sp.SectionQuizzesID))
		}
	}

	progressMap := make(map[uint]userModel.UserUnitModel)
	for _, u := range userUnits {
		u.SectionProgress = progressPerUnit[u.UnitID]
		progressMap[u.UnitID] = u

		// Update ke database kalau ada yang lengkap
		if completed, ok := completedMap[u.UnitID]; ok && len(completed) > 0 {
			_ = ctrl.DB.Model(&userModel.UserUnitModel{}).
				Where("id = ?", u.ID).
				Update("complete_section_quizzes", pq.Int64Array(completed)).Error
		}
	}

	type ResponseUnit struct {
		userModel.UnitModel
		UserProgress userModel.UserUnitModel `json:"user_progress"`
	}
	var result []ResponseUnit
	for _, unit := range units {
		userProgress := progressMap[unit.ID]
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
	out := make([]uint, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
