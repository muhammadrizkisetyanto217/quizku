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

	if err := uc.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Daftar unit news berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET /api/unit-news/unit/:unit_id
// Mengambil seluruh berita yang terkait dengan unit tertentu berdasarkan unit_id.
// Biasanya digunakan untuk menampilkan pengumuman terkini per unit di halaman pembelajaran.
func (uc *UnitNewsController) GetByUnitID(c *fiber.Ctx) error {
	unitID := c.Params("unit_id")
	var news []model.UnitNewsModel

	if err := uc.DB.
		Where("unit_news_unit_id = ?", unitID).
		Where("deleted_at IS NULL").
		Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if len(news) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita untuk unit ini tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita untuk unit berhasil diambil",
		"data":    news,
	})
}

// 🟢 GET /api/unit-news/:id
// Mengambil detail satu berita unit berdasarkan ID unik.
// Biasanya digunakan untuk membuka halaman detail pengumuman unit.
func (uc *UnitNewsController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	if err := uc.DB.First(&news, "unit_news_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita unit tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berita unit ditemukan",
		"data":    news,
	})
}

// 🟡 POST /api/unit-news
// Menambahkan berita atau pengumuman baru untuk unit tertentu.
// Setelah berhasil disimpan, akan memperbarui field JSON `update_news` di tabel `units` (jika tersedia).
func (uc *UnitNewsController) Create(c *fiber.Ctx) error {
	var news model.UnitNewsModel

	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Body request tidak valid",
		})
	}

	if err := uc.DB.Create(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateUnitNewsJSON(uc.DB, news.UnitNewsUnitID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berita unit berhasil ditambahkan",
		"data":    news,
	})
}

// 🟠 PUT /api/unit-news/:id
// Mengupdate berita unit berdasarkan ID.
// Setelah berhasil diupdate, field JSON `update_news` di tabel units akan diperbarui.
func (uc *UnitNewsController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var existing model.UnitNewsModel

	// Cek keberadaan data
	if err := uc.DB.First(&existing, "unit_news_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Berita unit tidak ditemukan",
		})
	}

	// Parsing update
	var updateData model.UnitNewsModel
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Body request tidak valid",
		})
	}

	updateData.UnitNewsID = existing.UnitNewsID // Jaga ID tetap

	if err := uc.DB.Save(&updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	updateUnitNewsJSON(uc.DB, updateData.UnitNewsUnitID)

	return c.JSON(fiber.Map{
		"message": "Berita unit berhasil diperbarui",
		"data":    updateData,
	})
}

// 🔴 DELETE /api/unit-news/:id
// Menghapus berita unit berdasarkan ID, kemudian memperbarui field `update_news` JSON di tabel units.
// 🔴 DELETE /api/unit-news/:id
func (uc *UnitNewsController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var news model.UnitNewsModel

	// Cari data berdasarkan ID
	if err := uc.DB.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Unit news not found",
		})
	}

	// Soft delete data
	if err := uc.DB.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Update JSON cache di tabel units
	updateUnitNewsJSON(uc.DB, news.UnitNewsUnitID) // ❗️Hati-hati terhadap kehilangan data

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Unit news with ID %v deleted successfully", news.UnitNewsID),
	})
}

// ⚙️ updateUnitNewsJSON
// Helper untuk memperbarui kolom `update_news` di tabel `units`.
// Digunakan untuk menyimpan data berita dalam bentuk JSON agar bisa ditampilkan cepat oleh frontend.
func updateUnitNewsJSON(db *gorm.DB, unitID uint) {
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
