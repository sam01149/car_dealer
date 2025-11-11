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

// PENJELASAN FILE password.go:
// File ini menyediakan fungsi untuk hashing dan validasi password dengan bcrypt
//
// Fungsi HashPassword:
// - Menerima plain text password (contoh: "password123")
// - Menggunakan bcrypt.GenerateFromPassword dengan cost default (10)
// - Bcrypt otomatis menambahkan random salt untuk setiap password
// - Return hash string (contoh: "$2a$10$abc...xyz")
// - Hash ini yang disimpan di database, bukan plain text
//
// Fungsi CheckPasswordHash:
// - Menerima plain text password dan hash dari database
// - Menggunakan bcrypt.CompareHashAndPassword untuk validasi
// - Return true jika password cocok, false jika tidak
// - Digunakan saat login untuk verifikasi password user
//
// Keamanan:
// - Bcrypt sengaja lambat (~100-300ms) untuk mencegah brute-force attack
// - Setiap password menghasilkan hash berbeda karena random salt
// - Hash tidak bisa di-decrypt kembali ke plain text (one-way hashing)
