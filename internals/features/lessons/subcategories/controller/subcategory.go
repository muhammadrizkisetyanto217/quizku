package controller

import (
	"fmt"
	"log"

	categoryModel "quizku/internals/features/lessons/categories/model"

	"quizku/internals/features/lessons/subcategories/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubcategoryController struct {
	DB *gorm.DB
}

func NewSubcategoryController(db *gorm.DB) *SubcategoryController {
	return &SubcategoryController{DB: db}
}

// 🟢 GET ALL SUBCATEGORIES: Ambil seluruh data subkategori
func (sc *SubcategoryController) GetSubcategories(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all subcategories")
	var subcategories []model.SubcategoryModel

	// 🔍 Query semua subkategori
	if err := sc.DB.Find(&subcategories).Error; err != nil {
		log.Println("[ERROR] Failed to fetch subcategories:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}

	// ✅ Kirim data subkategori
	log.Printf("[SUCCESS] Retrieved %d subcategories\n", len(subcategories))
	return c.JSON(fiber.Map{
		"message": "All subcategories fetched successfully",
		"total":   len(subcategories),
		"data":    subcategories,
	})
}

// 🟢 GET SUBCATEGORY BY ID: Ambil data subkategori berdasarkan ID
func (sc *SubcategoryController) GetSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Fetching subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Query berdasarkan ID
	if err := sc.DB.Where("subcategory_id = ?", id).First(&subcategory).Error; err != nil {
		log.Println("[ERROR] Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}

	// ✅ Kirim data
	log.Printf("[SUCCESS] Subcategory retrieved: ID=%d, Name=%s\n", subcategory.SubcategoryID, subcategory.SubcategoryName)
	return c.JSON(fiber.Map{
		"message": "Subcategory fetched successfully",
		"data":    subcategory,
	})
}

// 🟢 GET SUBCATEGORIES BY CATEGORY ID: Ambil data subkategori berdasarkan subcategory_category_id
func (sc *SubcategoryController) GetSubcategoriesByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	log.Printf("[INFO] Fetching subcategories with category ID: %s\n", categoryID)

	var subcategories []model.SubcategoryModel

	// 🔍 Query subkategori berdasarkan subcategory_category_id
	if err := sc.DB.Where("subcategory_category_id = ?", categoryID).Find(&subcategories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch subcategories for category ID %s: %v\n", categoryID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch subcategories"})
	}

	// ✅ Kirim hasil
	log.Printf("[SUCCESS] Retrieved %d subcategories for category ID %s\n", len(subcategories), categoryID)
	return c.JSON(fiber.Map{
		"message": "Subcategories fetched successfully by category",
		"total":   len(subcategories),
		"data":    subcategories,
	})
}

// 🟢 CREATE SUBCATEGORY: Tambah satu atau banyak subkategori dengan validasi
func (sc *SubcategoryController) CreateSubcategory(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk membuat subkategori")

	var single model.SubcategoryModel
	var multiple []model.SubcategoryModel

	// 🧠 Coba parsing sebagai array
	if err := c.BodyParser(&multiple); err == nil && len(multiple) > 0 {
		log.Printf("[DEBUG] Parsed sebagai array: %d subcategories\n", len(multiple))

		// ✅ Validasi setiap subkategori dalam array
		for i, item := range multiple {
			if item.SubcategoryName == "" || item.SubcategoryCategoryID == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "Field subcategory_name dan subcategory_category_id wajib diisi",
					"index": i,
				})
			}
			var exists int64
			if err := sc.DB.Table("categories").
				Where("category_id = ?", item.SubcategoryCategoryID).
				Count(&exists).Error; err != nil || exists == 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "subcategory_category_id tidak valid (kategori tidak ditemukan)",
					"index": i,
				})
			}
		}

		// 💾 Simpan seluruh array
		if err := sc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Gagal insert batch subcategories: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan subcategories"})
		}

		log.Printf("[SUCCESS] %d subcategories berhasil dibuat", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Subcategories berhasil dibuat",
			"data":    multiple,
		})
	}

	// 🔁 Parsing sebagai objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Gagal parsing request body single: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	log.Printf("[DEBUG] Subcategory tunggal yang dikirim: %+v", single)

	// ✅ Validasi tunggal
	if single.SubcategoryName == "" || single.SubcategoryCategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Field subcategory_name dan subcategory_category_id wajib diisi",
		})
	}

	var exists int64
	if err := sc.DB.Table("categories").
		Where("category_id = ?", single.SubcategoryCategoryID).
		Count(&exists).Error; err != nil || exists == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "subcategory_category_id tidak valid (kategori tidak ditemukan)",
		})
	}

	// 💾 Simpan subkategori tunggal
	if err := sc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Gagal insert subcategory tunggal: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan subcategory"})
	}

	log.Printf("[SUCCESS] Subcategory berhasil dibuat: ID=%d, Nama=%s", single.SubcategoryID, single.SubcategoryName)
	return c.Status(201).JSON(fiber.Map{
		"message": "Subcategory berhasil dibuat",
		"data":    single,
	})
}

// 🟢 UPDATE SUBCATEGORY: Perbarui subkategori berdasarkan ID
func (sc *SubcategoryController) UpdateSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Updating subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Cari data lama berdasarkan subcategory_id
	if err := sc.DB.Where("subcategory_id = ?", id).First(&subcategory).Error; err != nil {
		log.Println("[ERROR] Subcategory not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Subcategory not found"})
	}

	// 📝 Parse data baru dari request body
	if err := c.BodyParser(&subcategory); err != nil {
		log.Println("[ERROR] Body request tidak valid:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	// 💾 Simpan perubahan
	if err := sc.DB.Save(&subcategory).Error; err != nil {
		log.Println("[ERROR] Gagal update subcategory:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update subcategory"})
	}

	log.Printf("[SUCCESS] Subcategory updated: ID=%d, Name=%s", subcategory.SubcategoryID, subcategory.SubcategoryName)
	return c.JSON(fiber.Map{
		"message": "Subcategory berhasil diperbarui",
		"data":    subcategory,
	})
}

// 🟢 DELETE SUBCATEGORY: Hapus subkategori berdasarkan ID
func (sc *SubcategoryController) DeleteSubcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting subcategory with ID:", id)

	var subcategory model.SubcategoryModel

	// 🔍 Cari subkategori berdasarkan ID
	if err := sc.DB.Where("subcategory_id = ?", id).First(&subcategory).Error; err != nil {
		log.Println("[ERROR] Subcategory tidak ditemukan:", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Subcategory tidak ditemukan",
		})
	}

	// 🗑️ Hapus (soft delete menggunakan DeletedAt)
	if err := sc.DB.Delete(&subcategory).Error; err != nil {
		log.Println("[ERROR] Gagal hapus subcategory:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghapus subcategory",
		})
	}

	log.Printf("[SUCCESS] Subcategory dengan ID %s berhasil dihapus", id)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Subcategory dengan ID %s berhasil dihapus", id),
	})
}
// 🟢 GET CATEGORY WITH SUBCATEGORY AND THEMES: Ambil data lengkap kategori, subkategori, dan themes berdasarkan difficulty_id
func (sc *SubcategoryController) GetCategoryWithSubcategoryAndThemes(c *fiber.Ctx) error {
	difficultyID := c.Params("difficulty_id")
	log.Printf("[INFO] Fetching category, subcategory, and themes for difficulty ID: %s\n", difficultyID)

	var categories []categoryModel.CategoryModel

	// 🔍 Ambil semua kategori yang memiliki difficulty_id sesuai
	// dan preload semua subcategory aktif + themes di dalamnya
	if err := sc.DB.
		Where("category_difficulty_id = ?", difficultyID).
		Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("subcategory_status = ?", "active").
				Preload("ThemesOrLevels", func(db2 *gorm.DB) *gorm.DB {
					return db2.Order("themes_or_level_id ASC")
				})
		}).
		Order("category_id ASC").
		Find(&categories).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch categories: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data kategori lengkap",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d categories with subcategories and themes for difficulty ID %s\n", len(categories), difficultyID)
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data kategori lengkap",
		"data":    categories,
	})
}
