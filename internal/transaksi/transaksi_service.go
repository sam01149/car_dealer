package transaksi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"carapp.com/m/internal/auth" // Sesuaikan dengan modul Anda
	"carapp.com/m/internal/notifikasi"
	pb "carapp.com/m/proto" // Sesuaikan dengan modul Anda
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransaksiServiceServer adalah implementasi dari pb.TransaksiServiceServer
type TransaksiServiceServer struct {
	pb.UnimplementedTransaksiServiceServer
	DB *sql.DB
}

// NewTransaksiService membuat instance baru
func NewTransaksiService(db *sql.DB) *TransaksiServiceServer {
	return &TransaksiServiceServer{DB: db}
}

// BuyMobil menangani logika pembelian mobil (Fitur 3)
func (s *TransaksiServiceServer) BuyMobil(ctx context.Context, req *pb.BuyMobilRequest) (*pb.TransaksiJualResponse, error) {
	// 1. Dapatkan ID pembeli dari token
	pembeliID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak dapat mengambil UserID dari token")
	}

	if req.MobilId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "MobilID tidak boleh kosong")
	}

	log.Printf("TransaksiService: BuyMobil dipanggil oleh %s untuk mobil %s", pembeliID, req.MobilId)

	// 2. Mulai Transaksi Database (PENTING!)
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal memulai transaksi DB")
	}
	// Pastikan di-rollback jika ada error
	defer tx.Rollback()

	// 3. Kunci mobil dan cek status (PENTING: FOR UPDATE)
	var penjualID, statusMobil, merkMobil, modelMobil string
	var hargaJual float64
	queryCek := `SELECT owner_id, harga_jual, status, merk, model FROM mobils WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCek, req.MobilId).Scan(&penjualID, &hargaJual, &statusMobil, &merkMobil, &modelMobil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Mobil tidak ditemukan")
		}
		return nil, status.Errorf(codes.Internal, "Gagal mengecek mobil")
	}

	if statusMobil != "tersedia" {
		return nil, status.Errorf(codes.FailedPrecondition, "Mobil saat ini tidak tersedia")
	}
	if penjualID == pembeliID {
		return nil, status.Errorf(codes.FailedPrecondition, "Anda tidak bisa membeli mobil Anda sendiri")
	}

	// 4. Update status mobil
	queryUpdate := `UPDATE mobils SET status = 'terjual' WHERE id = $1`
	_, err = tx.ExecContext(ctx, queryUpdate, req.MobilId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal update status mobil")
	}

	// 5. Buat catatan transaksi
	var resp pb.TransaksiJualResponse
	statusTransaksi := "selesai" // Langsung selesai (simulasi pembayaran sukses)
	queryInsert := `
		INSERT INTO transaksi_jual (mobil_id, penjual_id, pembeli_id, total, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, mobil_id, penjual_id, pembeli_id, total, status
	`
	err = tx.QueryRowContext(ctx, queryInsert, req.MobilId, penjualID, pembeliID, hargaJual, statusTransaksi).
		Scan(&resp.Id, &resp.MobilId, &resp.PenjualId, &resp.PembeliId, &resp.Total, &resp.Status)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal mencatat transaksi")
	}

	// 6. Commit Transaksi DB
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal menyelesaikan transaksi")
	}

	log.Printf("Transaksi sukses: Mobil %s dibeli oleh %s", req.MobilId, pembeliID)
	// (Nanti di Langkah 6, kita buat notifikasi di sini)
	// --- 7. BUAT NOTIFIKASI (SETELAH COMMIT) ---
	pesanMobil := fmt.Sprintf("%s %s", merkMobil, modelMobil)
	tanggalSekarang := time.Now().Format("02 Jan 2006")

	// Notifikasi untuk Pembeli (jalankan di goroutine baru)
	go notifikasi.CreateNotification(s.DB, context.Background(), pembeliID, "beli",
		fmt.Sprintf("Anda melakukan pembelian mobil %s pada tanggal %s dengan harga Rp %.0f",
			pesanMobil, tanggalSekarang, hargaJual))

	// Notifikasi untuk Penjual (jalankan di goroutine baru)
	go notifikasi.CreateNotification(s.DB, context.Background(), penjualID, "jual",
		fmt.Sprintf("Anda melakukan penjualan mobil %s pada tanggal %s dengan harga Rp %.0f",
			pesanMobil, tanggalSekarang, hargaJual))

	return &resp, nil
}

// PENJELASAN FILE transaksi_service.go:
// File ini menangani transaksi jual beli mobil
//
// Fungsi BuyMobil (Pembelian Mobil):
// - Ambil pembeli_id dari context (user yang membeli)
// - Mulai database transaction (PENTING untuk data consistency)
// - Lock mobil dengan FOR UPDATE (prevent race condition)
// - Validasi: mobil harus tersedia, pembeli != penjual
// - Update status mobil jadi 'terjual'
// - Insert record ke transaksi_jual dengan status 'selesai'
// - Commit transaction
// - Buat notifikasi untuk pembeli dan penjual (goroutine background)
//
// Keamanan Transaction:
// - FOR UPDATE: Lock row mobil saat transaction (prevent double booking)
// - tx.Rollback(): Otomatis rollback jika ada error
// - tx.Commit(): Save semua perubahan sekaligus (atomic operation)
//
// Database Tables:
// - mobils: Data mobil, status diupdate saat transaksi
// - transaksi_jual: Record pembelian mobil
