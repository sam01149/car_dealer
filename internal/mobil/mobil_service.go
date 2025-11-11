package mobil

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"carapp.com/m/internal/auth"
	"carapp.com/m/internal/nhtsa"
	"carapp.com/m/internal/notifikasi"
	pb "carapp.com/m/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const cacheTTL = 24 * time.Hour // Sesuai rencana: TTL cache 24 jam

// MobilServiceServer adalah implementasi dari pb.MobilServiceServer
type MobilServiceServer struct {
	pb.UnimplementedMobilServiceServer
	DB *sql.DB
}

// NewMobilService membuat instance baru dari MobilServiceServer
func NewMobilService(db *sql.DB) *MobilServiceServer {
	return &MobilServiceServer{DB: db}
}

// GetMakes mengambil daftar merek, menggunakan cache
func (s *MobilServiceServer) GetMakes(ctx context.Context, req *pb.GetMakesRequest) (*pb.GetMakesResponse, error) {
	// 1. Coba ambil dari cache
	makes, err := s.getMakesFromCache(ctx)
	if err == nil {
		log.Println("GetMakes: Data disajikan dari cache DB")
		return &pb.GetMakesResponse{Makes: makes}, nil
	}

	// 2. Jika cache gagal (atau kedaluwarsa), ambil dari API
	log.Println("GetMakes: Cache kosong atau kedaluwarsa, mengambil dari NHTSA API...")
	apiMakes, err := nhtsa.FetchAllMakes()
	if err != nil {
		log.Printf("Gagal mengambil data dari NHTSA: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal mengambil data merek: %v", err)
	}

	// 3. Simpan ke cache (di latar belakang agar tidak memblokir user)
	go s.saveMakesToCache(apiMakes)

	// 4. Konversi ke format proto dan kembalikan
	var protoMakes []*pb.Make
	for _, m := range apiMakes {
		protoMakes = append(protoMakes, &pb.Make{
			BrandId: strconv.Itoa(m.MakeID), // Konversi int ke string
			Name:    m.MakeName,
		})
	}

	return &pb.GetMakesResponse{Makes: protoMakes}, nil
}

// CreateMobil menangani pembuatan data mobil baru
func (s *MobilServiceServer) CreateMobil(ctx context.Context, req *pb.CreateMobilRequest) (*pb.Mobil, error) {
	log.Println("MobilService: CreateMobil dipanggil")
	// 1. Ambil UserID dari context (yang diisi oleh middleware)
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak dapat mengambil UserID dari token")
	}

	// 2. Validasi input
	if req.Merk == "" || req.Model == "" || req.Tahun <= 1900 || req.HargaJual <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Data mobil tidak valid (Merk, Model, Tahun, Harga Jual)")
	}

	if req.Deskripsi == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Deskripsi harus diisi")
	}

	// Bulatkan harga untuk menghindari floating-point precision issue
	hargaJualBulat := math.Round(req.HargaJual)

	// 3. Simpan ke database
	query := `
		INSERT INTO mobils (
			owner_id, merk, model, tahun, kondisi, deskripsi, 
			harga_jual, foto_url, lokasi, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, owner_id, merk, model, tahun, kondisi, deskripsi, 
		          harga_jual, foto_url, lokasi, status, created_at
	`

	var mobil pb.Mobil
	var createdAt time.Time

	err := s.DB.QueryRowContext(ctx, query,
		userID,
		req.Merk,
		req.Model,
		req.Tahun,
		req.Kondisi,
		req.Deskripsi,
		hargaJualBulat,
		req.FotoUrl, // Simpan foto_url dari request
		req.Lokasi,
		"tersedia",
	).Scan(
		&mobil.Id,
		&mobil.OwnerId,
		&mobil.Merk,
		&mobil.Model,
		&mobil.Tahun,
		&mobil.Kondisi,
		&mobil.Deskripsi,
		&mobil.HargaJual,
		&mobil.FotoUrl, // Scan langsung ke FotoUrl
		&mobil.Lokasi,
		&mobil.Status,
		&createdAt,
	)

	if err != nil {
		log.Printf("Gagal menyimpan mobil ke DB: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal menyimpan mobil")
	}

	mobil.CreatedAt = timestamppb.New(createdAt)
	log.Printf("Mobil baru berhasil dibuat oleh UserID %s (MobilID: %s)", userID, mobil.Id)

	// Buat notifikasi untuk penjual
	pesanMobil := fmt.Sprintf("%d %s %s", mobil.Tahun, mobil.Merk, mobil.Model)
	go notifikasi.CreateNotification(s.DB, context.Background(), userID, "jual",
		fmt.Sprintf("Anda berhasil memasang iklan jual mobil %s dengan harga Rp %.0f pada tanggal %s.",
			pesanMobil, mobil.HargaJual, createdAt.Format("02 Jan 2006")))

	return &mobil, nil
}

// Implementasi dummy (Akan kita isi di Langkah 5)
func (s *MobilServiceServer) ListMobil(ctx context.Context, req *pb.ListMobilRequest) (*pb.ListMobilResponse, error) {
	log.Println("MobilService: ListMobil dipanggil")

	// Filter status dasar, default ke "tersedia"
	filterStatus := "tersedia"
	if req.FilterStatus != nil {
		filterStatus = *req.FilterStatus
	}

	// Logika paginasi sederhana
	limit := 50 // Sesuai permintaan Anda "50 mobil"
	if req.Limit > 0 {
		limit = int(req.Limit)
	}
	offset := 0
	if req.Page > 1 {
		offset = (int(req.Page) - 1) * limit
	}

	// Query untuk mengambil mobil
	query := `
		SELECT id, owner_id, merk, model, tahun, kondisi, deskripsi, 
		       harga_jual, foto_url, lokasi, status, created_at
		FROM mobils
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.DB.QueryContext(ctx, query, filterStatus, limit, offset)
	if err != nil {
		log.Printf("Gagal query ListMobil: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal mengambil data mobil")
	}
	defer rows.Close()

	var mobils []*pb.Mobil
	for rows.Next() {
		var mobil pb.Mobil
		var createdAt time.Time
		var fotoUrl sql.NullString

		err := rows.Scan(
			&mobil.Id, &mobil.OwnerId, &mobil.Merk, &mobil.Model, &mobil.Tahun,
			&mobil.Kondisi, &mobil.Deskripsi, &mobil.HargaJual, &fotoUrl,
			&mobil.Lokasi, &mobil.Status, &createdAt,
		)
		if err != nil {
			log.Printf("Gagal scan row mobil: %v", err)
			continue
		}

		// Set foto_url jika ada
		if fotoUrl.Valid {
			mobil.FotoUrl = fotoUrl.String
		}
		mobil.CreatedAt = timestamppb.New(createdAt)
		mobils = append(mobils, &mobil)
	}

	// Query untuk total (untuk paginasi)
	var total int32
	s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM mobils WHERE status = $1", filterStatus).Scan(&total)

	return &pb.ListMobilResponse{
		Mobils: mobils,
		Total:  total,
	}, nil
}

// GetMobil mengambil detail satu mobil (Fitur 3)
func (s *MobilServiceServer) GetMobil(ctx context.Context, req *pb.GetMobilRequest) (*pb.Mobil, error) {
	log.Printf("MobilService: GetMobil dipanggil untuk ID: %s", req.MobilId)

	if req.MobilId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "MobilID tidak boleh kosong")
	}

	query := `
		SELECT m.id, m.owner_id, u.name as owner_name, m.merk, m.model, m.tahun, m.kondisi, m.deskripsi, 
		       m.harga_jual, m.foto_url, m.lokasi, m.status, m.created_at
		FROM mobils m
		LEFT JOIN users u ON m.owner_id = u.id
		WHERE m.id = $1
	`
	var mobil pb.Mobil
	var createdAt time.Time
	var fotoUrl sql.NullString
	var ownerName string

	err := s.DB.QueryRowContext(ctx, query, req.MobilId).Scan(
		&mobil.Id, &mobil.OwnerId, &ownerName, &mobil.Merk, &mobil.Model, &mobil.Tahun,
		&mobil.Kondisi, &mobil.Deskripsi, &mobil.HargaJual, &fotoUrl,
		&mobil.Lokasi, &mobil.Status, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Mobil tidak ditemukan")
		}
		log.Printf("Gagal query GetMobil: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal mengambil data mobil")
	}

	// Set owner name
	mobil.OwnerName = ownerName
	log.Printf("Owner name: %s", ownerName)

	// Set foto_url jika ada
	if fotoUrl.Valid {
		mobil.FotoUrl = fotoUrl.String
	}
	mobil.CreatedAt = timestamppb.New(createdAt)

	return &mobil, nil
}

// getMakesFromCache helper untuk cek cache DB
func (s *MobilServiceServer) getMakesFromCache(ctx context.Context) ([]*pb.Make, error) {
	query := `SELECT brand_id, name, updated_at FROM brand_cache`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []*pb.Make
	var lastUpdate time.Time

	for rows.Next() {
		var make pb.Make
		var updatedAt time.Time
		if err := rows.Scan(&make.BrandId, &make.Name, &updatedAt); err != nil {
			return nil, err
		}
		makes = append(makes, &make)
		if updatedAt.After(lastUpdate) {
			lastUpdate = updatedAt
		}
	}

	// Cek apakah cache kosong atau kedaluwarsa
	if len(makes) == 0 || time.Since(lastUpdate) > cacheTTL {
		return nil, status.Error(codes.NotFound, "Cache kosong atau kedaluwarsa")
	}

	return makes, nil
}

// saveMakesToCache helper untuk menyimpan/update cache
func (s *MobilServiceServer) saveMakesToCache(apiMakes []nhtsa.NhtsaMake) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Gagal memulai transaksi cache: %v", err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO brand_cache (brand_id, name, raw, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (brand_id) DO UPDATE
		SET name = EXCLUDED.name, raw = EXCLUDED.raw, updated_at = NOW()
	`)
	if err != nil {
		log.Printf("Gagal menyiapkan statement cache: %v", err)
		return
	}
	defer stmt.Close()

	for _, m := range apiMakes {
		rawJson, _ := json.Marshal(m)
		brandIDStr := strconv.Itoa(m.MakeID)
		if _, err := stmt.ExecContext(ctx, brandIDStr, m.MakeName, rawJson); err != nil {
			log.Printf("Gagal menyimpan brand_id %s ke cache: %v", brandIDStr, err)
			continue // Lanjut ke item berikutnya
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Gagal commit transaksi cache merek: %v", err)
	}

	log.Printf("Sukses menyimpan %d merek ke cache DB", len(apiMakes))
}

// GetModelsForMake mengambil model untuk merek, menggunakan cache
func (s *MobilServiceServer) GetModelsForMake(ctx context.Context, req *pb.GetModelsForMakeRequest) (*pb.GetModelsForMakeResponse, error) {
	if req.BrandId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "BrandId tidak boleh kosong")
	}

	// 1. Coba ambil dari cache
	models, err := s.getModelsFromCache(ctx, req.BrandId)
	if err == nil {
		log.Printf("GetModelsForMake: Data untuk %s disajikan dari cache DB", req.BrandId)
		return &pb.GetModelsForMakeResponse{Models: models}, nil
	}

	// 2. Jika cache gagal, ambil dari API
	log.Printf("GetModelsForMake: Cache untuk %s kosong, mengambil dari NHTSA API...", req.BrandId)
	apiModels, err := nhtsa.FetchModelsForMakeID(req.BrandId)
	if err != nil {
		log.Printf("Gagal mengambil model dari NHTSA: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal mengambil data model: %v", err)
	}

	// 3. Simpan ke cache (di latar belakang)
	go s.saveModelsToCache(apiModels)

	// 4. Konversi dan kembalikan
	var protoModels []*pb.Model
	for _, m := range apiModels {
		protoModels = append(protoModels, &pb.Model{
			ModelId: strconv.Itoa(m.ModelID),
			BrandId: strconv.Itoa(m.MakeID),
			Name:    m.ModelName,
		})
	}

	return &pb.GetModelsForMakeResponse{Models: protoModels}, nil
}

// getModelsFromCache helper
func (s *MobilServiceServer) getModelsFromCache(ctx context.Context, brandID string) ([]*pb.Model, error) {
	query := `SELECT model_id, brand_id, name, updated_at FROM model_cache WHERE brand_id = $1`
	rows, err := s.DB.QueryContext(ctx, query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*pb.Model
	var lastUpdate time.Time

	for rows.Next() {
		var model pb.Model
		var updatedAt time.Time
		if err := rows.Scan(&model.ModelId, &model.BrandId, &model.Name, &updatedAt); err != nil {
			return nil, err
		}
		models = append(models, &model)
		if updatedAt.After(lastUpdate) {
			lastUpdate = updatedAt
		}
	}

	if len(models) == 0 || time.Since(lastUpdate) > cacheTTL {
		return nil, status.Error(codes.NotFound, "Cache model kosong atau kedaluwarsa")
	}

	return models, nil
}

// saveModelsToCache helper
func (s *MobilServiceServer) saveModelsToCache(apiModels []nhtsa.NhtsaModel) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Gagal memulai transaksi cache model: %v", err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO model_cache (model_id, brand_id, name, raw, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (model_id) DO UPDATE
		SET brand_id = EXCLUDED.brand_id, name = EXCLUDED.name, raw = EXCLUDED.raw, updated_at = NOW()
	`)
	if err != nil {
		log.Printf("Gagal menyiapkan statement cache model: %v", err)
		return
	}
	defer stmt.Close()

	for _, m := range apiModels {
		rawJson, _ := json.Marshal(m)
		modelIDStr := strconv.Itoa(m.ModelID)
		brandIDStr := strconv.Itoa(m.MakeID)
		if _, err := stmt.ExecContext(ctx, modelIDStr, brandIDStr, m.ModelName, rawJson); err != nil {
			log.Printf("Gagal menyimpan model_id %s ke cache: %v", modelIDStr, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Gagal commit transaksi cache model: %v", err)
	}

	log.Printf("Sukses menyimpan %d model ke cache DB", len(apiModels))
}

// PENJELASAN FILE mobil_service.go:
// File ini menangani semua operasi terkait mobil (CRUD + NHTSA data)
//
// Fungsi CreateMobil:
// - Ambil user_id dari context (owner mobil)
// - Validasi input (merk, model, tahun, harga_jual harus valid)
// - Bulatkan harga untuk menghindari floating-point precision issue
// - Insert mobil baru ke database dengan status 'tersedia'
// - Buat notifikasi untuk penjual (goroutine background)
// - Return data mobil yang baru dibuat
//
// Fungsi ListMobil:
// - Query daftar mobil dengan paginasi (default 50 mobil per page)
// - Filter berdasarkan status (default: 'tersedia')
// - Support limit dan offset untuk pagination
// - Order by created_at DESC (mobil terbaru di atas)
// - Handle harga_rental NULL dengan sql.NullFloat64
// - Return list mobil + total count
//
// Fungsi GetMobil:
// - Query detail satu mobil berdasarkan mobil_id
// - Return 404 NotFound jika mobil tidak ada
// - Support untuk public access (tidak perlu login)
//
// Fungsi GetMakes:
// - Coba ambil list merek dari cache DB
// - Jika cache kosong/expired -> fetch dari NHTSA API
// - Simpan ke cache di background (goroutine)
// - Return list merek mobil
//
// Fungsi GetModelsForMake:
// - Coba ambil list model untuk brand_id dari cache
// - Jika cache kosong -> fetch dari NHTSA API
// - Simpan ke cache di background
// - Return list model untuk merek tertentu
//
// Helper Functions:
// - getMakesFromCache: Query brand_cache, cek TTL 24 jam
// - saveMakesToCache: UPSERT ke brand_cache dengan transaction
// - getModelsFromCache: Query model_cache untuk brand_id
// - saveModelsToCache: UPSERT ke model_cache
//
// Database Tables:
// - mobils: Data mobil user
// - brand_cache: Cache merek dari NHTSA
// - model_cache: Cache model dari NHTSA

// PENJELASAN FILE mobil_service.go:
// File ini menangani semua operasi terkait mobil (CRUD + NHTSA data)
//
// Fungsi CreateMobil:
// - Ambil user_id dari context (owner mobil)
// - Validasi input (merk, model, tahun, harga_jual harus valid)
// - Bulatkan harga untuk menghindari floating-point precision issue
// - Insert mobil baru ke database dengan status 'tersedia'
// - Buat notifikasi untuk penjual (goroutine background)
// - Return data mobil yang baru dibuat
//
// Fungsi ListMobil:
// - Query daftar mobil dengan paginasi (default 50 mobil per page)
// - Filter berdasarkan status (default: 'tersedia')
// - Support limit dan offset untuk pagination
// - Order by created_at DESC (mobil terbaru di atas)
// - Handle harga_rental NULL dengan sql.NullFloat64
// - Return list mobil + total count
//
// Fungsi GetMobil:
// - Query detail satu mobil berdasarkan mobil_id
// - Return 404 NotFound jika mobil tidak ada
// - Support untuk public access (tidak perlu login)
//
// Fungsi GetMakes:
// - Coba ambil list merek dari cache DB
// - Jika cache kosong/expired -> fetch dari NHTSA API
// - Simpan ke cache di background (goroutine)
// - Return list merek mobil
//
// Fungsi GetModelsForMake:
// - Coba ambil list model untuk brand_id dari cache
// - Jika cache kosong -> fetch dari NHTSA API
// - Simpan ke cache di background
// - Return list model untuk merek tertentu
//
// Helper Functions:
// - getMakesFromCache: Query brand_cache, cek TTL 24 jam
// - saveMakesToCache: UPSERT ke brand_cache dengan transaction
// - getModelsFromCache: Query model_cache untuk brand_id
// - saveModelsToCache: UPSERT ke model_cache
//
// Database Tables:
// - mobils: Data mobil user
// - brand_cache: Cache merek dari NHTSA
// - model_cache: Cache model dari NHTSA
