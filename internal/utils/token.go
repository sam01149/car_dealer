package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims kustom untuk data di dalam token
type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// getJwtSecret mengambil secret key dari .env
func getJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}
	return []byte(secret)
}

// GenerateToken membuat JWT baru untuk user.
func GenerateToken(userID, email, role string) (string, error) {
	// Set claims
	claims := &JwtCustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token valid 72 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret
	t, err := token.SignedString(getJwtSecret())
	if err != nil {
		return "", err
	}

	return t, nil
}

// ValidateToken memverifikasi token.
// (Kita akan gunakan ini nanti di middleware)
func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// PENJELASAN FILE token.go:
// File ini menangani pembuatan dan validasi JWT (JSON Web Token) untuk autentikasi
//
// Struct JwtCustomClaims:
// - Menyimpan data user dalam token: UserID, Email, Role
// - Juga berisi ExpiresAt dan IssuedAt dari jwt.RegisteredClaims
// - Data ini bisa diakses setelah token divalidasi
//
// Fungsi getJwtSecret:
// - Mengambil secret key dari environment variable JWT_SECRET_KEY
// - Secret ini digunakan untuk sign dan verify token
// - Harus dijaga kerahasiaannya (jangan commit ke git!)
//
// Fungsi GenerateToken:
// - Dipanggil setelah login/register berhasil
// - Membuat JWT token dengan masa berlaku 72 jam
// - Token berisi user_id, email, role yang ter-encrypt
// - Return token string untuk dikirim ke client
//
// Fungsi ValidateToken:
// - Dipanggil di middleware untuk setiap request
// - Parse token dan verifikasi signature dengan secret key
// - Cek apakah token sudah expired
// - Return claims (user info) jika valid, error jika tidak
//
// Keamanan:
// - Token di-sign dengan JWT_SECRET_KEY (prevent tampering)
// - Token expire otomatis setelah 72 jam
// - Payload bisa di-decode tapi tidak bisa diubah tanpa secret
// - Jangan simpan data sensitif (password, credit card) di token
