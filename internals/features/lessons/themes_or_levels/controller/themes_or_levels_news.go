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
	id := c.Params("themes_or_levels_id") // Ambil parameter ID dari URL
	var news []model.ThemesOrLevelsNewsModel

	// Query berdasarkan themes_or_levels_id dan pastikan belum terhapus (soft delete)
	if err := tc.DB.
		Where("themes_or_levels_id = ?", id).
		Where("deleted_at IS NULL").
		Find(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Jika tidak ditemukan data apapun, kembalikan 404
	if len(news) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "No news found for this themes_or_levels_id",
		})
	}

	// Kirim data news berdasarkan themes_or_levels_id
	return c.JSON(fiber.Map{
		"message": "News for the selected themes/levels retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET /themes-or-levels-news/:id
// Mengambil satu berita berdasarkan ID unik
func (tc *ThemesOrLevelsNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	// Ambil satu data berdasarkan primary key
	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	// Kirim respons jika ditemukan
	return c.JSON(fiber.Map{
		"message": "Themes/Levels news found successfully",
		"data":    news,
	})
}

// 🟡 POST /themes-or-levels-news
// Menambahkan berita baru untuk themes_or_levels tertentu.
// Setelah berhasil disimpan, akan memperbarui field JSON `update_news` di tabel themes_or_levels.
func (tc *ThemesOrLevelsNewsController) Create(c *fiber.Ctx) error {
	var news model.ThemesOrLevelsNewsModel

	// Parse body dari request
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Simpan berita baru ke database
	if err := tc.DB.Create(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Perbarui field update_news di tabel themes_or_levels (cache JSON)
	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Themes/Levels news created successfully",
		"data":    news,
	})
}

// 🟠 PUT /themes-or-levels-news/:id
// Mengupdate isi berita berdasarkan ID, lalu menyegarkan field `update_news` di themes_or_levels
func (tc *ThemesOrLevelsNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	// Cari data lama berdasarkan ID
	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	// Overwrite data lama dengan body baru dari request
	if err := c.BodyParser(&news); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Simpan perubahan ke database
	if err := tc.DB.Save(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Update field update_news (JSON array) di themes_or_levels
	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.JSON(fiber.Map{
		"message": "Themes/Levels news updated successfully",
		"data":    news,
	})
}

// 🔴 DELETE /themes-or-levels-news/:id
// Menghapus satu berita berdasarkan ID, lalu menyegarkan field `update_news` di tabel themes_or_levels
func (tc *ThemesOrLevelsNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.ThemesOrLevelsNewsModel

	// Cari data berdasarkan ID
	if err := tc.DB.First(&news, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Themes/Levels news not found",
		})
	}

	// Hapus data (soft delete jika pakai gorm.Model dengan DeletedAt)
	if err := tc.DB.Delete(&news).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Perbarui cache JSON pada themes_or_levels
	updateThemesOrLevelsNewsJSON(tc.DB, news.ThemesOrLevelsID)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Themes/Levels news with ID %v deleted successfully", news.ID),
	})
}

// ⚙️ updateThemesOrLevelsNewsJSON
// Helper internal untuk menyegarkan field `update_news` pada tabel themes_or_levels.
// Field ini menyimpan array JSON dari semua news aktif terkait theme tersebut.
// Tujuannya adalah untuk efisiensi frontend yang hanya perlu membaca satu kolom.
func updateThemesOrLevelsNewsJSON(db *gorm.DB, themeID uint) {
	var newsList []model.ThemesOrLevelsNewsModel

	// Ambil semua news untuk theme terkait, urutkan berdasarkan created_at terbaru
	if err := db.Where("themes_or_level_id = ?", themeID).
		Order("created_at desc").
		Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch themes/levels news for update:", err)
		return
	}

	// Ubah ke bentuk JSON array
	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal themes/levels news:", err)
		return
	}

	// Simpan ke field update_news (kolom JSON) di themes_or_levels
	res := db.Table("themes_or_levels").
		Where("id = ?", themeID).
		Update("update_news", datatypes.JSON(newsData))

	log.Println("[DEBUG] Rows affected (themes_or_levels):", res.RowsAffected)
}
