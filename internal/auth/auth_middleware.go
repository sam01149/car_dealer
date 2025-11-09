package auth

import (
	"context"
	"log"
	"strings"

	"carapp.com/m/internal/utils" // Sesuaikan dengan modul Anda
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Definisikan kunci untuk context
type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

// AuthInterceptor adalah gRPC Unary Interceptor untuk validasi JWT
func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// Daftar service/method yang tidak perlu dicek token-nya
	publicMethods := map[string]bool{
		"/carapp.AuthService/Login":                 true,
		"/carapp.AuthService/Register":              true,
		"/carapp.NhtsaDataService/GetMakes":         true, // Diubah dari MobilService
		"/carapp.NhtsaDataService/GetModelsForMake": true, // Diubah dari MobilService

		// --- TAMBAHAN BARU ---
		"/carapp.MobilService/ListMobil": true, // Publik bisa lihat daftar mobil
		"/carapp.MobilService/GetMobil":  true,
	}

	// Cek apakah method ini publik
	if publicMethods[info.FullMethod] {
		// Jika ya, langsung teruskan ke handler tanpa cek token
		return handler(ctx, req)
	}

	log.Printf("--> AuthInterceptor: Memvalidasi method %s", info.FullMethod)

	// 1. Ambil metadata dari context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata tidak ditemukan")
	}

	// 2. Ambil nilai 'authorization'
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Token authorization tidak ditemukan")
	}

	// 3. Token biasanya dalam format "Bearer <token>"
	authHeader := authHeaders[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "Format token salah")
	}
	tokenString := parts[1]

	// 4. Validasi token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token tidak valid: %v", err)
	}

	// 5. Token valid. Simpan info user di context
	ctxWithUser := context.WithValue(ctx, UserIDKey, claims.UserID)
	ctxWithUser = context.WithValue(ctxWithUser, UserEmailKey, claims.Email)
	ctxWithUser = context.WithValue(ctxWithUser, UserRoleKey, claims.Role)

	log.Printf("Token valid untuk UserID: %s", claims.UserID)

	// 6. Teruskan ke handler asli dengan context yang sudah berisi info user
	return handler(ctxWithUser, req)
}
