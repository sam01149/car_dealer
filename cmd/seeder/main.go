package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"carapp.com/m/internal/db"
	"carapp.com/m/internal/utils"
	"github.com/joho/godotenv"
)

// Struct untuk parsing JSON dari Marketcheck
type MarketcheckResponse struct {
	Listings []struct {
		ID      string `json:"id"`
		Heading string `json:"heading"` // "2019 Honda Civic..."
		VIN     string `json:"vin"`
		Build   struct {
			Make      string `json:"make"`
			Model     string `json:"model"`
			Year      int    `json:"year"`
			BodyType  string `json:"body_type"`
			DriveType string `json:"drivetrain"`
			FuelType  string `json:"fuel_type"`
		} `json:"build"`
		Media struct {
			PhotoLinks []string `json:"photo_links"`
		} `json:"media"`
		Price         float64 `json:"price"`
		InventoryType string  `json:"inventory_type"` // "used", "new"
		Dealer        struct {
			City  string `json:"city"`
			State string `json:"state"`
		} `json:"dealer"`
		Miles    int    `json:"miles"`
		Exterior string `json:"exterior_color"`
		Interior string `json:"interior_color"`
	} `json:"listings"`
	NumFound int `json:"num_found"`
}

func main() {
	log.Println("===========================================")
	log.Println("🚗 Memulai Seeder Inventaris Dealer...")
	log.Println("===========================================")
	log.Println("")

	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		// Coba dari cmd/seeder
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan env variables...")
		}
	}

	// 2. Koneksi DB
	log.Println("📡 Menghubungkan ke database...")
	dbConn := db.ConnectDB()
	defer dbConn.Close()
	log.Println("✅ Berhasil terhubung ke DB.")
	log.Println("")

	// 3. Dapatkan/Buat User "Dealer"
	dealerUserID, err := getOrCreateDealerUser(dbConn)
	if err != nil {
		log.Fatalf("❌ Gagal membuat user Dealer: %v", err)
	}
	log.Printf("✅ Menggunakan ID Dealer: %s", dealerUserID)
	log.Println("")

	// RESET MOBIL LAMA (Hapus semua mobil dealer)
	log.Println("🗑️  Menghapus mobil lama dari dealer...")

	// Hapus transaksi terlebih dahulu untuk menghindari foreign key constraint
	_, err = dbConn.Exec("DELETE FROM transaksi_jual WHERE penjual_id = $1 OR pembeli_id = $1", dealerUserID)
	if err != nil {
		log.Printf("⚠️  Warning: Gagal menghapus transaksi lama: %v", err)
	}

	// Sekarang hapus mobil
	_, err = dbConn.Exec("DELETE FROM mobils WHERE owner_id = $1", dealerUserID)
	if err != nil {
		log.Fatalf("❌ Gagal menghapus mobil lama: %v", err)
	}
	log.Println("✅ Mobil lama berhasil dihapus.")
	log.Println("")

	// 4. Ambil API Key
	apiKey := os.Getenv("MARKETCHECK_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ MARKETCHECK_API_KEY tidak ditemukan di .env")
	}

	// 5. Panggil Marketcheck API
	log.Println("🌐 Memanggil Marketcheck API untuk mengambil stok...")

	// Query untuk mobil bekas di berbagai lokasi populer
	// Kita ambil 50 mobil untuk variasi yang lebih baik
	url := fmt.Sprintf("https://api.marketcheck.com/v2/search/car/active?api_key=%s&car_type=used&rows=50&start=0", apiKey)

	log.Printf("📞 Request URL: %s", strings.Replace(url, apiKey, "***", -1))

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("❌ Gagal memanggil API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("❌ API mengembalikan status %d: %s", resp.StatusCode, resp.Status)
	}

	var apiResponse MarketcheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Fatalf("❌ Gagal parsing JSON: %v", err)
	}

	if len(apiResponse.Listings) == 0 {
		log.Fatal("❌ API tidak mengembalikan mobil, cek API Key/Query Anda.")
	}

	log.Printf("✅ Sukses mengambil %d mobil dari Marketcheck (dari total %d).", len(apiResponse.Listings), apiResponse.NumFound)
	log.Println("")

	// 6. Simpan ke Database
	log.Println("💾 Menyimpan mobil ke database...")
	ctx := context.Background()
	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Gagal memulai transaksi DB: %v", err)
	}
	defer tx.Rollback() // Rollback jika ada error

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO mobils (
			owner_id, merk, model, tahun, kondisi, deskripsi, 
			harga_jual, foto_url, lokasi, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		log.Fatalf("❌ Gagal menyiapkan statement DB: %v", err)
	}
	defer stmt.Close()

	count := 0
	skipped := 0

	for i, mobil := range apiResponse.Listings {
		// Skip mobil tanpa harga atau data tidak lengkap
		if mobil.Price <= 0 || mobil.Build.Make == "" || mobil.Build.Model == "" || mobil.Build.Year == 0 {
			skipped++
			continue
		}

		// Konversi kondisi ke bahasa Indonesia
		kondisi := "bekas"
		if mobil.InventoryType == "new" {
			kondisi = "baru"
		}

		// Buat deskripsi yang informatif
		deskripsi := fmt.Sprintf("%s - %d miles. %s exterior, %s interior. %s. VIN: %s",
			mobil.Heading,
			mobil.Miles,
			mobil.Exterior,
			mobil.Interior,
			mobil.Build.BodyType,
			mobil.VIN,
		)

		// KONVERSI USD KE IDR
		// Kurs: 1 USD = 15.800 IDR (approximate)
		const USD_TO_IDR = 15800.0
		hargaJualIDR := mobil.Price * USD_TO_IDR

		// Ambil foto pertama jika ada
		fotoUrl := ""
		if len(mobil.Media.PhotoLinks) > 0 {
			fotoUrl = mobil.Media.PhotoLinks[0]
		}

		// Lokasi
		lokasi := "Jakarta, Indonesia" // Default
		if mobil.Dealer.City != "" && mobil.Dealer.State != "" {
			lokasi = fmt.Sprintf("%s, %s", mobil.Dealer.City, mobil.Dealer.State)
		}

		_, err := stmt.ExecContext(ctx,
			dealerUserID,
			mobil.Build.Make,
			mobil.Build.Model,
			mobil.Build.Year,
			kondisi,
			deskripsi,
			hargaJualIDR,
			fotoUrl,
			lokasi,
			"tersedia",
		)
		if err != nil {
			log.Printf("⚠️  Gagal menyimpan mobil #%d (%s): %v", i+1, mobil.Heading, err)
			skipped++
			continue
		}
		count++

		// Progress indicator
		if (i+1)%10 == 0 {
			log.Printf("   📝 Progress: %d/%d mobil diproses...", i+1, len(apiResponse.Listings))
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("❌ Gagal commit ke DB: %v", err)
	}

	log.Println("")
	log.Println("===========================================")
	log.Printf("✅ SELESAI! Berhasil menyimpan %d mobil baru", count)
	if skipped > 0 {
		log.Printf("⚠️  %d mobil dilewati (data tidak lengkap)", skipped)
	}
	log.Println("===========================================")
	log.Println("")
	log.Println("🎉 Database Anda sekarang terisi dengan inventaris mobil real!")
	log.Println("💡 Jalankan frontend Next.js dan lihat hasilnya di http://localhost:3000")
}

// getOrCreateDealerUser adalah helper untuk membuat akun "Dealer"
func getOrCreateDealerUser(db *sql.DB) (string, error) {
	ctx := context.Background()
	dealerEmail := "dealer@carapp.com"

	var userID string
	// Cek apakah user dealer sudah ada
	err := db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", dealerEmail).Scan(&userID)
	if err == nil {
		// User sudah ada, kembalikan ID-nya
		log.Println("   ℹ️  User 'dealer@carapp.com' sudah ada, menggunakan yang existing.")
		return userID, nil
	}

	// Jika tidak ada (sql.ErrNoRows), buat baru
	log.Println("   📝 User 'dealer@carapp.com' tidak ditemukan, membuat baru...")

	// Hash password dummy
	hashedPassword, err := utils.HashPassword("dealer123")
	if err != nil {
		return "", err
	}

	query := `INSERT INTO users (name, email, password_hash, phone, role)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id`

	err = db.QueryRowContext(ctx, query, "Dealer Resmi", dealerEmail, hashedPassword, "+62812345678", "admin").Scan(&userID)
	if err != nil {
		return "", err
	}

	log.Println("   ✅ Berhasil membuat user 'Dealer Resmi' (Email: dealer@carapp.com, Password: dealer123, Phone: +62812345678)")
	return userID, nil
}

// PENJELASAN FILE cmd/seeder/main.go:
// File ini untuk mengisi database dengan data mobil dummy dari Marketcheck API
//
// Flow Seeder:
// 1. Load .env untuk konfigurasi (DB_SOURCE, MARKETCHECK_API_KEY)
// 2. Koneksi ke database PostgreSQL
// 3. Buat atau cari user "Dealer" (email: dealer@carapp.com)
// 4. Hapus mobil lama milik dealer (RESET data)
// 5. Call Marketcheck API untuk ambil 50 mobil bekas
// 6. Loop setiap mobil dan insert ke database
// 7. Convert harga dari USD ke IDR (kurs 1 USD = 15,800 IDR)
// 8. Set harga rental = 0.5% dari harga jual per hari (min 100k)
// 9. Commit transaction ke database
//
// Fungsi getOrCreateDealerUser:
// - Cek apakah user dealer sudah ada di database
// - Jika sudah ada, return ID-nya
// - Jika belum ada, buat user baru dengan:
//   * Name: "Dealer Resmi"
//   * Email: "dealer@carapp.com"
//   * Password: "dealer123" (di-hash dengan bcrypt)
//   * Role: "admin"
//
// Struct MarketcheckResponse:
// - Untuk parsing JSON dari Marketcheck API
// - Berisi: merek, model, tahun, harga, jarak tempuh, warna, lokasi, dll
//
// Database:
// - Seeder menggunakan transaction untuk insert batch (lebih efisien)
// - ON CONFLICT DO NOTHING untuk skip mobil duplicate
// - Progress indicator setiap 10 mobil
//
// Catatan:
// - API Key Marketcheck harus valid (dari https://www.marketcheck.com)
// - Seeder bisa dijalankan berulang kali (akan reset data dealer)
// - Untuk testing/development, tidak untuk production

// PENJELASAN FILE cmd/seeder/main.go:
// File ini untuk mengisi database dengan data mobil dummy dari Marketcheck API
//
// Flow Seeder:
// 1. Load .env untuk konfigurasi (DB_SOURCE, MARKETCHECK_API_KEY)
// 2. Koneksi ke database PostgreSQL
// 3. Buat atau cari user "Dealer" (email: dealer@carapp.com)
// 4. Hapus mobil lama milik dealer (RESET data)
// 5. Call Marketcheck API untuk ambil 50 mobil bekas
// 6. Loop setiap mobil dan insert ke database
// 7. Convert harga dari USD ke IDR (kurs 1 USD = 15,800 IDR)
// 8. Set harga rental = 0.5% dari harga jual per hari (min 100k)
// 9. Commit transaction ke database
//
// Fungsi getOrCreateDealerUser:
// - Cek apakah user dealer sudah ada di database
// - Jika sudah ada, return ID-nya
// - Jika belum ada, buat user baru dengan:
//   * Name: "Dealer Resmi"
//   * Email: "dealer@carapp.com"
//   * Password: "dealer123" (di-hash dengan bcrypt)
//   * Role: "admin"
//
// Struct MarketcheckResponse:
// - Untuk parsing JSON dari Marketcheck API
// - Berisi: merek, model, tahun, harga, jarak tempuh, warna, lokasi, dll
//
// Database:
// - Seeder menggunakan transaction untuk insert batch (lebih efisien)
// - ON CONFLICT DO NOTHING untuk skip mobil duplicate
// - Progress indicator setiap 10 mobil
//
// Catatan:
// - API Key Marketcheck harus valid (dari https://www.marketcheck.com)
// - Seeder bisa dijalankan berulang kali (akan reset data dealer)
// - Untuk testing/development, tidak untuk production
