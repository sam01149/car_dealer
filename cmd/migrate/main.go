package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("==========================================")
	log.Println("🔄 Running Migration 002...")
	log.Println("==========================================")
	log.Println("")

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env tidak ditemukan, menggunakan environment variables")
	}

	// Connect to database
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("❌ DB_SOURCE tidak ditemukan")
	}

	database, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("❌ Database tidak dapat diakses: %v", err)
	}

	log.Println("✅ Terkoneksi ke database")
	log.Println("")

	// Run migrations
	migrations := []string{
		"DROP TABLE IF EXISTS transaksi_rental CASCADE",
		"ALTER TABLE mobils DROP COLUMN IF EXISTS harga_rental_per_hari",
		"ALTER TABLE mobils ADD COLUMN IF NOT EXISTS foto_url TEXT",
		"UPDATE mobils SET status = 'tersedia' WHERE status = 'dirental'",
	}

	for i, migration := range migrations {
		log.Printf("🔄 Running step %d/%d...", i+1, len(migrations))
		if _, err := database.Exec(migration); err != nil {
			log.Printf("❌ Error: %v", err)
			log.Printf("   SQL: %s", migration)
		} else {
			log.Printf("✅ Success")
		}
	}

	log.Println("")
	log.Println("==========================================")
	log.Println("✅ MIGRATION 002 SELESAI!")
	log.Println("==========================================")
	log.Println("")
	log.Println("📝 Perubahan:")
	log.Println("   • Menghapus tabel transaksi_rental")
	log.Println("   • Menghapus kolom harga_rental_per_hari dari mobils")
	log.Println("   • Menambahkan kolom foto_url ke mobils")
	log.Println("   • Update status mobil dari 'dirental' ke 'tersedia'")
	log.Println("")
}
