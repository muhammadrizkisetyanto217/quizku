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

	model "quizku/internals/features/certificates/user_certificates/model"
	categoryModel "quizku/internals/features/lessons/categories/model"
	subcategoryModel "quizku/internals/features/lessons/subcategories/model"
	themesModel "quizku/internals/features/lessons/themes_or_levels/model"
	unitModel "quizku/internals/features/lessons/units/model"
	userProfileModel "quizku/internals/features/users/user/model"

	issuedCertificateService "quizku/internals/features/certificates/user_certificates/service"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	var cert model.UserCertificate
	// 🔍 Gunakan kolom semantik jika sudah refactor model
	if err := ctrl.DB.Where("user_cert_id = ?", id).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sertifikat tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Detail sertifikat ditemukan",
		"data":    cert,
	})
}

// ✅ Untuk User: Get all certificates miliknya sendiri
// ✅ Untuk User: Get all certificates miliknya sendiri
func (ctrl *IssuedCertificateController) GetByID(c *fiber.Ctx) error {
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

	var issuedCerts []model.UserCertificate
	if err := ctrl.DB.Where("user_cert_user_id = ?", userID).Find(&issuedCerts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil sertifikat"})
	}
	issuedMap := make(map[uint]model.UserCertificate)
	for _, cert := range issuedCerts {
		issuedMap[cert.UserCertSubcategoryID] = cert
	}

	var categories []categoryModel.CategoryModel
	if err := ctrl.DB.Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
		return db.Where("subcategory_status = ?", "active").Preload("ThemesOrLevels")
	}).Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil kategori"})
	}

	var userSubcats []subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_subcategory_user_id = ?", userID).Find(&userSubcats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil user_subcategory"})
	}
	userSubcatMap := make(map[int]subcategoryModel.UserSubcategoryModel)
	for _, us := range userSubcats {
		existing, ok := userSubcatMap[us.UserSubcategorySubcategoryID]
		if !ok || us.UpdatedAt.After(existing.UpdatedAt) {
			userSubcatMap[us.UserSubcategorySubcategoryID] = us
		}
	}

	var userThemes []themesModel.UserThemesOrLevelsModel
	if err := ctrl.DB.Where("user_theme_user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil user_themes_or_levels"})
	}
	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, ut := range userThemes {
		userThemeMap[ut.UserThemeThemesOrLevelID] = ut
	}

	var units []unitModel.UnitModel
	if err := ctrl.DB.Find(&units).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil units"})
	}
	unitMap := make(map[uint][]unitModel.UnitModel)
	for _, u := range units {
		unitMap[u.UnitThemesOrLevelID] = append(unitMap[u.UnitThemesOrLevelID], u)
	}

	type VersionMap struct {
		SubcategoryID uint
		VersionNumber int
	}
	var versionList []VersionMap
	if err := ctrl.DB.Table("certificate_versions").Select("cert_versions_subcategory_id, MAX(cert_versions_number) as version_number").Group("cert_versions_subcategory_id").Scan(&versionList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil versi sertifikat"})
	}
	versionMap := make(map[uint]int)
	for _, v := range versionList {
		versionMap[v.SubcategoryID] = v.VersionNumber
	}

	type ThemeWithProgress struct {
		ThemeID              uint                  `json:"theme_id"`
		ThemeName            string                `json:"theme_name"`
		ThemeStatus          string                `json:"theme_status"`
		ThemeShortDesc       string                `json:"theme_short_description"`
		ThemeLongDesc        string                `json:"theme_long_description"`
		ThemeTotalUnits      pq.Int64Array         `json:"theme_total_units"`
		ThemeImageURL        string                `json:"theme_image_url"`
		ThemeCreatedAt       time.Time             `json:"theme_created_at"`
		ThemeUpdatedAt       *time.Time            `json:"theme_updated_at"`
		ThemeSubcategoryID   uint                  `json:"theme_subcategory_id"`
		ThemeGradeResult     int                   `json:"theme_grade_result"`
		ThemeCompleteUnit    datatypes.JSON        `json:"theme_complete_unit"`
		UserHasThemeProgress bool                  `json:"user_has_theme_progress"`
		Units                []unitModel.UnitModel `json:"theme_units"`
	}
	type SubcategoryWithProgress struct {
		SubcategoryID      uint                `json:"subcategory_id"`
		SubcategoryName    string              `json:"subcategory_name"`
		SubcategoryStatus  string              `json:"subcategory_status"`
		SubcategoryDesc    string              `json:"subcategory_description"`
		SubcategoryUnits   pq.Int64Array       `json:"subcategory_total_themes"`
		SubcategoryImage   string              `json:"subcategory_image_url"`
		SubcategoryCreated time.Time           `json:"subcategory_created_at"`
		SubcategoryUpdated *time.Time          `json:"subcategory_updated_at"`
		CategoriesID       uint                `json:"category_id"`
		SubcategoryGrade   int                 `json:"subcategory_grade_result"`
		CompletedThemes    datatypes.JSONMap   `json:"completed_themes"`
		IssuedVersion      int                 `json:"certificate_version_issued"`
		CurrentVersion     *int                `json:"certificate_version_current"`
		UserSubcategoryID  uint                `json:"user_subcategory_id"`
		UserID             uuid.UUID           `json:"user_id"`
		ThemesProgress     []ThemeWithProgress `json:"themes_with_progress"`
		UserHasProgress    bool                `json:"user_has_subcategory_progress"`
		IssuedAt           time.Time           `json:"certificate_issued_at"`
		SlugURL            string              `json:"certificate_slug_url"`
		IsUpToDate         bool                `json:"user_cert_is_up_to_date"`
	}
	type CategoriesWithSubcat struct {
		CategoriesID         uint                      `json:"category_id"`
		CategoriesName       string                    `json:"category_name"`
		CategoriesStatus     string                    `json:"category_status"`
		CategoriesShort      string                    `json:"category_short_description"`
		CategoriesLong       string                    `json:"category_long_description"`
		CategoriesSubTotal   pq.Int64Array             `json:"category_total_subcategories"`
		CategoriesImage      string                    `json:"category_image_url"`
		CategoriesDifficulty uint                      `json:"category_difficulty_id"`
		CreatedAt            time.Time                 `json:"category_created_at"`
		UpdatedAt            *time.Time                `json:"category_updated_at"`
		Subcategories        []SubcategoryWithProgress `json:"subcategories_progress"`
	}

	var result []CategoriesWithSubcat
	for _, cat := range categories {
		subcatList := []SubcategoryWithProgress{}
		for _, sub := range cat.Subcategories {
			issued, ok := issuedMap[sub.SubcategoryID]
			if !ok {
				continue
			}
			us, hasProgress := userSubcatMap[int(sub.SubcategoryID)]
			if !hasProgress {
				continue
			}

			themes := []ThemeWithProgress{}
			for _, theme := range sub.ThemesOrLevels {
				ut := userThemeMap[theme.ThemesOrLevelID]
				rawJSON, err := json.Marshal(ut.UserThemeCompleteUnit)
				if err != nil {
					log.Println("[WARNING] Marshal theme_complete_unit gagal:", err)
					rawJSON = []byte("{}")
				}

				themes = append(themes, ThemeWithProgress{
					ThemeID:              theme.ThemesOrLevelID,
					ThemeName:            theme.ThemesOrLevelName,
					ThemeStatus:          theme.ThemesOrLevelStatus,
					ThemeShortDesc:       theme.ThemesOrLevelDescriptionShort,
					ThemeLongDesc:        theme.ThemesOrLevelDescriptionLong,
					ThemeTotalUnits:      theme.ThemesOrLevelTotalUnit,
					ThemeImageURL:        theme.ThemesOrLevelImageURL,
					ThemeCreatedAt:       theme.CreatedAt,
					ThemeUpdatedAt:       theme.UpdatedAt,
					ThemeSubcategoryID:   uint(theme.ThemesOrLevelSubcategoryID),
					ThemeGradeResult:     ut.UserThemeGradeResult,
					ThemeCompleteUnit:    datatypes.JSON(rawJSON),
					UserHasThemeProgress: ut.UserThemeGradeResult > 0 || (ut.UserThemeCompleteUnit != nil && len(ut.UserThemeCompleteUnit) > 0),
					Units:                unitMap[theme.ThemesOrLevelID],
				})
			}

			subcatList = append(subcatList, SubcategoryWithProgress{
				SubcategoryID:      sub.SubcategoryID,
				SubcategoryName:    sub.SubcategoryName,
				SubcategoryStatus:  sub.SubcategoryStatus,
				SubcategoryDesc:    sub.SubcategoryDescriptionLong,
				SubcategoryUnits:   sub.SubcategoryTotalThemesOrLevels,
				SubcategoryImage:   sub.SubcategoryImageURL,
				SubcategoryCreated: sub.CreatedAt,
				SubcategoryUpdated: sub.UpdatedAt,
				CategoriesID:       sub.SubcategoryCategoryID,
				SubcategoryGrade:   us.UserSubcategoryGradeResult,
				CompletedThemes:    us.UserSubcategoryCompleteThemesOrLevels,
				IssuedVersion:      versionMap[sub.SubcategoryID],
				CurrentVersion:     &us.UserSubcategoryCurrentVersion,
				UserSubcategoryID:  us.UserSubcategoryID,
				UserID:             userID,
				ThemesProgress:     themes,
				UserHasProgress:    true,
				IssuedAt:           issued.UserCertIssuedAt,
				SlugURL:            issued.UserCertSlugURL,
				IsUpToDate:         issued.UserCertIsUpToDate,
			})
		}
		if len(subcatList) > 0 {
			result = append(result, CategoriesWithSubcat{
				CategoriesID:         cat.CategoryID,
				CategoriesName:       cat.CategoryName,
				CategoriesStatus:     cat.CategoryName,
				CategoriesShort:      cat.CategoryDescriptionShort,
				CategoriesLong:       cat.CategoryDescriptionLong,
				CategoriesSubTotal:   cat.CategoryTotalSubcategories,
				CategoriesImage:      cat.CategoryImageURL,
				CategoriesDifficulty: cat.CategoryDifficultyID,
				CreatedAt:            cat.CreatedAt,
				UpdatedAt:            &cat.UpdatedAt,
				Subcategories:        subcatList,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data sertifikat lengkap",
		"data":    result,
	})
}

// ✅ GET /api/certificates/by-subcategory/:subcategory_id
// ✅ Untuk User: Ambil sertifikat miliknya berdasarkan subcategory_id
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

	var cert model.UserCertificate
	if err := ctrl.DB.Where("user_cert_user_id = ? AND user_cert_subcategory_id = ?", userID, subcategoryID).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sertifikat tidak ditemukan"})
	}

	var sub subcategoryModel.SubcategoryModel
	if err := ctrl.DB.Preload("ThemesOrLevels").First(&sub, subcategoryID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil subcategory"})
	}

	var us subcategoryModel.UserSubcategoryModel
	if err := ctrl.DB.Where("user_subcategory_user_id = ? AND user_subcategory_subcategory_id = ?", userID, subcategoryID).Order("updated_at DESC").First(&us).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil progress user_subcategory"})
	}

	var themeIDs []uint
	for _, theme := range sub.ThemesOrLevels {
		themeIDs = append(themeIDs, theme.ThemesOrLevelID)
	}

	userThemeMap := make(map[uint]themesModel.UserThemesOrLevelsModel)
	for _, themeID := range themeIDs {
		var ut themesModel.UserThemesOrLevelsModel
		err := ctrl.DB.
			Where("user_theme_user_id = ? AND user_theme_themes_or_level_id = ?", userID, themeID).
			Order("updated_at DESC").
			First(&ut).Error
		if err == nil {
			userThemeMap[themeID] = ut
		}
	}

	var units []unitModel.UnitModel
	if err := ctrl.DB.Find(&units).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil units"})
	}
	unitMap := make(map[uint][]unitModel.UnitModel)
	for _, u := range units {
		unitMap[u.UnitThemesOrLevelID] = append(unitMap[u.UnitThemesOrLevelID], u)
	}

	var versionIssued int
	_ = ctrl.DB.Table("certificate_versions").
		Select("MAX(cert_versions_number)").
		Where("cert_versions_subcategory_id = ?", subcategoryID).
		Scan(&versionIssued)

	isUpToDate, err := issuedCertificateService.CheckAndUpdateIsUpToDate(ctrl.DB, userID, subcategoryID, cert, us, sub, versionIssued)
	if err != nil {
		log.Println("[WARNING] Gagal validasi IsUpToDate:", err.Error())
		isUpToDate = cert.UserCertIsUpToDate
	}

	type UnitTitleOnly struct {
		ID   uint   `json:"unit_id"`
		Name string `json:"unit_name"`
	}

	type ThemeWithProgress struct {
		ThemesOrLevelID            uint            `json:"themes_or_level_id"`
		ThemesOrLevelName          string          `json:"themes_or_level_name"`
		ThemesOrLevelStatus        string          `json:"themes_or_level_status"`
		ThemesOrLevelTotalUnit     pq.Int64Array   `json:"themes_or_level_total_unit"`
		ThemesOrLevelSubcategoryID int             `json:"themes_or_level_subcategory_id"`
		UserThemeGradeResult       int             `json:"user_theme_grade_result"`
		UserThemeCompleteUnit      datatypes.JSON  `json:"user_theme_complete_unit"`
		UserHasThemeProgress       bool            `json:"user_has_theme_progress"`
		Units                      []UnitTitleOnly `json:"units"`
	}

	type SubcategoryWithProgress struct {
		ID                     uint                `json:"subcategory_id"`
		Name                   string              `json:"subcategory_name"`
		FullName               string              `json:"full_name"`
		Status                 string              `json:"subcategory_status"`
		DescriptionLong        string              `json:"subcategory_description_long"`
		TotalThemesOrLevels    pq.Int64Array       `json:"subcategory_total_themes_or_levels"`
		CategoriesID           uint                `json:"categories_id"`
		GradeResult            int                 `json:"user_subcategory_grade_result"`
		CompleteThemesOrLevels datatypes.JSONMap   `json:"user_subcategory_completed"`
		IssuedVersion          int                 `json:"certificate_version_issued"`
		CurrentVersion         *int                `json:"certificate_version_current"`
		UserSubcategoryID      uint                `json:"user_subcategory_id"`
		UserID                 uuid.UUID           `json:"user_id"`
		ThemesOrLevels         []ThemeWithProgress `json:"themes_or_levels"`
		HasProgressSubcategory bool                `json:"user_has_subcategory_progress"`
		IssuedAt               time.Time           `json:"user_cert_issued_at"`
		SlugURL                string              `json:"user_cert_slug_url"`
		IsUpToDate             bool                `json:"user_cert_is_up_to_date"`
	}

	var themes []ThemeWithProgress
	for _, theme := range sub.ThemesOrLevels {
		ut, ok := userThemeMap[theme.ThemesOrLevelID]
		var gradeResult int
		var completeUnit datatypes.JSON
		var hasProgress bool

		if ok {
			gradeResult = ut.UserThemeGradeResult
			rawJSON, _ := json.Marshal(ut.UserThemeCompleteUnit)
			completeUnit = datatypes.JSON(rawJSON)
			hasProgress = gradeResult > 0 || (ut.UserThemeCompleteUnit != nil && len(ut.UserThemeCompleteUnit) > 0)
		} else {
			gradeResult = 0
			completeUnit = datatypes.JSON([]byte("{}"))
			hasProgress = false
		}

		var unitTitles []UnitTitleOnly
		for _, u := range unitMap[theme.ThemesOrLevelID] {
			unitTitles = append(unitTitles, UnitTitleOnly{ID: u.UnitID, Name: u.UnitName})
		}

		themes = append(themes, ThemeWithProgress{
			ThemesOrLevelID:            theme.ThemesOrLevelID,
			ThemesOrLevelName:          theme.ThemesOrLevelName,
			ThemesOrLevelStatus:        theme.ThemesOrLevelStatus,
			ThemesOrLevelTotalUnit:     theme.ThemesOrLevelTotalUnit,
			ThemesOrLevelSubcategoryID: theme.ThemesOrLevelSubcategoryID,
			UserThemeGradeResult:       gradeResult,
			UserThemeCompleteUnit:      completeUnit,
			UserHasThemeProgress:       hasProgress,
			Units:                      unitTitles,
		})
	}

	currentVersionPtr := func(v int) *int {
		if v > 0 {
			return &v
		}
		return nil
	}(us.UserSubcategoryCurrentVersion)

	return c.JSON(fiber.Map{
		"message": "Berhasil ambil data sertifikat berdasarkan subcategory",
		"data": SubcategoryWithProgress{
			ID:                     sub.SubcategoryID,
			Name:                   sub.SubcategoryName,
			FullName:               profile.FullName,
			Status:                 sub.SubcategoryStatus,
			DescriptionLong:        sub.SubcategoryDescriptionLong,
			TotalThemesOrLevels:    sub.SubcategoryTotalThemesOrLevels,
			CategoriesID:           sub.SubcategoryCategoryID,
			GradeResult:            us.UserSubcategoryGradeResult,
			CompleteThemesOrLevels: us.UserSubcategoryCompleteThemesOrLevels,
			IssuedVersion:          versionIssued,
			CurrentVersion:         currentVersionPtr,
			UserSubcategoryID:      us.UserSubcategoryID,
			UserID:                 userID,
			ThemesOrLevels:         themes,
			HasProgressSubcategory: true,
			IssuedAt:               cert.UserCertIssuedAt,
			SlugURL:                cert.UserCertSlugURL,
			IsUpToDate:             isUpToDate,
		},
	})
}

// ✅ Untuk Public: Get certificate by slug (tanpa login)
// ✅ Untuk Public: Get certificate by slug (tanpa login)
func (ctrl *IssuedCertificateController) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug sertifikat tidak boleh kosong",
		})
	}

	// 🔍 Cari sertifikat berdasarkan slug
	var cert model.UserCertificate
	if err := ctrl.DB.Where("user_cert_slug_url = ?", slug).First(&cert).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Sertifikat tidak ditemukan",
		})
	}

	// 🔍 Ambil current version dari tabel user_subcategory
	var currentVersion int
	err := ctrl.DB.Table("user_subcategory").
		Select("current_version").
		Where("user_id = ? AND subcategory_id = ?", cert.UserCertUserID, cert.UserCertSubcategoryID).
		Scan(&currentVersion).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil versi terkini user_subcategory",
		})
	}

	// 🔍 Ambil versi sertifikat terbaru dari certificate_versions
	var latestVersionIssued int
	err = ctrl.DB.Table("certificate_versions").
		Select("MAX(cert_versions_number)").
		Where("cert_versions_subcategory_id = ?", cert.UserCertSubcategoryID).
		Scan(&latestVersionIssued).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil versi sertifikat terbaru",
		})
	}

	// ✅ Struct respons publik yang semantik
	type PublicCertificateResponse struct {
		CertificateUserID              uuid.UUID `json:"certificate_user_id"`
		CertificateSubcategoryID       uint      `json:"certificate_subcategory_id"`
		CertificateSlugURL             string    `json:"certificate_slug_url"`
		CertificateIssuedAt            time.Time `json:"certificate_issued_at"`
		CertificateIsUpToDate          bool      `json:"certificate_is_up_to_date"`
		CertificateVersionCurrentSaved int       `json:"certificate_version_current"`
		CertificateVersionLatestIssued int       `json:"certificate_version_issued"`
	}

	// 🚀 Kirim respons
	response := PublicCertificateResponse{
		CertificateUserID:              cert.UserCertUserID,
		CertificateSubcategoryID:       cert.UserCertSubcategoryID,
		CertificateSlugURL:             cert.UserCertSlugURL,
		CertificateIssuedAt:            cert.UserCertIssuedAt,
		CertificateIsUpToDate:          cert.UserCertIsUpToDate,
		CertificateVersionCurrentSaved: currentVersion,
		CertificateVersionLatestIssued: latestVersionIssued,
	}

	return c.JSON(response)
}
