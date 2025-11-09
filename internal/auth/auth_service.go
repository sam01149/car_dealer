package auth

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"carapp.com/m/internal/utils" // Sesuaikan dengan nama modul Anda
	pb "carapp.com/m/proto"       // Sesuaikan dengan nama modul Anda
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthServiceServer adalah implementasi dari pb.AuthServiceServer
type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	DB *sql.DB
}

// NewAuthService membuat instance baru dari AuthServiceServer
func NewAuthService(db *sql.DB) *AuthServiceServer {
	return &AuthServiceServer{DB: db}
}

// Register menangani pendaftaran user baru
func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	log.Printf("Menerima permintaan Register untuk email: %s", req.Email)

	// 1. Validasi input (sederhana)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Nama, Email, dan Password tidak boleh kosong")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Gagal hash password: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal memproses pendaftaran")
	}

	// 3. Simpan user ke database
	var userID, userName, userEmail, userPhone, userRole string
	var createdAt time.Time

	// Kita set 'client' sebagai role default
	defaultRole := "client"

	query := `INSERT INTO users (name, email, password_hash, phone, role)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, name, email, phone, role, created_at`

	err = s.DB.QueryRowContext(ctx, query, req.Name, req.Email, hashedPassword, req.Phone, defaultRole).
		Scan(&userID, &userName, &userEmail, &userPhone, &userRole, &createdAt)

	if err != nil {
		// Cek jika email sudah terdaftar (unique constraint violation)
		if strings.Contains(err.Error(), "unique constraint") {
			log.Printf("Email sudah terdaftar: %s", req.Email)
			return nil, status.Errorf(codes.AlreadyExists, "Email sudah terdaftar")
		}
		log.Printf("Gagal menyimpan user ke DB: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal menyimpan data user")
	}

	// 4. Buat JWT Token
	token, err := utils.GenerateToken(userID, userEmail, userRole)
	if err != nil {
		log.Printf("Gagal membuat token: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal membuat token")
	}

	// 5. Kembalikan response
	return &pb.AuthResponse{
		User: &pb.User{
			Id:        userID,
			Name:      userName,
			Email:     userEmail,
			Phone:     userPhone,
			Role:      userRole,
			CreatedAt: timestamppb.New(createdAt),
		},
		Token: token,
	}, nil
}

// Login menangani login user
func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	log.Printf("Menerima permintaan Login untuk email: %s", req.Email)

	// 1. Validasi input
	if req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Email dan Password tidak boleh kosong")
	}

	// 2. Cari user di database
	var userID, userName, userEmail, userPhone, userRole, hashedPassword string
	var createdAt time.Time

	query := `SELECT id, name, email, phone, role, password_hash, created_at FROM users WHERE email = $1`

	err := s.DB.QueryRowContext(ctx, query, req.Email).
		Scan(&userID, &userName, &userEmail, &userPhone, &userRole, &hashedPassword, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User tidak ditemukan: %s", req.Email)
			return nil, status.Errorf(codes.NotFound, "Email atau Password salah")
		}
		log.Printf("Gagal query DB: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal memproses login")
	}

	// 3. Verifikasi password
	if !utils.CheckPasswordHash(req.Password, hashedPassword) {
		log.Printf("Password salah untuk: %s", req.Email)
		return nil, status.Errorf(codes.Unauthenticated, "Email atau Password salah")
	}

	// 4. Buat JWT Token
	token, err := utils.GenerateToken(userID, userEmail, userRole)
	if err != nil {
		log.Printf("Gagal membuat token: %v", err)
		return nil, status.Errorf(codes.Internal, "Gagal membuat token")
	}

	// 5. Kembalikan response
	return &pb.AuthResponse{
		User: &pb.User{
			Id:        userID,
			Name:      userName,
			Email:     userEmail,
			Phone:     userPhone,
			Role:      userRole,
			CreatedAt: timestamppb.New(createdAt),
		},
		Token: token,
	}, nil
}
