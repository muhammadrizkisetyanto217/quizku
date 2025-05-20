package controller

import (
	"encoding/json"
	"log"

	certificateModel "quizku/internals/features/certificates/certificate_versions/model"
	categoryModel "quizku/internals/features/lessons/categories/model"
	subcategoryModel "quizku/internals/features/lessons/subcategories/model"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"
	unitModel "quizku/internals/features/lessons/units/model"

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

// 🟢 CREATE USER SUBCATEGORY: Inisialisasi user_subcategory, user_themes, dan user_units saat pertama kali user memilih subkategori
// 🟢 CREATE USER SUBCATEGORY: Inisialisasi user_subcategory, user_themes, dan user_units saat pertama kali user memilih subkategori
func (ctrl *UserSubcategoryController) Create(c *fiber.Ctx) error {
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

	type RequestBody struct {
		SubcategoryID uint `json:"subcategory_id"`
	}
	var body RequestBody
	if err := c.BodyParser(&body); err != nil || body.SubcategoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SubcategoryID tidak boleh kosong atau nol"})
	}

	tx := ctrl.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi database"})
	}

	var subcategory subcategoryModel.SubcategoryModel
	if err := tx.First(&subcategory, body.SubcategoryID).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Subcategory tidak ditemukan"})
	}

	var certVersion struct{ CertificateVersionNumber int }
	tx.Table("certificate_versions").
		Select("certificate_version_number").
		Where("subcategory_id = ?", body.SubcategoryID).
		Order("certificate_version_number DESC").
		Limit(1).
		Scan(&certVersion)

	now := time.Now()
	userSubcat := subcategoryModel.UserSubcategoryModel{
		UserSubcategoryUserID:         userID,
		UserSubcategorySubcategoryID:  int(body.SubcategoryID),
		CreatedAt:                     now,
		UpdatedAt:                     now,
		UserSubcategoryCurrentVersion: certVersion.CertificateVersionNumber,
	}
	if err := tx.Create(&userSubcat).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user_subcategory"})
	}

	var themes []themesModel.ThemesOrLevelsModel
	if err := tx.Where("themes_or_level_subcategory_id = ?", body.SubcategoryID).Find(&themes).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data themes"})
	}

	var themeIDs []uint
	var userThemes []themesModel.UserThemesOrLevelsModel
	for _, theme := range themes {
		themeIDs = append(themeIDs, theme.ThemesOrLevelID)
		userThemes = append(userThemes, themesModel.UserThemesOrLevelsModel{
			UserThemeUserID:          userID,
			UserThemeThemesOrLevelID: theme.ThemesOrLevelID,
			UserThemeCompleteUnit:    datatypes.JSONMap{},
			UserThemeGradeResult:     0,
			CreatedAt:                now,
		})
	}
	if len(userThemes) > 0 {
		if err := tx.CreateInBatches(&userThemes, 100).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user_themes"})
		}
	}

	var units []unitModel.UnitModel
	if err := tx.Where("unit_themes_or_level_id IN ?", themeIDs).Find(&units).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data unit"})
	}

	var userUnits []unitModel.UserUnitModel
	for _, unit := range units {
		userUnits = append(userUnits, unitModel.UserUnitModel{
			UserUnitUserID:                 userID,
			UserUnitUnitID:                 unit.UnitID,
			UserUnitAttemptReading:         0,
			UserUnitAttemptEvaluation:      datatypes.JSON([]byte(`{"evaluation_attempt_count":0,"evaluation_final_grade":0}`)),
			UserUnitCompleteSectionQuizzes: datatypes.JSON([]byte("[]")),
			UserUnitGradeQuiz:              0,
			UserUnitGradeExam:              0,
			UserUnitGradeResult:            0,
			UserUnitIsPassed:               false,
			CreatedAt:                      now,
			UpdatedAt:                      now,
		})
	}
	if len(userUnits) > 0 {
		if err := tx.CreateInBatches(&userUnits, 100).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user_units"})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal commit transaksi"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "UserSubcategory, UserThemes, dan UserUnits berhasil dibuat",
		"data":    userSubcat,
	})
}

// 🟢 GET USER_SUBCATEGORY BY USER ID: Ambil data user_subcategory milik user tertentu (hanya 1 record terbaru)
// 🟢 GET USER_SUBCATEGORY BY USER ID: Ambil data user_subcategory milik user tertentu (hanya 1 record terbaru)
func (ctrl *UserSubcategoryController) GetByUserId(c *fiber.Ctx) error {
	id := c.Params("id")

	// 🔐 Validasi UUID dari parameter
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID user tidak valid",
		})
	}

	// 🔍 Ambil data user_subcategory milik user tersebut (hanya satu record terbaru)
	var userSub subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.
		Where("user_subcategory_user_id = ?", userID).
		Order("updated_at DESC").
		First(&userSub).Error; err != nil {

		// 🛑 Jika data tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data tidak ditemukan",
			})
		}

		// ❌ Jika error lain saat query
		log.Println("[ERROR] Gagal ambil user_subcategory:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data",
		})
	}

	// ✅ Kirim data jika ditemukan
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": userSub,
	})
}

// 🟢 GET USER PROGRESS WITH CATEGORY STRUCTURE: Ambil data lengkap kategori, subkategori, themes, dan progress user
// ✅ grade_result & is_passed hanya berasal dari service exam, bukan dihitung ulang di sini
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
		Where("category_difficulty_id = ?", difficultyID).
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("subcategory_status = ?", "active").Preload("ThemesOrLevels")
		}).
		Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	var userSubcats []subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_subcategory_user_id = ?", userID).Order("updated_at DESC").Find(&userSubcats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_subcategory"})
	}
	userSubcatMap := make(map[int]subcategoryModel.UserSubcategoryModel)
	for _, us := range userSubcats {
		if existing, ok := userSubcatMap[us.UserSubcategorySubcategoryID]; !ok || us.UpdatedAt.After(existing.UpdatedAt) {
			userSubcatMap[us.UserSubcategorySubcategoryID] = us
		}
	}

	var userThemes []themesModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_theme_user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_themes_or_levels"})
	}
	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, ut := range userThemes {
		if existing, ok := userThemeMap[ut.UserThemeThemesOrLevelID]; !ok || ut.UpdatedAt.After(existing.UpdatedAt) {
			userThemeMap[ut.UserThemeThemesOrLevelID] = ut
		}
	}

	type UserThemeProgress struct {
		ThemeID               uint           `json:"theme_id"`
		ThemeName             string         `json:"theme_name"`
		ThemeStatus           string         `json:"theme_status"`
		ThemeShortDesc        string         `json:"theme_short_description"`
		ThemeLongDesc         string         `json:"theme_long_description"`
		ThemeTotalUnits       pq.Int64Array  `json:"theme_total_units"`
		ThemeImageURL         string         `json:"theme_image_url"`
		ThemeSubcategoryID    uint           `json:"theme_subcategory_id"`
		UserThemeGradeResult  int            `json:"user_theme_grade_result"`
		UserThemeCompleteUnit datatypes.JSON `json:"user_theme_complete_unit"`
		UserHasThemeProgress  bool           `json:"user_has_theme_progress"`
	}
	type UserSubcategoryProgress struct {
		SubcategoryID              uint                `json:"subcategory_id"`
		SubcategoryName            string              `json:"subcategory_name"`
		SubcategoryStatus          string              `json:"subcategory_status"`
		SubcategoryLongDesc        string              `json:"subcategory_long_description"`
		SubcategoryTotalThemes     pq.Int64Array       `json:"subcategory_total_themes"`
		SubcategoryImageURL        string              `json:"subcategory_image_url"`
		CreatedAt                  time.Time           `json:"created_at"`
		UpdatedAt                  *time.Time          `json:"updated_at"`
		CategoryID                 uint                `json:"category_id"`
		UserSubcategoryGradeResult int                 `json:"user_subcategory_grade_result"`
		UserSubcategoryCompleted   datatypes.JSONMap   `json:"user_subcategory_completed"`
		CertificateVersionIssued   int                 `json:"certificate_version_issued"`
		CertificateVersionCurrent  *int                `json:"certificate_version_current"`
		UserSubcategoryID          uint                `json:"user_subcategory_id"`
		UserID                     uuid.UUID           `json:"user_id"`
		ThemesWithProgress         []UserThemeProgress `json:"themes_with_progress"`
		UserHasSubcategoryProgress bool                `json:"user_has_subcategory_progress"`
	}
	type CategoryWithUserProgress struct {
		CategoryID            uint                      `json:"category_id"`
		CategoryName          string                    `json:"category_name"`
		CategoryStatus        string                    `json:"category_status"`
		CategoryShortDesc     string                    `json:"category_short_description"`
		CategoryLongDesc      string                    `json:"category_long_description"`
		CategoryTotalSub      pq.Int64Array             `json:"category_total_sub"`
		CategoryImageURL      string                    `json:"category_image_url"`
		CategoryDifficultyID  uint                      `json:"category_difficulty_id"`
		CreatedAt             time.Time                 `json:"created_at"`
		UpdatedAt             *time.Time                `json:"updated_at"`
		SubcategoriesProgress []UserSubcategoryProgress `json:"subcategories_progress"`
	}

	var result []CategoryWithUserProgress
	totalThemeGrade := 0
	totalThemeWithProgress := 0

	for _, cat := range categories {
		subcatProgressList := []UserSubcategoryProgress{}
		for _, sub := range cat.Subcategories {
			us, hasProgress := userSubcatMap[int(sub.SubcategoryID)]
			if !hasProgress {
				continue
			}

			var certVersion certificateModel.CertificateVersionModel
			var versionNumber *int = nil
			if err := ctrl.DB.Where("cert_versions_subcategory_id = ?", sub.SubcategoryID).Order("cert_versions_number DESC").First(&certVersion).Error; err == nil {
				versionNumber = &certVersion.CertVersionNumber
			}
			issuedVersion := 0
			if versionNumber != nil {
				issuedVersion = *versionNumber
			}

			themeProgressList := []UserThemeProgress{}
			for _, theme := range sub.ThemesOrLevels {
				ut, ok := userThemeMap[theme.ThemesOrLevelID]
				if !ok {
					continue
				}
				rawJSON, _ := json.Marshal(ut.UserThemeCompleteUnit)
				themeProgressList = append(themeProgressList, UserThemeProgress{
					ThemeID:               theme.ThemesOrLevelID,
					ThemeName:             theme.ThemesOrLevelName,
					ThemeStatus:           theme.ThemesOrLevelStatus,
					ThemeShortDesc:        theme.ThemesOrLevelDescriptionShort,
					ThemeLongDesc:         theme.ThemesOrLevelDescriptionLong,
					ThemeTotalUnits:       theme.ThemesOrLevelTotalUnit,
					ThemeImageURL:         theme.ThemesOrLevelImageURL,
					ThemeSubcategoryID:    uint(theme.ThemesOrLevelSubcategoryID),
					UserThemeGradeResult:  ut.UserThemeGradeResult,
					UserThemeCompleteUnit: datatypes.JSON(rawJSON),
					UserHasThemeProgress:  ut.UserThemeGradeResult > 0 || (ut.UserThemeCompleteUnit != nil && len(ut.UserThemeCompleteUnit) > 0),
				})
				if ut.UserThemeGradeResult > 0 {
					totalThemeGrade += ut.UserThemeGradeResult
					totalThemeWithProgress++
				}
			}

			subcatProgressList = append(subcatProgressList, UserSubcategoryProgress{
				SubcategoryID:              sub.SubcategoryID,
				SubcategoryName:            sub.SubcategoryName,
				SubcategoryStatus:          sub.SubcategoryStatus,
				SubcategoryLongDesc:        sub.SubcategoryDescriptionLong,
				SubcategoryTotalThemes:     sub.SubcategoryTotalThemesOrLevels,
				SubcategoryImageURL:        sub.SubcategoryImageURL,
				CreatedAt:                  sub.CreatedAt,
				UpdatedAt:                  sub.UpdatedAt,
				CategoryID:                 sub.SubcategoryCategoryID,
				UserSubcategoryGradeResult: us.UserSubcategoryGradeResult,
				UserSubcategoryCompleted:   us.UserSubcategoryCompleteThemesOrLevels,
				CertificateVersionIssued:   issuedVersion,
				CertificateVersionCurrent:  &us.UserSubcategoryCurrentVersion,
				UserSubcategoryID:          us.UserSubcategoryID,
				UserID:                     userID,
				ThemesWithProgress:         themeProgressList,
				UserHasSubcategoryProgress: true,
			})
		}

		result = append(result, CategoryWithUserProgress{
			CategoryID:            cat.CategoryID,
			CategoryName:          cat.CategoryName,
			CategoryStatus:        cat.CategoryStatus,
			CategoryShortDesc:     cat.CategoryDescriptionShort,
			CategoryLongDesc:      cat.CategoryDescriptionLong,
			CategoryTotalSub:      cat.CategoryTotalSubcategories,
			CategoryImageURL:      cat.CategoryImageURL,
			CategoryDifficultyID:  cat.CategoryDifficultyID,
			CreatedAt:             cat.CreatedAt,
			UpdatedAt:             &cat.UpdatedAt,
			SubcategoriesProgress: subcatProgressList,
		})
	}

	type UserCombinedProgressSummary struct {
		UserID                  uuid.UUID `json:"user_id"`
		TotalThemesWithProgress int       `json:"total_themes_with_progress"`
		AccumulatedThemesGrade  int       `json:"accumulated_themes_grade"`
		AverageThemeGrade       int       `json:"average_theme_grade"`
	}

	summary := UserCombinedProgressSummary{
		UserID:                  userID,
		TotalThemesWithProgress: totalThemeWithProgress,
		AccumulatedThemesGrade:  totalThemeGrade,
		AverageThemeGrade:       0,
	}
	if totalThemeWithProgress > 0 {
		summary.AverageThemeGrade = totalThemeGrade / totalThemeWithProgress
	}

	return c.JSON(fiber.Map{
		"message":       "Berhasil ambil data lengkap",
		"data":          result,
		"user_progress": summary,
	})
}
