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
	log.Println("🧹 RESET DATA CLIENT")
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

	// Get dealer user ID
	var dealerID string
	err = database.QueryRow("SELECT id FROM users WHERE email = 'dealer@carapp.com'").Scan(&dealerID)
	if err != nil {
		log.Fatalf("❌ Gagal mendapatkan dealer ID: %v", err)
	}
	log.Printf("📌 Dealer ID: %s", dealerID)
	log.Println("")

	// Reset operations
	operations := []struct {
		name string
		sql  string
	}{
		{
			name: "Hapus semua notifikasi dari user client",
			sql:  "DELETE FROM notifikasi WHERE user_id != $1",
		},
		{
			name: "Hapus semua transaksi jual",
			sql:  "DELETE FROM transaksi_jual",
		},
		{
			name: "Hapus semua mobil dari user client",
			sql:  "DELETE FROM mobils WHERE owner_id != $1",
		},
		{
			name: "Hapus semua user client (bukan dealer)",
			sql:  "DELETE FROM users WHERE id != $1",
		},
		{
			name: "Reset mobil dealer ke status tersedia",
			sql:  "UPDATE mobils SET status = 'tersedia' WHERE owner_id = $1",
		},
	}

	for i, op := range operations {
		log.Printf("🔄 Step %d/%d: %s...", i+1, len(operations), op.name)
		
		var result sql.Result
		var err error
		
		if op.sql == "DELETE FROM transaksi_jual" {
			// Transaksi jual tidak perlu parameter dealer ID
			result, err = database.Exec(op.sql)
		} else {
			result, err = database.Exec(op.sql, dealerID)
		}
		
		if err != nil {
			log.Printf("❌ Error: %v", err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			log.Printf("✅ Success - %d rows affected", rowsAffected)
		}
	}

	log.Println("")
	log.Println("==========================================")
	log.Println("✅ RESET SELESAI!")
	log.Println("==========================================")
	log.Println("")
	log.Println("📝 Yang tersisa di database:")
	log.Println("   • User dealer (dealer@carapp.com)")
	log.Println("   • 42 mobil dari dealer (status: tersedia)")
	log.Println("   • Notifikasi dealer")
	log.Println("")
	log.Println("🎯 Aplikasi siap untuk testing dari awal!")
	log.Println("")
}
