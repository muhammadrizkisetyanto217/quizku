package controller

import (
	"encoding/json"
	"fmt"
	"log"

	certificateModel "quizku/internals/features/certificates/certificate_versions/model"
	categoryModel "quizku/internals/features/lessons/categories/model"
	subcategoryModel "quizku/internals/features/lessons/subcategory/model"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"
	unitModel "quizku/internals/features/lessons/units/model"
	"quizku/internals/features/quizzes/quizzes/model"
	sectionQuizzesModel "quizku/internals/features/quizzes/quizzes/model"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserSubcategoryController struct {
	DB *gorm.DB
}

func NewUserSubcategoryController(db *gorm.DB) *UserSubcategoryController {
	return &UserSubcategoryController{DB: db}
}

func (ctrl *UserSubcategoryController) Create(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	userID, err := uuid.Parse(fmt.Sprintf("%v", userIDStr))
	if err != nil || userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID dari token tidak valid",
		})
	}

	type RequestBody struct {
		SubcategoryID uint `json:"subcategory_id"`
	}
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if body.SubcategoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SubcategoryID tidak boleh kosong atau nol",
		})
	}

	tx := ctrl.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memulai transaksi database",
		})
	}

	var subcategory subcategoryModel.SubcategoryModel
	if err := tx.First(&subcategory, body.SubcategoryID).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data subcategory",
		})
	}

	var certVersion struct {
		VersionNumber int
	}
	tx.Table("certificate_versions").
		Select("version_number").
		Where("subcategory_id = ?", body.SubcategoryID).
		Order("version_number DESC").
		Limit(1).
		Scan(&certVersion)

	now := time.Now()
	input := subcategoryModel.UserSubcategoryModel{
		UserID:         userID,
		SubcategoryID:  int(body.SubcategoryID),
		CreatedAt:      now,
		UpdatedAt:      now,
		CurrentVersion: certVersion.VersionNumber,
	}
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal simpan user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data user_subcategory",
		})
	}

	var themes []themesModel.ThemesOrLevelsModel
	if err := tx.Where("subcategories_id = ?", body.SubcategoryID).Find(&themes).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil themes:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data themes yang terkait",
		})
	}

	var themeIDs []uint
	var userThemes []themesModel.UserThemesOrLevelsModel
	for _, theme := range themes {
		themeIDs = append(themeIDs, theme.ID)
		userThemes = append(userThemes, themesModel.UserThemesOrLevelsModel{
			UserID:           userID,
			ThemesOrLevelsID: theme.ID,
			CompleteUnit:     datatypes.JSONMap{},
			GradeResult:      0,
			CreatedAt:        now,
		})
	}

	if len(userThemes) > 0 {
		if err := tx.CreateInBatches(&userThemes, 100).Error; err != nil {
			tx.Rollback()
			log.Println("[ERROR] Gagal simpan user_themes batch:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan data user_themes",
			})
		}
	}

	var units []unitModel.UnitModel
	if err := tx.Where("themes_or_level_id IN ?", themeIDs).Find(&units).Error; err != nil {
		tx.Rollback()
		log.Println("[ERROR] Gagal ambil units:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data unit",
		})
	}

	var userUnits []unitModel.UserUnitModel
	for _, unit := range units {
		var sectionQuizIDs []int64
		if err := tx.Model(&sectionQuizzesModel.SectionQuizzesModel{}).
			Where("unit_id = ?", unit.ID).
			Pluck("id", &sectionQuizIDs).Error; err != nil {
			tx.Rollback()
			log.Printf("[ERROR] Gagal ambil section_quizzes untuk unit_id %d: %v", unit.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Gagal mengambil section_quizzes untuk unit_id %d", unit.ID),
			})
		}

		userUnits = append(userUnits, unitModel.UserUnitModel{
			UserID:                 userID,
			UnitID:                 unit.ID,
			AttemptReading:         0,
			AttemptEvaluation:      datatypes.JSON([]byte(`{"attempt":0,"grade_evaluation":0}`)),
			CompleteSectionQuizzes: datatypes.JSON([]byte(`[]`)),
			GradeExam:              0,
			IsPassed:               false,
			GradeResult:            0,
			CreatedAt:              now,
			UpdatedAt:              now,
			SectionProgress:        []model.UserSectionQuizzesModel{},
		})
	}

	if len(userUnits) > 0 {
		if err := tx.CreateInBatches(&userUnits, 100).Error; err != nil {
			tx.Rollback()
			log.Println("[ERROR] Gagal simpan user_units batch:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan data user_units",
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("[ERROR] Commit transaksi gagal:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal commit transaksi database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "UserSubcategory, UserThemes, dan UserUnits berhasil dibuat",
		"data":    input,
	})
}

func (ctrl *UserSubcategoryController) GetByUserId(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID user tidak valid",
		})
	}

	var userSub subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.
		Select("*"). // <-- Ini penting untuk memastikan current_version ikut terambil
		Where("user_id = ?", userID).
		First(&userSub).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data tidak ditemukan",
			})
		}
		log.Println("[ERROR] Gagal ambil user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": userSub,
	})
}

// ✅ Refactored: grade_result & is_passed hanya diubah oleh service exam
func (ctrl *UserSubcategoryController) GetWithProgressByParam(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - user_id not found in token"})
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - invalid user_id format"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - invalid user_id UUID"})
	}

	difficultyID := c.Params("difficulty_id")
	if difficultyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "difficulty_id wajib diisi"})
	}

	var categories []categoryModel.CategoryModel
	if err := ctrl.DB.
		Where("difficulty_id = ?", difficultyID).
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active").Preload("ThemesOrLevels")
		}).
		Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	var userSubcat []subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&userSubcat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_subcategory"})
	}
	userSubcatMap := make(map[int]subcategoryModel.UserSubcategoryModel)
	for _, us := range userSubcat {
		existing, ok := userSubcatMap[us.SubcategoryID]
		if !ok || us.UpdatedAt.After(existing.UpdatedAt) {
			userSubcatMap[us.SubcategoryID] = us
		}
	}

	var userThemes []themesModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_themes_or_levels"})
	}
	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, ut := range userThemes {
		userThemeMap[ut.ThemesOrLevelsID] = ut
	}

	type ThemeWithProgress struct {
		ID               uint           `json:"id"`
		Name             string         `json:"name"`
		Status           string         `json:"status"`
		DescriptionShort string         `json:"description_short"`
		DescriptionLong  string         `json:"description_long"`
		TotalUnit        pq.Int64Array  `json:"total_unit"`
		ImageURL         string         `json:"image_url"`
		UpdateNews       datatypes.JSON `json:"update_news"`
		CreatedAt        time.Time      `json:"created_at"`
		UpdatedAt        *time.Time     `json:"updated_at"`
		SubcategoriesID  uint           `json:"subcategories_id"`
		GradeResult      int            `json:"grade_result"`
		CompleteUnit     datatypes.JSON `json:"complete_unit"`
		HasProgressTheme bool           `json:"has_progress_theme"`
	}

	type SubcategoryWithProgress struct {
		ID                     uint                `json:"id"`
		Name                   string              `json:"name"`
		Status                 string              `json:"status"`
		DescriptionLong        string              `json:"description_long"`
		TotalThemesOrLevels    pq.Int64Array       `json:"total_themes_or_levels"`
		ImageURL               string              `json:"image_url"`
		UpdateNews             datatypes.JSON      `json:"update_news"`
		CreatedAt              time.Time           `json:"created_at"`
		UpdatedAt              *time.Time          `json:"updated_at"`
		CategoriesID           uint                `json:"categories_id"`
		GradeResult            int                 `json:"grade_result"`
		CompleteThemesOrLevels datatypes.JSONMap   `json:"complete_themes_or_levels"`
		IssuedVersion          int                 `json:"issued_version"`
		CurrentVersion         *int                `json:"current_version"`
		UserSubcategoryID      uint                `json:"user_subcategory_id"`
		UserID                 uuid.UUID           `json:"user_id"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
		HasProgressSubcategory bool                `json:"has_progress_subcategory"`
	}

	type CategoryWithSubcat struct {
		ID                 uint                      `json:"id"`
		Name               string                    `json:"name"`
		Status             string                    `json:"status"`
		DescriptionShort   string                    `json:"description_short"`
		DescriptionLong    string                    `json:"description_long"`
		TotalSubcategories pq.Int64Array             `json:"total_subcategories"`
		ImageURL           string                    `json:"image_url"`
		UpdateNews         datatypes.JSON            `json:"update_news"`
		DifficultyID       uint                      `json:"difficulty_id"`
		CreatedAt          time.Time                 `json:"created_at"`
		UpdatedAt          *time.Time                `json:"updated_at"`
		Subcategories      []SubcategoryWithProgress `json:"subcategories"`
	}

	var result []CategoryWithSubcat
	totalGrade := 0
	totalCount := 0

	for _, cat := range categories {
		subcatWithProgress := []SubcategoryWithProgress{}

		for _, sub := range cat.Subcategories {
			us, hasProgress := userSubcatMap[int(sub.ID)]
			if !hasProgress {
				continue
			}

			var certVersion certificateModel.CertificateVersionModel
			var versionNumber *int = nil
			if err := ctrl.DB.
				Where("subcategory_id = ?", sub.ID).
				Order("version_number DESC").
				First(&certVersion).Error; err == nil {
				versionNumber = &certVersion.VersionNumber
			}

			issuedVersion := 0
			if versionNumber != nil {
				issuedVersion = *versionNumber
			}

			themes := []ThemeWithProgress{}
			for _, theme := range sub.ThemesOrLevels {
				userTheme := userThemeMap[theme.ID]
				rawJSON, _ := json.Marshal(userTheme.CompleteUnit)

				themes = append(themes, ThemeWithProgress{
					ID:               theme.ID,
					Name:             theme.Name,
					Status:           theme.Status,
					DescriptionShort: theme.DescriptionShort,
					DescriptionLong:  theme.DescriptionLong,
					TotalUnit:        theme.TotalUnit,
					ImageURL:         theme.ImageURL,
					CreatedAt:        theme.CreatedAt,
					UpdatedAt:        theme.UpdatedAt,
					SubcategoriesID:  uint(theme.SubcategoriesID),
					GradeResult:      userTheme.GradeResult,
					CompleteUnit:     datatypes.JSON(rawJSON),
					HasProgressTheme: userTheme.GradeResult > 0 || (userTheme.CompleteUnit != nil && len(userTheme.CompleteUnit) > 0),
				})

				if userTheme.GradeResult > 0 {
					totalGrade += userTheme.GradeResult
					totalCount++
				}
			}

			subcatWithProgress = append(subcatWithProgress, SubcategoryWithProgress{
				ID:                     sub.ID,
				Name:                   sub.Name,
				Status:                 sub.Status,
				DescriptionLong:        sub.DescriptionLong,
				TotalThemesOrLevels:    sub.TotalThemesOrLevels,
				ImageURL:               sub.ImageURL,
				// UpdateNews:             sub.UpdateNews,
				CreatedAt:              sub.CreatedAt,
				UpdatedAt:              sub.UpdatedAt,
				CategoriesID:           sub.CategoriesID,
				GradeResult:            us.GradeResult,
				CompleteThemesOrLevels: us.CompleteThemesOrLevels,
				IssuedVersion:          issuedVersion,
				CurrentVersion:         &us.CurrentVersion,
				UserSubcategoryID:      us.ID,
				UserID:                 userID,
				ThemesOrLevels:         themes,
				HasProgressSubcategory: true,
			})
		}

		result = append(result, CategoryWithSubcat{
			ID:                 cat.ID,
			Name:               cat.Name,
			Status:             cat.Status,
			DescriptionShort:   cat.DescriptionShort,
			DescriptionLong:    cat.DescriptionLong,
			TotalSubcategories: cat.TotalSubcategories,
			ImageURL:           cat.ImageURL,
			// UpdateNews:         cat.UpdateNews,
			DifficultyID:       cat.DifficultyID,
			CreatedAt:          cat.CreatedAt,
			// UpdatedAt:          cat.UpdatedAt,
			Subcategories:      subcatWithProgress,
		})
	}

	type CombinedProgress struct {
		UserID       uuid.UUID `json:"user_id"`
		AverageGrade int       `json:"average_grade"`
		DataCount    int       `json:"data_count"`
	}
	combined := CombinedProgress{
		UserID:       userID,
		AverageGrade: 0,
		DataCount:    totalCount,
	}
	if totalCount > 0 {
		combined.AverageGrade = totalGrade / totalCount
	}

	return c.JSON(fiber.Map{
		"message":       "Berhasil ambil data lengkap",
		"data":          result,
		"user_progress": combined,
	})
}
