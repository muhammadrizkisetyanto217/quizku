package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"quizku/internals/features/lessons/themes_or_levels/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ThemesOrLevelsNewsController struct {
	DB *gorm.DB
}

func NewThemesOrLevelsNewsController(db *gorm.DB) *ThemesOrLevelsNewsController {
	return &ThemesOrLevelsNewsController{DB: db}
}

// 🟢 GET /themes-or-levels-news
// Mengambil seluruh daftar news (berita/pengumuman) untuk semua themes_or_levels
func (tc *ThemesOrLevelsNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.ThemesOrLevelsNewsModel

	// Ambil semua data dari tabel themes_or_levels_news
	if err := tc.DB.Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Jika berhasil, kirim response dengan seluruh data news
	return c.JSON(fiber.Map{
		"message": "Themes/Levels news list retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET /themes-or-levels-news/themes-or-levels/:themes_or_levels_id
// Mengambil semua news yang terkait dengan satu themes_or_levels tertentu
func (tc *ThemesOrLevelsNewsController) GetByThemesOrLevelsID(c *fiber.Ctx) error {
	themesOrLevelID := c.Params("themes_or_levels_id")
	var newsList []model.ThemesOrLevelsNewsModel

	if err := tc.DB.
		Where("themes_news_themes_or_level_id = ?", themesOrLevelID).
		Where("deleted_at IS NULL").
		Find(&newsList).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal mengambil news berdasarkan themes_or_levels_id: " + err.Error(),
		})
	}

	if len(newsList) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News tidak ditemukan untuk themes_or_levels_id tersebut",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Daftar news untuk themes_or_levels berhasil diambil",
		"data":    newsList,
	})
}

// 🟢 GET /themes-or-levels-news/:id
// Mengambil satu berita berdasarkan ID unik
func (tc *ThemesOrLevelsNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	if err := tc.DB.First(&news, "themes_news_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News tidak ditemukan untuk ID tersebut",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News berhasil ditemukan",
		"data":    news,
	})
}

// 🟡 POST /themes-or-levels-news
// Menambahkan berita baru untuk themes_or_levels tertentu.
// Setelah berhasil disimpan, akan memperbarui field JSON `update_news` di tabel themes_or_levels.
func (tc *ThemesOrLevelsNewsController) Create(c *fiber.Ctx) error {
	var input model.ThemesOrLevelsNewsModel

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Format body tidak valid",
		})
	}

	if err := tc.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menyimpan news: " + err.Error(),
		})
	}

	// 🧠 Perbarui cache update_news pada tabel themes_or_levels
	updateThemesOrLevelsNewsJSON(tc.DB, input.ThemesNewsThemesOrLevelID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "News berhasil ditambahkan",
		"data":    input,
	})
}

// 🟠 PUT /themes-or-levels-news/:id
// Mengupdate isi berita berdasarkan ID, lalu menyegarkan field `update_news` di themes_or_levels
func (tc *ThemesOrLevelsNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var existing model.ThemesOrLevelsNewsModel

	// 🔍 Cari data lama berdasarkan ID
	if err := tc.DB.First(&existing, "themes_news_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News tidak ditemukan",
		})
	}

	// 📥 Parsing body request untuk overwrite field (tanpa overwrite ID)
	var input model.ThemesOrLevelsNewsModel
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Format request body tidak valid",
		})
	}

	// 📝 Update field satu per satu untuk menghindari overwrite ID
	existing.ThemesNewsTitle = input.ThemesNewsTitle
	existing.ThemesNewsDescription = input.ThemesNewsDescription
	existing.ThemesNewsIsPublic = input.ThemesNewsIsPublic
	existing.ThemesNewsThemesOrLevelID = input.ThemesNewsThemesOrLevelID

	// 💾 Simpan perubahan ke database
	if err := tc.DB.Save(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal memperbarui news: " + err.Error(),
		})
	}

	// 🔄 Update cache update_news di themes_or_levels
	updateThemesOrLevelsNewsJSON(tc.DB, existing.ThemesNewsThemesOrLevelID)

	return c.JSON(fiber.Map{
		"message": "News berhasil diperbarui",
		"data":    existing,
	})
}

// 🔴 DELETE /themes-or-levels-news/:id
// Menghapus satu berita berdasarkan ID, lalu menyegarkan field `update_news` di tabel themes_or_levels
func (tc *ThemesOrLevelsNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	// 🔍 Cari data berdasarkan ID
	if err := tc.DB.First(&news, "themes_news_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "News tidak ditemukan",
		})
	}

	// 🗑️ Hapus data (soft delete)
	if err := tc.DB.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal menghapus news: " + err.Error(),
		})
	}

	// 🔄 Perbarui field JSON `update_news` di themes_or_levels
	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesNewsThemesOrLevelID)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("News dengan ID %v berhasil dihapus", news.ThemesNewsID),
	})
}

// ⚙️ updateThemesOrLevelsNewsJSON
// Helper internal untuk menyegarkan field `themes_or_level_update_news` pada tabel themes_or_levels.
// Field ini menyimpan array JSON dari semua news aktif terkait theme tersebut.
// Tujuannya adalah untuk efisiensi frontend yang hanya perlu membaca satu kolom.
func updateThemesOrLevelsNewsJSON(db *gorm.DB, themeID uint) {
	var newsList []model.ThemesOrLevelsNewsModel

	// 🔍 Ambil semua news aktif (belum terhapus) untuk theme terkait, urutkan berdasarkan created_at terbaru
	if err := db.
		Where("themes_news_themes_or_level_id = ?", themeID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Gagal mengambil daftar news untuk update:", err)
		return
	}

	// 🔄 Konversi menjadi JSON array
	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Gagal mengubah news menjadi JSON:", err)
		return
	}

	// 💾 Simpan JSON ke kolom themes_or_level_update_news pada tabel themes_or_levels
	if err := db.Table("themes_or_levels").
		Where("themes_or_level_id = ?", themeID).
		Update("themes_or_level_update_news", datatypes.JSON(newsData)).Error; err != nil {
		log.Println("[ERROR] Gagal update kolom themes_or_level_update_news:", err)
	} else {
		log.Printf("[INFO] Kolom update_news berhasil disegarkan untuk theme_id: %d", themeID)
	}
}
