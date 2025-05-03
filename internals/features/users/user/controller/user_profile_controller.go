package controller

import (
	"log"

	"quizku/internals/features/users/user/models"
	helper "quizku/internals/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersProfileController struct {
	DB *gorm.DB
}

func NewUsersProfileController(db *gorm.DB) *UsersProfileController {
	return &UsersProfileController{DB: db}
}

func (upc *UsersProfileController) GetProfiles(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all user profiles")

	var profiles []models.UsersProfileModel
	if err := upc.DB.Find(&profiles).Error; err != nil {
		log.Println("[ERROR] Failed to fetch user profiles:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to fetch user profiles")
	}

	return helper.Success(c, "User profiles fetched successfully", profiles)
}

func (upc *UsersProfileController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	log.Println("[INFO] Fetching user profile with user_id:", userID)

	var profile models.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	return helper.Success(c, "User profile fetched successfully", profile)
}

func (upc *UsersProfileController) CreateProfile(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user profile")

	// Ambil user_id dari JWT
	userID := c.Locals("user_id")
	if userID == nil {
		log.Println("[ERROR] user_id not found in context")
		return helper.Error(c, fiber.StatusUnauthorized, "Unauthorized: no user_id")
	}

	var input models.UsersProfileModel
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Invalid request format")
	}

	// Set user_id dari token ke model
	input.UserID = userID.(uuid.UUID)

	var existingProfile models.UsersProfileModel
	result := upc.DB.Where("user_id = ?", input.UserID).First(&existingProfile)

	if result.RowsAffected > 0 {
		if err := upc.DB.Model(&existingProfile).Updates(input).Error; err != nil {
			log.Println("[ERROR] Failed to update user profile:", err)
			return helper.Error(c, fiber.StatusInternalServerError, "Failed to update user profile")
		}
		log.Println("[SUCCESS] User profile updated:", input.UserID)
		return helper.Success(c, "User profile updated successfully", existingProfile)
	}

	if err := upc.DB.Create(&input).Error; err != nil {
		log.Println("[ERROR] Failed to create user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to create user profile")
	}

	log.Println("[SUCCESS] User profile created:", input.UserID)
	return helper.SuccessWithCode(c, fiber.StatusCreated, "User profile created successfully", input)
}

func (upc *UsersProfileController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	log.Println("[INFO] Updating user profile with user_id:", userID)

	var profile models.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	// Ambil data baru dari body
	if err := c.BodyParser(&profile); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Invalid request format")
	}

	profile.UserID = userID.(uuid.UUID) // Pastikan user_id tetap konsisten

	if err := upc.DB.Save(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to update user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to update user profile")
	}

	return helper.Success(c, "User profile updated successfully", profile)
}

func (upc *UsersProfileController) DeleteProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	log.Println("[INFO] Deleting user profile with user_id:", userID)

	var profile models.UsersProfileModel
	if err := upc.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	if err := upc.DB.Delete(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to delete user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to delete user profile")
	}

	return helper.Success(c, "User profile deleted successfully", nil)
}
