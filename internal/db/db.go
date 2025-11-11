package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

func ConnectDB() *sql.DB {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable is not set")
	}

	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal mem-ping database: %v", err)
	}

	log.Println("Berhasil terhubung ke database PostgreSQL!")
	return db
}

// PENJELASAN FILE db.go:
// File ini berfungsi untuk membuat koneksi ke database PostgreSQL
//
// Fungsi ConnectDB:
// - Membaca connection string dari environment variable DB_SOURCE
// - Format: "postgresql://user:password@host:port/dbname?sslmode=disable"
// - Membuka koneksi dengan driver "postgres" (github.com/lib/pq)
// - Melakukan Ping untuk memastikan koneksi berhasil
// - Return pointer *sql.DB yang siap digunakan oleh service lain
//
// Catatan:
// - DB_SOURCE harus diset di file .env
// - Koneksi ini dibuat sekali di main.go dan dibagikan ke semua service
// - sql.DB sudah mengelola connection pooling secara otomatis
