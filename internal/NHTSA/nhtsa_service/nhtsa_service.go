package nhtsa_service

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"carapp.com/m/internal/nhtsa" // Sesuaikan dengan modul Anda
	pb "carapp.com/m/proto"       // Sesuaikan dengan modul Anda
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const cacheTTL = 24 * time.Hour

// NhtsaDataServiceServer adalah implementasi dari pb.NhtsaDataServiceServer
type NhtsaDataServiceServer struct {
	pb.UnimplementedNhtsaDataServiceServer
	DB *sql.DB
}

// NewNhtsaDataService membuat instance baru
func NewNhtsaDataService(db *sql.DB) *NhtsaDataServiceServer {
	return &NhtsaDataServiceServer{DB: db}
}

// GetMakes (Logika ini DIPINDAHKAN dari mobil_service.go)
func (s *NhtsaDataServiceServer) GetMakes(ctx context.Context, req *pb.GetMakesRequest) (*pb.GetMakesResponse, error) {
	// ... (Salin-tempel SEMUA kode fungsi GetMakes dari internal/mobil/mobil_service.go ke sini) ...
	// ... (Termasuk helper getMakesFromCache dan saveMakesToCache) ...
	log.Println("NhtsaDataService: GetMakes dipanggil")
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

// GetModelsForMake (Logika ini DIPINDAHKAN dari mobil_service.go)
func (s *NhtsaDataServiceServer) GetModelsForMake(ctx context.Context, req *pb.GetModelsForMakeRequest) (*pb.GetModelsForMakeResponse, error) {
	// ... (Salin-tempel SEMUA kode fungsi GetModelsForMake dari internal/mobil/mobil_service.go ke sini) ...
	// ... (Termasuk helper getModelsFromCache dan saveModelsToCache) ...
	log.Println("NhtsaDataService: GetModelsForMake dipanggil")
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

// Helper: Ambil makes dari cache DB
func (s *NhtsaDataServiceServer) getMakesFromCache(ctx context.Context) ([]*pb.Make, error) {
	query := `
		SELECT brand_id, name 
		FROM nhtsa_makes_cache
		WHERE cached_at > NOW() - INTERVAL '24 hours'
		ORDER BY name
	`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []*pb.Make
	for rows.Next() {
		var m pb.Make
		if err := rows.Scan(&m.BrandId, &m.Name); err != nil {
			return nil, err
		}
		makes = append(makes, &m)
	}

	if len(makes) == 0 {
		return nil, sql.ErrNoRows
	}

	return makes, nil
}

// Helper: Simpan makes ke cache DB
func (s *NhtsaDataServiceServer) saveMakesToCache(apiMakes []nhtsa.NhtsaMake) {
	// Hapus cache lama
	_, err := s.DB.Exec("DELETE FROM nhtsa_makes_cache")
	if err != nil {
		log.Printf("Gagal menghapus cache lama: %v", err)
		return
	}

	// Simpan data baru
	for _, m := range apiMakes {
		query := `INSERT INTO nhtsa_makes_cache (brand_id, name, cached_at)
		          VALUES ($1, $2, NOW())
		          ON CONFLICT (brand_id) DO UPDATE SET name = $2, cached_at = NOW()`
		_, err := s.DB.Exec(query, strconv.Itoa(m.MakeID), m.MakeName)
		if err != nil {
			log.Printf("Gagal menyimpan make %s: %v", m.MakeName, err)
		}
	}
	log.Println("Cache makes berhasil disimpan")
}

// Helper: Ambil models dari cache DB
func (s *NhtsaDataServiceServer) getModelsFromCache(ctx context.Context, brandID string) ([]*pb.Model, error) {
	query := `
		SELECT model_id, brand_id, name
		FROM nhtsa_models_cache
		WHERE brand_id = $1 AND cached_at > NOW() - INTERVAL '24 hours'
		ORDER BY name
	`
	rows, err := s.DB.QueryContext(ctx, query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*pb.Model
	for rows.Next() {
		var m pb.Model
		if err := rows.Scan(&m.ModelId, &m.BrandId, &m.Name); err != nil {
			return nil, err
		}
		models = append(models, &m)
	}

	if len(models) == 0 {
		return nil, sql.ErrNoRows
	}

	return models, nil
}

// Helper: Simpan models ke cache DB
func (s *NhtsaDataServiceServer) saveModelsToCache(apiModels []nhtsa.NhtsaModel) {
	for _, m := range apiModels {
		query := `INSERT INTO nhtsa_models_cache (model_id, brand_id, name, cached_at)
		          VALUES ($1, $2, $3, NOW())
		          ON CONFLICT (model_id) DO UPDATE SET name = $3, cached_at = NOW()`
		_, err := s.DB.Exec(query, strconv.Itoa(m.ModelID), strconv.Itoa(m.MakeID), m.ModelName)
		if err != nil {
			log.Printf("Gagal menyimpan model %s: %v", m.ModelName, err)
		}
	}
	log.Println("Cache models berhasil disimpan")
}
