package controller

import (
	"log"
	"strconv"

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
	id := c.Params("id")
	log.Println("[INFO] Fetching user profile with ID:", id)

	var profile models.UsersProfileModel
	if err := upc.DB.First(&profile, id).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	return helper.Success(c, "User profile fetched successfully", profile)
}

func (upc *UsersProfileController) CreateProfile(c *fiber.Ctx) error {
	log.Println("[INFO] Creating or updating user profile")

	var input models.UsersProfileModel
	if err := c.BodyParser(&input); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Invalid request format")
	}

	if input.UserID == uuid.Nil {
		log.Println("[ERROR] Missing user_id")
		return helper.Error(c, fiber.StatusBadRequest, "user_id is required")
	}

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
	id := c.Params("id")
	log.Println("[INFO] Updating user profile with ID:", id)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("[ERROR] Invalid ID format:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Invalid ID format")
	}

	var profile models.UsersProfileModel
	if err := upc.DB.First(&profile, idInt).Error; err != nil {
		log.Println("[ERROR] User profile not found:", err)
		return helper.Error(c, fiber.StatusNotFound, "User profile not found")
	}

	if err := c.BodyParser(&profile); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return helper.Error(c, fiber.StatusBadRequest, "Invalid request format")
	}

	profile.ID = uint(idInt) // Pastikan ID tetap konsisten

	if err := upc.DB.Save(&profile).Error; err != nil {
		log.Println("[ERROR] Failed to update user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to update user profile")
	}

	return helper.Success(c, "User profile updated successfully", profile)
}

func (upc *UsersProfileController) DeleteProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("[INFO] Deleting user profile with ID:", id)

	if err := upc.DB.Delete(&models.UsersProfileModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete user profile:", err)
		return helper.Error(c, fiber.StatusInternalServerError, "Failed to delete user profile")
	}

	return helper.Success(c, "User profile deleted successfully", nil)
}
