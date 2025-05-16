package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	model "quizku/internals/features/certificates/issued_certificates/model"
	categoryModel "quizku/internals/features/lessons/categories/model"
	subcategoryModel "quizku/internals/features/lessons/subcategories/model"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"
	unitModel "quizku/internals/features/lessons/units/model"
	userProfileModel "quizku/internals/features/users/user/model"

	issuedCertificateService "quizku/internals/features/certificates/issued_certificates/service"
)

type IssuedCertificateController struct {
	DB *gorm.DB
}

func NewIssuedCertificateController(db *gorm.DB) *IssuedCertificateController {
	return &IssuedCertificateController{DB: db}
}

// ✅ GET /api/certificates/:id
// ✅ GetByIDUser: Ambil detail sertifikat berdasarkan ID (hanya untuk admin atau keperluan umum)
func (ctrl *IssuedCertificateController) GetByIDUser(c *fiber.Ctx) error {
	// 🔹 Ambil parameter ID dari URL
	idStr := c.Params("id")

	// 🔍 Konversi ke integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// ❌ Jika gagal konversi, kirim error 400
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	// 🔍 Ambil data sertifikat dari database berdasarkan ID
	var cert model.IssuedCertificateModel
	if err := ctrl.DB.First(&cert, id).Error; err != nil {
		// ❌ Jika tidak ditemukan, kirim error 404
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Sertifikat tidak ditemukan",
		})
	}

	// ✅ Berhasil ambil sertifikat, kirim response JSON
	return c.JSON(fiber.Map{
		"message": "Detail sertifikat ditemukan",
		"data":    cert,
	})
}


// ✅ Untuk User: Get all certificates miliknya sendiri
func (ctrl *IssuedCertificateController) GetByID(c *fiber.Ctx) error {
	// 🔐 Ambil user_id dari token (autentikasi)
	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id format"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	// ✅ Ambil semua sertifikat milik user
	var issuedCerts []model.IssuedCertificateModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&issuedCerts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil sertifikat"})
	}
	issuedMap := make(map[uint]model.IssuedCertificateModel)
	for _, cert := range issuedCerts {
		issuedMap[cert.SubcategoryID] = cert
	}

	// ✅ Ambil semua kategori (beserta subkategori & themes aktif)
	var categories []categoryModel.CategoryModel
	if err := ctrl.DB.
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "active").Preload("ThemesOrLevels")
		}).
		Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	// ✅ Ambil semua progress user_subcategory
	var userSubcats []subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&userSubcats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil user_subcategory"})
	}
	userSubcatMap := make(map[int]subcategoryModel.UserSubcategoryModel)
	for _, us := range userSubcats {
		existing, ok := userSubcatMap[us.SubcategoryID]
		if !ok || us.UpdatedAt.After(existing.UpdatedAt) {
			userSubcatMap[us.SubcategoryID] = us
		}
	}

	// ✅ Ambil semua progress user_themes_or_levels
	var userThemes []themesModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil user_themes_or_levels"})
	}
	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, ut := range userThemes {
		userThemeMap[ut.ThemesOrLevelsID] = ut
	}

	// ✅ Ambil semua unit
	var units []unitModel.UnitModel
	if err := ctrl.DB.Find(&units).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil units"})
	}
	unitMap := make(map[uint][]unitModel.UnitModel)
	for _, u := range units {
		unitMap[u.ThemesOrLevelID] = append(unitMap[u.ThemesOrLevelID], u)
	}

	// ✅ Ambil versi maksimum sertifikat per subkategori
	type VersionMap struct {
		SubcategoryID uint
		VersionNumber int
	}
	var versionList []VersionMap
	if err := ctrl.DB.
		Table("certificate_versions").
		Select("subcategory_id, MAX(version_number) as version_number").
		Group("subcategory_id").
		Scan(&versionList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil versi sertifikat"})
	}
	versionMap := make(map[uint]int)
	for _, v := range versionList {
		versionMap[v.SubcategoryID] = v.VersionNumber
	}

	// ✅ Struct untuk response: ThemeWithProgress, SubcategoryWithProgress, CategoryWithSubcat
	type ThemeWithProgress struct {
		ID               uint                  `json:"id"`
		Name             string                `json:"name"`
		Status           string                `json:"status"`
		DescriptionShort string                `json:"description_short"`
		DescriptionLong  string                `json:"description_long"`
		TotalUnit        pq.Int64Array         `json:"total_unit"`
		ImageURL         string                `json:"image_url"`
		UpdateNews       datatypes.JSON        `json:"update_news"`
		CreatedAt        time.Time             `json:"created_at"`
		UpdatedAt        *time.Time            `json:"updated_at"`
		SubcategoriesID  uint                  `json:"subcategories_id"`
		GradeResult      int                   `json:"grade_result"` // ✅ nilai dari exam
		CompleteUnit     datatypes.JSON        `json:"complete_unit"`
		HasProgressTheme bool                  `json:"has_progress_theme"`
		Units            []unitModel.UnitModel `json:"units"`
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
		GradeResult            int                 `json:"grade_result"` // ✅ nilai akhir dari semua theme (hasil exam)
		CompleteThemesOrLevels datatypes.JSONMap   `json:"complete_themes_or_levels"`
		IssuedVersion          int                 `json:"issued_version"`
		CurrentVersion         *int                `json:"current_version"`
		UserSubcategoryID      uint                `json:"user_subcategory_id"`
		UserID                 uuid.UUID           `json:"user_id"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
		HasProgressSubcategory bool                `json:"has_progress_subcategory"`
		IssuedAt               time.Time           `json:"issued_at"`
		SlugURL                string              `json:"slug_url"`
		IsUpToDate             bool                `json:"is_up_to_date"`
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

	// ✅ Bangun respons akhir
	var result []CategoryWithSubcat
	for _, cat := range categories {
		var subcatList []SubcategoryWithProgress

		for _, sub := range cat.Subcategories {
			if _, ok := issuedMap[sub.ID]; !ok {
				continue // ❌ Skip jika tidak punya sertifikat
			}
			us, hasProgress := userSubcatMap[int(sub.ID)]
			if !hasProgress {
				continue
			}

			var themes []ThemeWithProgress
			for _, theme := range sub.ThemesOrLevels {
				ut := userThemeMap[theme.ID]
				rawJSON, _ := json.Marshal(ut.CompleteUnit)

				themes = append(themes, ThemeWithProgress{
					ID:               theme.ID,
					Name:             theme.Name,
					Status:           theme.Status,
					DescriptionShort: theme.DescriptionShort,
					DescriptionLong:  theme.DescriptionLong,
					TotalUnit:        theme.TotalUnit,
					ImageURL:         theme.ImageURL,
					// UpdateNews:       theme.UpdateNews,
					CreatedAt:        theme.CreatedAt,
					UpdatedAt:        theme.UpdatedAt,
					SubcategoriesID:  uint(theme.SubcategoriesID),
					GradeResult:      ut.GradeResult, // ✅ Nilai hasil exam per theme
					CompleteUnit:     datatypes.JSON(rawJSON),
					HasProgressTheme: ut.GradeResult > 0 || (ut.CompleteUnit != nil && len(ut.CompleteUnit) > 0),
					Units:            unitMap[theme.ID],
				})
			}

			currentVersionPtr := func(v int) *int {
				if v > 0 {
					return &v
				}
				return nil
			}(us.CurrentVersion)

			issued := issuedMap[sub.ID]

			subcatList = append(subcatList, SubcategoryWithProgress{
				ID:                  sub.ID,
				Name:                sub.Name,
				Status:              sub.Status,
				DescriptionLong:     sub.DescriptionLong,
				TotalThemesOrLevels: sub.TotalThemesOrLevels,
				ImageURL:            sub.ImageURL,
				// UpdateNews:             sub.UpdateNews,
				CreatedAt:              sub.CreatedAt,
				UpdatedAt:              sub.UpdatedAt,
				CategoriesID:           sub.CategoriesID,
				GradeResult:            us.GradeResult, // ✅ Nilai akhir subkategori (total)
				CompleteThemesOrLevels: us.CompleteThemesOrLevels,
				IssuedVersion:          versionMap[sub.ID],
				CurrentVersion:         currentVersionPtr,
				UserSubcategoryID:      us.ID,
				UserID:                 userID,
				ThemesOrLevels:         themes,
				HasProgressSubcategory: true,
				IssuedAt:               issued.IssuedAt,
				SlugURL:                issued.SlugURL,
				IsUpToDate:             issued.IsUpToDate,
			})
		}

		if len(subcatList) > 0 {
			result = append(result, CategoryWithSubcat{
				ID:                 cat.ID,
				Name:               cat.Name,
				Status:             cat.Status,
				DescriptionShort:   cat.DescriptionShort,
				DescriptionLong:    cat.DescriptionLong,
				TotalSubcategories: cat.TotalSubcategories,
				ImageURL:           cat.ImageURL,
				// UpdateNews:         cat.UpdateNews,
				DifficultyID: cat.DifficultyID,
				CreatedAt:    cat.CreatedAt,
				// UpdatedAt:          cat.UpdatedAt,
				Subcategories: subcatList,
			})
		}
	}

	// ✅ Kirim respons JSON
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data sertifikat lengkap",
		"data":    result,
	})
}

func (ctrl *IssuedCertificateController) GetBySubcategoryID(c *fiber.Ctx) error {
	subcategoryIDStr := c.Params("subcategory_id")
	subcategoryID, err := strconv.Atoi(subcategoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subcategory_id"})
	}

	userIDVal := c.Locals("user_id")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id format"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid UUID"})
	}

	var profile userProfileModel.UsersProfileModel
	if err := ctrl.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil profile user"})
	}

	var cert model.IssuedCertificateModel
	if err := ctrl.DB.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sertifikat tidak ditemukan"})
	}

	var sub subcategoryModel.SubcategoryModel
	if err := ctrl.DB.Preload("ThemesOrLevels").First(&sub, subcategoryID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil subcategory"})
	}

	var us subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_id = ? AND subcategory_id = ?", userID, subcategoryID).Order("updated_at DESC").First(&us).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_subcategory"})
	}

	// Ambil theme IDs
	var themeIDs []uint
	for _, theme := range sub.ThemesOrLevels {
		themeIDs = append(themeIDs, theme.ID)
	}

	// ✅ Loop per theme untuk ambil user progress paling baru (native GORM)
	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, themeID := range themeIDs {
		var ut themesModel.UserThemesOrLevelsModel
		err := ctrl.DB.
			Where("user_id = ? AND themes_or_levels_id = ?", userID, themeID).
			Order("updated_at DESC").
			First(&ut).Error
		if err == nil {
			userThemeMap[themeID] = ut
		}
	}

	// Ambil semua unit
	var units []unitModel.UnitModel
	if err := ctrl.DB.Find(&units).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil units"})
	}
	unitMap := make(map[uint][]unitModel.UnitModel)
	for _, u := range units {
		unitMap[u.ThemesOrLevelID] = append(unitMap[u.ThemesOrLevelID], u)
	}

	// Ambil versi issued tertinggi
	var versionIssued int
	_ = ctrl.DB.Table("certificate_versions").
		Select("MAX(version_number)").
		Where("subcategory_id = ?", subcategoryID).
		Scan(&versionIssued)

		// 🔁 Cek dan update is_up_to_date jika perlu
	isUpToDate, err := issuedCertificateService.CheckAndUpdateIsUpToDate(ctrl.DB, userID, subcategoryID, cert, us, sub, versionIssued)
	if err != nil {
		log.Println("[WARNING] Gagal validasi IsUpToDate:", err.Error())
		isUpToDate = cert.IsUpToDate // fallback ke yang ada
	}

	type UnitTitleOnly struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	type ThemeWithProgress struct {
		ID               uint            `json:"id"`
		Name             string          `json:"name"`
		Status           string          `json:"status"`
		TotalUnit        pq.Int64Array   `json:"total_unit"`
		SubcategoriesID  uint            `json:"subcategories_id"`
		GradeResult      int             `json:"grade_result"`
		CompleteUnit     datatypes.JSON  `json:"complete_unit"`
		HasProgressTheme bool            `json:"has_progress_theme"`
		Units            []UnitTitleOnly `json:"units"`
	}
	type SubcategoryWithProgress struct {
		ID                     uint                `json:"id"`
		Name                   string              `json:"name"`
		FullName               string              `json:"full_name"`
		Status                 string              `json:"status"`
		DescriptionLong        string              `json:"description_long"`
		TotalThemesOrLevels    pq.Int64Array       `json:"total_themes_or_levels"`
		CategoriesID           uint                `json:"categories_id"`
		GradeResult            int                 `json:"grade_result"`
		CompleteThemesOrLevels datatypes.JSONMap   `json:"complete_themes_or_levels"`
		IssuedVersion          int                 `json:"issued_version"`
		CurrentVersion         *int                `json:"current_version"`
		UserSubcategoryID      uint                `json:"user_subcategory_id"`
		UserID                 uuid.UUID           `json:"user_id"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
		HasProgressSubcategory bool                `json:"has_progress_subcategory"`
		IssuedAt               time.Time           `json:"issued_at"`
		SlugURL                string              `json:"slug_url"`
		IsUpToDate             bool                `json:"is_up_to_date"`
	}

	// Susun data theme beserta progress-nya
	var themes []ThemeWithProgress
	for _, theme := range sub.ThemesOrLevels {
		ut, ok := userThemeMap[theme.ID]
		var gradeResult int
		var completeUnit datatypes.JSON
		var hasProgress bool

		if ok {
			gradeResult = ut.GradeResult
			rawJSON, _ := json.Marshal(ut.CompleteUnit)
			completeUnit = datatypes.JSON(rawJSON)
			hasProgress = gradeResult > 0 || (ut.CompleteUnit != nil && len(ut.CompleteUnit) > 0)
		} else {
			gradeResult = 0
			completeUnit = datatypes.JSON([]byte("{}"))
			hasProgress = false
		}

		var unitTitles []UnitTitleOnly
		for _, u := range unitMap[theme.ID] {
			unitTitles = append(unitTitles, UnitTitleOnly{
				ID:   u.ID,
				Name: u.Name,
			})
		}

		themes = append(themes, ThemeWithProgress{
			ID:               theme.ID,
			Name:             theme.Name,
			Status:           theme.Status,
			TotalUnit:        theme.TotalUnit,
			SubcategoriesID:  sub.ID,
			GradeResult:      gradeResult,
			CompleteUnit:     completeUnit,
			HasProgressTheme: hasProgress,
			Units:            unitTitles,
		})
	}

	// Convert current_version jadi pointer
	currentVersionPtr := func(v int) *int {
		if v > 0 {
			return &v
		}
		return nil
	}(us.CurrentVersion)

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil data sertifikat berdasarkan subcategory",
		"data": SubcategoryWithProgress{
			ID:                     sub.ID,
			Name:                   sub.Name,
			Status:                 sub.Status,
			FullName:               profile.FullName,
			DescriptionLong:        sub.DescriptionLong,
			TotalThemesOrLevels:    sub.TotalThemesOrLevels,
			CategoriesID:           sub.CategoriesID,
			GradeResult:            us.GradeResult,
			CompleteThemesOrLevels: us.CompleteThemesOrLevels,
			IssuedVersion:          versionIssued,
			CurrentVersion:         currentVersionPtr,
			UserSubcategoryID:      us.ID,
			UserID:                 userID,
			ThemesOrLevels:         themes,
			HasProgressSubcategory: true,
			IssuedAt:               cert.IssuedAt,
			SlugURL:                cert.SlugURL,
			IsUpToDate:             isUpToDate,
		},
	})
}

// ✅ Untuk Public: Get certificate by slug (tanpa login)
func (ctrl *IssuedCertificateController) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug tidak boleh kosong"})
	}

	var cert model.IssuedCertificateModel
	if err := ctrl.DB.Where("slug_url = ?", slug).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sertifikat tidak ditemukan"})
	}

	var versionCurrent int
	err := ctrl.DB.Table("user_subcategory").
		Select("current_version").
		Where("user_id = ? AND subcategory_id = ?", cert.UserID, cert.SubcategoryID).
		Scan(&versionCurrent).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil current_version"})
	}

	var versionIssued int
	err = ctrl.DB.Table("certificate_versions").
		Select("version_number").
		Where("subcategory_id = ?", cert.SubcategoryID).
		Order("version_number DESC").
		Limit(1).
		Scan(&versionIssued).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil version_issued"})
	}

	type PublicCertificateResponse struct {
		model.IssuedCertificateModel
		VersionCurrent int `json:"version_current"`
		VersionIssued  int `json:"version_issued"`
	}
	resp := PublicCertificateResponse{
		IssuedCertificateModel: cert,
		VersionCurrent:         versionCurrent,
		VersionIssued:          versionIssued,
	}

	return c.JSON(resp)
}
