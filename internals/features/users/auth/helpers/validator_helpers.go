package helpers

import (
	"errors"
	"regexp"
	"strings"
)

// Validasi Register
func ValidateRegisterInput(name, email, password string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return errors.New("Nama minimal 3 karakter")
	}
	if !isValidEmail(email) {
		return errors.New("Format email tidak valid")
	}
	if len(password) < 8 {
		return errors.New("Password minimal 8 karakter")
	}
	return nil
}

// Validasi Login
func ValidateLoginInput(identifier, password string) error {
	if len(strings.TrimSpace(identifier)) < 3 {
		return errors.New("Email atau Username minimal 3 karakter")
	}
	if len(password) < 8 {
		return errors.New("Password minimal 8 karakter")
	}
	return nil
}

// Validasi Email (regex simple)
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Validasi Ganti Password
func ValidateChangePassword(oldPassword, newPassword string) error {
	if len(oldPassword) < 8 || len(newPassword) < 8 {
		return errors.New("Password minimal 8 karakter")
	}
	if oldPassword == newPassword {
		return errors.New("Password baru harus berbeda dengan password lama")
	}
	return nil
}

// Validasi Reset Password
func ValidateResetPassword(email, newPassword string) error {
	if !isValidEmail(email) {
		return errors.New("Format email tidak valid")
	}
	if len(newPassword) < 8 {
		return errors.New("Password baru minimal 8 karakter")
	}
	return nil
}

// Validasi untuk cek jawaban keamanan
func ValidateSecurityAnswerInput(email, answer string) error {
	if !isValidEmail(email) {
		return errors.New("Format email tidak valid")
	}
	if len(strings.TrimSpace(answer)) == 0 {
		return errors.New("Jawaban keamanan tidak boleh kosong")
	}
	return nil
}
