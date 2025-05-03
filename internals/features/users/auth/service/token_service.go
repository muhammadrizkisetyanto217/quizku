package service

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"quizku/internals/configs"
	authHelper "quizku/internals/features/users/auth/helpers"
	authModel "quizku/internals/features/users/auth/models"
	authRepo "quizku/internals/features/users/auth/repository"
	userModel "quizku/internals/features/users/user/models"
	helpers "quizku/internals/helpers"
)

// ========================== REFRESH TOKEN ==========================
func RefreshToken(db *gorm.DB, c *fiber.Ctx) error {
	// 1️⃣ Ambil refresh_token dari cookie (default)
	refreshToken := c.Cookies("refresh_token")

	// 2️⃣ Atau fallback ke body JSON jika tidak ada di cookie
	if refreshToken == "" {
		var payload struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&payload); err != nil || payload.RefreshToken == "" {
			return helpers.Error(c, fiber.StatusUnauthorized, "No refresh token provided")
		}
		refreshToken = payload.RefreshToken
	}

	// 🔍 Cek token ada di database
	rt, err := authRepo.FindRefreshToken(db, refreshToken)
	if err != nil {
		return helpers.Error(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	// 🧠 Validasi isi refresh token secara manual
	claims := jwt.MapClaims{}
	parser := jwt.Parser{SkipClaimsValidation: true}
	_, err = parser.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.JWTRefreshSecret), nil
	})
	if err != nil {
		log.Println("[ERROR] Failed to parse refresh token:", err)
		return helpers.Error(c, fiber.StatusUnauthorized, "Malformed refresh token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return helpers.Error(c, fiber.StatusUnauthorized, "Refresh token missing expiration")
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return helpers.Error(c, fiber.StatusUnauthorized, "Refresh token expired")
	}

	// 🧑‍💼 Ambil user dari DB
	user, err := authRepo.FindUserByID(db, rt.UserID)
	if err != nil {
		return helpers.Error(c, fiber.StatusUnauthorized, "User not found")
	}

	// 🔁 Kembalikan access_token baru + refresh_token baru
	return issueTokens(c, db, *user)
}

// ========================== ISSUE TOKEN ==========================
func issueTokens(c *fiber.Ctx, db *gorm.DB, user userModel.UserModel) error {
	accessTokenDuration := 60 * time.Minute
	refreshTokenDuration := 7 * 24 * time.Hour

	// 🔐 Generate access token
	accessToken, accessExp, err := generateToken(user, configs.JWTSecret, accessTokenDuration)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to generate access token")
	}

	// 🔐 Generate refresh token
	refreshToken, refreshExp, err := generateToken(user, configs.JWTRefreshSecret, refreshTokenDuration)
	if err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	// 📝 Debug log durasi token
	log.Printf("[DEBUG] Access Token Exp:  %v (%s)", accessExp.Unix(), accessExp.Format(time.RFC3339))
	log.Printf("[DEBUG] Refresh Token Exp: %v (%s)", refreshExp.Unix(), refreshExp.Format(time.RFC3339))

	// 💾 Simpan refresh token ke database
	rt := authModel.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExp,
	}
	if err := authRepo.CreateRefreshToken(db, &rt); err != nil {
		return helpers.Error(c, fiber.StatusInternalServerError, "Failed to save refresh token")
	}

	// 🍪 Simpan refresh token di cookie aman
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  refreshExp,
	})

	// ✅ Kembalikan response (plus refresh token untuk debug saja, hapus di production)
	return helpers.Success(c, "Login successful", fiber.Map{
		"access_token":        accessToken,
		"refresh_token_debug": refreshToken,      // ⛔ DEBUG SAJA!
		"access_exp_unix":     accessExp.Unix(),  // debug tambahan
		"refresh_exp_unix":    refreshExp.Unix(), // debug tambahan
		"user": fiber.Map{
			"id":        user.ID,
			"user_name": user.UserName,
			"email":     user.Email,
			"role":      user.Role,
		},
	})
}

// ========================== GENERATE TOKEN ==========================
func generateToken(user userModel.UserModel, secret string, duration time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"id":        user.ID.String(),
		"user_name": user.UserName,
		"role":      user.Role,
		"exp":       exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, exp, err
}

func generateDummyPassword() string {
	hash, _ := authHelper.HashPassword("RandomDummyPassword123!")
	return hash
}

func CheckSecurityAnswer(db *gorm.DB, c *fiber.Ctx) error {
	var input struct {
		Email  string `json:"email"`
		Answer string `json:"security_answer"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helpers.Error(c, fiber.StatusBadRequest, "Invalid request format")
	}

	if err := authHelper.ValidateSecurityAnswerInput(input.Email, input.Answer); err != nil {
		return helpers.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := authRepo.FindUserByEmail(db, input.Email)
	if err != nil {
		return helpers.Error(c, fiber.StatusNotFound, "User not found")
	}

	if strings.TrimSpace(input.Answer) != strings.TrimSpace(user.SecurityAnswer) {
		return helpers.Error(c, fiber.StatusBadRequest, "Incorrect security answer")
	}

	return helpers.Success(c, "Security answer correct", fiber.Map{
		"email": user.Email,
	})
}
