package notifikasi

import (
	"context"
	"database/sql"
	"log"
)

// CreateNotification adalah helper internal untuk membuat notifikasi.
// Kita akan memanggil ini sebagai goroutine (go ...) agar tidak memblokir
// respons utama ke pengguna.
func CreateNotification(db *sql.DB, ctx context.Context, userID, tipe, pesan string) {
	query := `
		INSERT INTO notifikasi (user_id, tipe, pesan, priority)
		VALUES ($1, $2, $3, 'normal')
	`
	_, err := db.ExecContext(ctx, query, userID, tipe, pesan)
	if err != nil {
		// Dalam produksi, ini harusnya di-log ke sistem monitoring
		log.Printf("GAGAL MEMBUAT NOTIFIKASI untuk UserID %s: %v", userID, err)
	} else {
		log.Printf("Notifikasi dibuat untuk UserID %s, Tipe: %s", userID, tipe)
	}
}

// PENJELASAN FILE notifikasi_helper.go:
// File ini menyediakan helper function untuk membuat notifikasi
//
// Fungsi CreateNotification:
// - Helper internal yang dipanggil dari service lain (transaksi, mobil)
// - Dipanggil sebagai goroutine (go CreateNotification...) agar tidak blocking
// - Insert notifikasi baru ke tabel 'notifikasi' di database
// - Parameter: userID (penerima), tipe (jual/beli/rental), pesan (deskripsi)
// - Priority default: 'normal'
// - Jika gagal, hanya log error (tidak mengganggu flow utama)
//
// Use case:
// - Setelah user beli mobil -> buat notifikasi untuk pembeli dan penjual
// - Setelah user rental mobil -> buat notifikasi untuk penyewa dan pemilik
// - Setelah user posting mobil -> buat notifikasi untuk penjual

// PENJELASAN FILE notifikasi_helper.go:
// File ini menyediakan helper function untuk membuat notifikasi
//
// Fungsi CreateNotification:
// - Helper internal yang dipanggil dari service lain (transaksi, mobil)
// - Dipanggil sebagai goroutine (go CreateNotification...) agar tidak blocking
// - Insert notifikasi baru ke tabel 'notifikasi' di database
// - Parameter: userID (penerima), tipe (jual/beli/rental), pesan (deskripsi)
// - Priority default: 'normal'
// - Jika gagal, hanya log error (tidak mengganggu flow utama)
//
// Use case:
// - Setelah user beli mobil -> buat notifikasi untuk pembeli dan penjual
// - Setelah user rental mobil -> buat notifikasi untuk penyewa dan pemilik
// - Setelah user posting mobil -> buat notifikasi untuk penjual

// PENJELASAN FILE notifikasi_helper.go:
// File ini menyediakan helper function untuk membuat notifikasi
//
// Fungsi CreateNotification:
// - Helper internal yang dipanggil dari service lain (transaksi, mobil)
// - Dipanggil sebagai goroutine (go CreateNotification...) agar tidak blocking
// - Insert notifikasi baru ke tabel 'notifikasi' di database
// - Parameter: userID (penerima), tipe (jual/beli/rental), pesan (deskripsi)
// - Priority default: 'normal'
// - Jika gagal, hanya log error (tidak mengganggu flow utama)
//
// Use case:
// - Setelah user beli mobil -> buat notifikasi untuk pembeli dan penjual
// - Setelah user rental mobil -> buat notifikasi untuk penyewa dan pemilik
// - Setelah user posting mobil -> buat notifikasi untuk penjual
