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