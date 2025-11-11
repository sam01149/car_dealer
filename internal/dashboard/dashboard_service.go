package dashboard

import (
	"context"
	"database/sql"
	"log"
	"sync" // Kita akan menggunakan WaitGroup untuk query paralel

	"carapp.com/m/internal/auth" // Sesuaikan dengan modul Anda
	pb "carapp.com/m/proto"      // Sesuaikan modul Anda
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DashboardServiceServer adalah implementasi dari pb.DashboardServiceServer
type DashboardServiceServer struct {
	pb.UnimplementedDashboardServiceServer
	DB *sql.DB
}

// NewDashboardService membuat instance baru
func NewDashboardService(db *sql.DB) *DashboardServiceServer {
	return &DashboardServiceServer{DB: db}
}

// GetDashboard mengambil data ringkasan untuk dashboard user
func (s *DashboardServiceServer) GetDashboard(ctx context.Context, req *emptypb.Empty) (*pb.DashboardSummary, error) {
	// 1. Dapatkan UserID dari token
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak dapat mengambil UserID dari token")
	}

	log.Printf("DashboardService: Mengambil dashboard untuk UserID %s", userID)

	var resp pb.DashboardSummary
	var wg sync.WaitGroup             // Untuk menjalankan query secara bersamaan
	var errChan = make(chan error, 4) // Channel untuk menangkap error dari goroutine

	// 2. Jalankan 4 query secara paralel
	wg.Add(4)

	// Goroutine 1: Hitung total mobil milik user
	go func() {
		defer wg.Done()
		query := `SELECT COUNT(*) FROM mobils WHERE owner_id = $1`
		err := s.DB.QueryRowContext(ctx, query, userID).Scan(&resp.TotalMobilAnda)
		if err != nil {
			errChan <- err
		}
	}()

	// Goroutine 2: Hitung transaksi aktif
	go func() {
		defer wg.Done()
		// Hanya transaksi jual yang "diproses" yang melibatkan user
		query := `
			SELECT COUNT(*) FROM transaksi_jual 
			WHERE (penjual_id = $1 OR pembeli_id = $1) AND status = 'diproses'
		`
		err := s.DB.QueryRowContext(ctx, query, userID).Scan(&resp.TransaksiAktif)
		if err != nil {
			errChan <- err
		}
	}()

	// Goroutine 3: Hitung pendapatan terakhir (Total penjualan mobil)
	go func() {
		defer wg.Done()
		// Menggunakan COALESCE untuk memastikan 0 jika hasilnya NULL
		query := `SELECT COALESCE(SUM(total), 0) FROM transaksi_jual WHERE penjual_id = $1 AND status = 'selesai'`
		err := s.DB.QueryRowContext(ctx, query, userID).Scan(&resp.PendapatanTerakhir)
		if err != nil {
			errChan <- err
		}
	}()

	// Goroutine 4: Hitung notifikasi baru (belum dibaca)
	go func() {
		defer wg.Done()
		query := `SELECT COUNT(*) FROM notifikasi WHERE user_id = $1 AND read_at IS NULL`
		err := s.DB.QueryRowContext(ctx, query, userID).Scan(&resp.NotifikasiBaru)
		if err != nil {
			errChan <- err
		}
	}()

	// 3. Tunggu semua query selesai
	wg.Wait()
	close(errChan)

	// Cek apakah ada error dari goroutine
	if err := <-errChan; err != nil {
		log.Printf("Gagal mengambil data dashboard: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal memuat data dashboard")
	}

	log.Printf("Dashboard untuk %s: Mobil=%d, Transaksi=%d, Notif=%d",
		userID, resp.TotalMobilAnda, resp.TransaksiAktif, resp.NotifikasiBaru)

	// 4. Kembalikan response
	return &resp, nil
}

// PENJELASAN FILE dashboard_service.go:
// File ini menyediakan summary data untuk dashboard user
//
// Fungsi GetDashboard:
// - Mengambil user_id dari context (sudah tervalidasi di middleware)
// - Jalankan 4 query secara PARALEL menggunakan goroutine & WaitGroup:
//   1. Total mobil milik user (COUNT dari mobils WHERE owner_id)
//   2. Transaksi aktif (COUNT transaksi_jual status='diproses')
//   3. Pendapatan terakhir (SUM total dari transaksi_jual WHERE penjual_id dan status='selesai')
//   4. Notifikasi baru (COUNT notifikasi WHERE user_id dan read_at IS NULL)
// - Gunakan errChan untuk capture error dari goroutine
// - Return DashboardSummary dengan semua data aggregate
//
// Keuntungan parallel query:
// - Lebih cepat daripada query sequential (4 query jadi 1x waktu query terlama)
// - Efisien untuk dashboard yang butuh banyak data
