package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword membuat hash bcrypt dari password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash membandingkan password dengan hash-nya.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // true jika cocok, false jika tidak
}
