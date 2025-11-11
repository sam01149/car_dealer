package notifikasi

import (
	"database/sql"
	"log"
	"time"

	"carapp.com/m/internal/auth" // Sesuaikan dengan nama modul Anda
	pb "carapp.com/m/proto"      // Sesuaikan dengan nama modul Anda
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NotifikasiServiceServer adalah implementasi dari pb.NotifikasiServiceServer
type NotifikasiServiceServer struct {
	pb.UnimplementedNotifikasiServiceServer
	DB *sql.DB
}

// NewNotifikasiService membuat instance baru
func NewNotifikasiService(db *sql.DB) *NotifikasiServiceServer {
	return &NotifikasiServiceServer{DB: db}
}

// GetNotifications adalah streaming RPC (Fitur 6)
func (s *NotifikasiServiceServer) GetNotifications(req *pb.GetNotificationsRequest, stream pb.NotifikasiService_GetNotificationsServer) error {
	ctx := stream.Context()

	// 1. Dapatkan UserID dari token
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Tidak dapat mengambil UserID dari token")
	}

	log.Printf("NotifikasiService: Memulai stream notifikasi untuk UserID %s", userID)

	// 2. Query notifikasi terbaru (misal, 50 terakhir)
	// Implementasi sederhana: kirim semua notifikasi historis dan tutup stream.
	query := `
		SELECT id, user_id, tipe, pesan, priority, read_at, created_at
		FROM notifikasi
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50 -- Batasi 50 notifikasi terbaru
	`
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Gagal query notifikasi: %v", err)
		return status.Errorf(codes.Internal, "Gagal mengambil notifikasi")
	}
	defer rows.Close()

	// 3. Kirim notifikasi satu per satu ke client via stream
	for rows.Next() {
		var notif pb.Notifikasi
		var createdAt time.Time
		var readAt sql.NullTime // Gunakan NullTime untuk kolom 'read_at'

		err := rows.Scan(
			&notif.Id,
			&notif.UserId,
			&notif.Tipe,
			&notif.Pesan,
			&notif.Priority,
			&readAt,
			&createdAt,
		)
		if err != nil {
			log.Printf("Gagal scan notifikasi: %v", err)
			continue // Lewati notifikasi yg error
		}

		// Konversi ke format Protobuf Timestamp
		notif.CreatedAt = timestamppb.New(createdAt)
		if readAt.Valid {
			notif.ReadAt = timestamppb.New(readAt.Time)
		}

		// Kirim ke stream
		if err := stream.Send(&notif); err != nil {
			log.Printf("Gagal mengirim notifikasi ke stream: %v", err)
			// Kemungkinan client sudah disconnect
			return status.Errorf(codes.Aborted, "Stream client ditutup")
		}
	}

	log.Printf("Selesai streaming notifikasi historis untuk UserID %s", userID)

	// Catatan: Untuk notifikasi "real-time" (opsional), kita bisa menambahkan
	// time.Ticker di sini untuk terus mengecek DB setiap 10 detik.
	// Tapi untuk MVP, mengirim data historis sudah cukup.

	return nil
}

// PENJELASAN FILE notifikasi_service.go:
// File ini menangani streaming notifikasi ke client (server-side streaming)
//
// Fungsi GetNotifications (Server Streaming RPC):
// - Mengambil user_id dari context (diisi oleh auth middleware)
// - Query 50 notifikasi terbaru user dari database
// - Kirim notifikasi satu per satu ke client via stream
// - Client akan menerima notifikasi secara real-time
//
// Flow:
// 1. Client buka stream connection
// 2. Server validate token di middleware
// 3. Server query notifikasi dari DB (ORDER BY created_at DESC LIMIT 50)
// 4. Loop setiap notifikasi dan kirim via stream.Send()
// 5. Jika error atau selesai, tutup stream
//
// Database:
// - read_at bisa NULL (notifikasi belum dibaca)
// - Menggunakan sql.NullTime untuk handle NULL value

// PENJELASAN FILE notifikasi_service.go:
// File ini menangani streaming notifikasi ke client (server-side streaming)
//
// Fungsi GetNotifications (Server Streaming RPC):
// - Mengambil user_id dari context (diisi oleh auth middleware)
// - Query 50 notifikasi terbaru user dari database
// - Kirim notifikasi satu per satu ke client via stream
// - Client akan menerima notifikasi secara real-time
//
// Flow:
// 1. Client buka stream connection
// 2. Server validate token di middleware
// 3. Server query notifikasi dari DB (ORDER BY created_at DESC LIMIT 50)
// 4. Loop setiap notifikasi dan kirim via stream.Send()
// 5. Jika error atau selesai, tutup stream

// PENJELASAN FILE notifikasi_service.go:
// File ini menangani streaming notifikasi ke client (server-side streaming)
//
// Fungsi GetNotifications (Server Streaming RPC):
// - Mengambil user_id dari context (diisi oleh auth middleware)
// - Query 50 notifikasi terbaru user dari database
// - Kirim notifikasi satu per satu ke client via stream
// - Client akan menerima notifikasi secara real-time
//
// Flow:
// 1. Client buka stream connection
// 2. Server validate token di middleware
// 3. Server query notifikasi dari DB (ORDER BY created_at DESC LIMIT 50)
// 4. Loop setiap notifikasi dan kirim via stream.Send()
// 5. Jika error atau selesai, tutup stream
//
// Database:
// - read_at bisa NULL (notifikasi belum dibaca)
// - Menggunakan sql.NullTime untuk handle NULL value

// PENJELASAN FILE notifikasi_service.go:
// File ini menangani streaming notifikasi ke client (server-side streaming)
//
// Fungsi GetNotifications (Server Streaming RPC):
// - Mengambil user_id dari context (diisi oleh auth middleware)
// - Query 50 notifikasi terbaru user dari database
// - Kirim notifikasi satu per satu ke client via stream
// - Client akan menerima notifikasi secara real-time
//
// Flow:
// 1. Client buka stream connection
// 2. Server validate token di middleware
// 3. Server query notifikasi dari DB (ORDER BY created_at DESC LIMIT 50)
// 4. Loop setiap notifikasi dan kirim via stream.Send()
// 5. Jika error atau selesai, tutup stream
//
// Database:
// - read_at bisa NULL (notifikasi belum dibaca)
// - Menggunakan sql.NullTime untuk handle NULL value
