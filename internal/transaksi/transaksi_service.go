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

	// Notifikasi untuk Pembeli (jalankan di goroutine baru)
	go notifikasi.CreateNotification(s.DB, context.Background(), pembeliID, "beli",
		fmt.Sprintf("Anda telah berhasil membeli mobil %s seharga Rp%.2f.", pesanMobil, hargaJual))

	// Notifikasi untuk Penjual (jalankan di goroutine baru)
	go notifikasi.CreateNotification(s.DB, context.Background(), penjualID, "jual",
		fmt.Sprintf("Mobil Anda %s telah terjual kepada user lain seharga Rp%.2f.", pesanMobil, hargaJual))

	return &resp, nil
}

// RentMobil menangani logika rental mobil (Fitur 5)
func (s *TransaksiServiceServer) RentMobil(ctx context.Context, req *pb.RentMobilRequest) (*pb.TransaksiRentalResponse, error) {
	// 1. Dapatkan ID penyewa dari token
	penyewaID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak dapat mengambil UserID dari token")
	}

	log.Printf("TransaksiService: RentMobil dipanggil oleh %s untuk mobil %s", penyewaID, req.MobilId)

	// 2. Validasi Tanggal
	layout := "2006-01-02"
	tglMulai, errMulai := time.Parse(layout, req.TanggalMulai)
	tglSelesai, errSelesai := time.Parse(layout, req.TanggalSelesai)
	if errMulai != nil || errSelesai != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Format tanggal salah, gunakan YYYY-MM-DD")
	}
	if tglMulai.After(tglSelesai) || tglMulai.Before(time.Now().Truncate(24*time.Hour)) {
		return nil, status.Errorf(codes.InvalidArgument, "Periode tanggal tidak valid")
	}

	durasiHari := int(tglSelesai.Sub(tglMulai).Hours()/24) + 1

	// 3. Mulai Transaksi Database
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal memulai transaksi DB")
	}
	defer tx.Rollback()

	// 4. Kunci mobil dan cek status
	var pemilikID, statusMobil, merkMobil, modelMobil string // <-- Tambah merk/model
	var hargaRental float64

	// --- UBAH QUERY INI ---
	queryCek := `SELECT owner_id, harga_rental_per_hari, status, merk, model FROM mobils WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCek, req.MobilId).Scan(&pemilikID, &hargaRental, &statusMobil, &merkMobil, &modelMobil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal mengecek mobil")
	}

	if statusMobil != "tersedia" {
		return nil, status.Errorf(codes.FailedPrecondition, "Mobil saat ini tidak tersedia untuk rental")
	}
	if hargaRental <= 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Mobil ini tidak untuk dirental")
	}

	// 5. Update status mobil
	queryUpdate := `UPDATE mobils SET status = 'dirental' WHERE id = $1`
	_, err = tx.ExecContext(ctx, queryUpdate, req.MobilId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal update status mobil")
	}

	// 6. Buat catatan transaksi
	totalBiaya := float64(durasiHari) * hargaRental
	var resp pb.TransaksiRentalResponse
	statusTransaksi := "aktif"
	dendaPerHari := 50000.0 // Contoh denda

	queryInsert := `
		INSERT INTO transaksi_rental (mobil_id, pemilik_id, penyewa_id, tanggal_mulai, tanggal_selesai, total, status, denda_per_hari)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, mobil_id, pemilik_id, penyewa_id, tanggal_mulai, tanggal_selesai, total, status
	`
	err = tx.QueryRowContext(ctx, queryInsert,
		req.MobilId, pemilikID, penyewaID, tglMulai, tglSelesai, totalBiaya, statusTransaksi, dendaPerHari,
	).Scan(&resp.Id, &resp.MobilId, &resp.PemilikId, &resp.PenyewaId, &resp.TanggalMulai, &resp.TanggalSelesai, &resp.Total, &resp.Status)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal mencatat transaksi rental")
	}

	// 7. Commit Transaksi DB
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "Gagal menyelesaikan transaksi rental")
	}

	log.Printf("Transaksi sukses: Mobil %s dirental oleh %s", req.MobilId, penyewaID)
	pesanMobil := fmt.Sprintf("%s %s", merkMobil, modelMobil)

	// Notifikasi untuk Penyewa
	go notifikasi.CreateNotification(s.DB, context.Background(), penyewaID, "rental",
		fmt.Sprintf("Anda berhasil merental %s dari %s s/d %s.", pesanMobil, req.TanggalMulai, req.TanggalSelesai))

	// Notifikasi untuk Pemilik
	go notifikasi.CreateNotification(s.DB, context.Background(), pemilikID, "rental",
		fmt.Sprintf("Mobil Anda %s telah dirental dari %s s/d %s.", pesanMobil, req.TanggalMulai, req.TanggalSelesai))

	return &resp, nil
}

// CompleteRental (dummy untuk saat ini, bisa diimplementasikan penuh jika mau)
func (s *TransaksiServiceServer) CompleteRental(ctx context.Context, req *pb.CompleteRentalRequest) (*pb.TransaksiRentalResponse, error) {
	// TODO:
	// 1. Dapatkan UserID dari context, pastikan dia pemilik atau penyewa
	// 2. Cari transaksi_rental berdasarkan req.RentalId
	// 3. Cek apakah tanggal pengembalian telat (Hitung denda)
	// 4. Update status transaksi_rental -> 'selesai'
	// 5. Update status mobils -> 'tersedia'
	// 6. Buat notifikasi
	return nil, status.Errorf(codes.Unimplemented, "method CompleteRental not implemented")
}
