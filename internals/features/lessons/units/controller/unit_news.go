package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"quizku/internals/features/lessons/units/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UnitNewsController struct {
	DB *gorm.DB
}

func NewUnitNewsController(db *gorm.DB) *UnitNewsController {
	return &UnitNewsController{DB: db}
}

// 🟢 GET /api/unit-news
// Mengambil seluruh daftar berita atau pengumuman yang terkait dengan semua unit.
// Biasanya digunakan oleh admin atau frontend untuk menampilkan daftar lengkap unit news.
func (uc *UnitNewsController) GetAll(c *fiber.Ctx) error {
	var news []model.UnitNewsModel

	// Ambil semua data dari tabel unit_news tanpa filter
	if err := uc.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unit news list retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET /api/unit-news/unit/:unit_id
// Mengambil seluruh berita yang terkait dengan unit tertentu berdasarkan unit_id.
// Biasanya digunakan untuk menampilkan pengumuman terkini per unit di halaman pembelajaran.
func (uc *UnitNewsController) GetByUnitID(c *fiber.Ctx) error {
	unitID := c.Params("unit_id")
	var news []model.UnitNewsModel

	// Ambil semua news yang aktif (belum terhapus) berdasarkan unit_id
	if err := uc.DB.
		Where("unit_id = ?", unitID).
		Where("deleted_at IS NULL").
		Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Jika tidak ada data ditemukan, balikan 404
	if len(news) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "No news found for this unit_id",
		})
	}

	return c.JSON(fiber.Map{
		"message": "News for the selected unit retrieved successfully",
		"data":    news,
	})
}

// 🟢 GET /api/unit-news/:id
// Mengambil detail satu berita unit berdasarkan ID unik.
// Biasanya digunakan untuk membuka halaman detail pengumuman unit.
func (uc *UnitNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	// Cari data berdasarkan primary key
	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unit news found successfully",
		"data":    news,
	})
}

// 🟡 POST /api/unit-news
// Menambahkan berita atau pengumuman baru untuk unit tertentu.
// Setelah berhasil disimpan, akan memperbarui field JSON `update_news` di tabel `units` (jika tersedia).
func (uc *UnitNewsController) Create(c *fiber.Ctx) error {
	var news model.UnitNewsModel

	// Parse isi request body menjadi model
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Simpan data ke database
	if err := uc.DB.Create(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Update cache JSON news di tabel units
	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Unit news created successfully",
		"data":    news,
	})
}

// 🟠 PUT /api/unit-news/:id
// Mengupdate berita unit berdasarkan ID.
// Setelah berhasil diupdate, field JSON `update_news` di tabel units akan diperbarui.
func (uc *UnitNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	// Cek apakah data dengan ID tersebut ada
	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	// Parse request body baru
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Simpan perubahan
	if err := uc.DB.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Refresh data JSON di tabel units
	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.JSON(fiber.Map{
		"message": "Unit news updated successfully",
		"data":    news,
	})
}

// 🔴 DELETE /api/unit-news/:id
// Menghapus berita unit berdasarkan ID, kemudian memperbarui field `update_news` JSON di tabel units.
func (uc *UnitNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	// Ambil data lama
	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	// Hapus dari database
	if err := uc.DB.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Update cache JSON di tabel units
	updateUnitNewsJSON(uc.DB, news.UnitID)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Unit news with ID %v deleted successfully", news.ID),
	})
}

// ⚙️ updateUnitNewsJSON
// Helper untuk memperbarui kolom `update_news` di tabel `units`.
// Digunakan untuk menyimpan data berita dalam bentuk JSON agar bisa ditampilkan cepat oleh frontend.
func updateUnitNewsJSON(db *gorm.DB, unitID int) {
	var newsList []model.UnitNewsModel

	// Ambil semua berita berdasarkan unit_id, diurutkan dari terbaru
	if err := db.Where("unit_id = ?", unitID).Order("created_at desc").Find(&newsList).Error; err != nil {
		log.Println("[ERROR] Failed to fetch unit news for update:", err)
		return
	}

	// Ubah ke format JSON
	newsData, err := json.Marshal(newsList)
	if err != nil {
		log.Println("[ERROR] Failed to marshal unit news:", err)
		return
	}

	// Simpan ke kolom update_news di tabel units
	db.Table("units").
		Where("id = ?", unitID).
		Update("update_news", datatypes.JSON(newsData))
}
