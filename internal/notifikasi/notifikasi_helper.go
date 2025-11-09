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
