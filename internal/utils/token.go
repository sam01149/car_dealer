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
