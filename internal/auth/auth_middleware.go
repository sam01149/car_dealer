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

// StreamAuthInterceptor adalah gRPC Stream Interceptor untuk validasi JWT
func StreamAuthInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// Daftar stream methods yang tidak perlu dicek token
	publicStreamMethods := map[string]bool{
		// Tambahkan stream public methods di sini jika ada
	}

	// Cek apakah method ini publik
	if publicStreamMethods[info.FullMethod] {
		return handler(srv, ss)
	}

	log.Printf("--> StreamAuthInterceptor: Memvalidasi stream method %s", info.FullMethod)

	// 1. Ambil metadata dari context
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("GAGAL: Metadata tidak ditemukan untuk stream method %s", info.FullMethod)
		return status.Errorf(codes.Unauthenticated, "Metadata tidak ditemukan")
	}

	log.Printf("Metadata ditemukan: %v", md)

	// 2. Ambil nilai 'authorization'
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		log.Printf("GAGAL: Token authorization tidak ditemukan di metadata")
		log.Printf("Available metadata keys: %v", md)
		return status.Errorf(codes.Unauthenticated, "Token authorization tidak ditemukan")
	}

	log.Printf("Authorization header ditemukan: %s", authHeaders[0][:20]+"...")

	// 3. Token biasanya dalam format "Bearer <token>"
	authHeader := authHeaders[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return status.Errorf(codes.Unauthenticated, "Format token salah")
	}
	tokenString := parts[1]

	// 4. Validasi token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Token tidak valid: %v", err)
	}

	// 5. Token valid. Simpan info user di context
	ctxWithUser := context.WithValue(ctx, UserIDKey, claims.UserID)
	ctxWithUser = context.WithValue(ctxWithUser, UserEmailKey, claims.Email)
	ctxWithUser = context.WithValue(ctxWithUser, UserRoleKey, claims.Role)

	log.Printf("Stream token valid untuk UserID: %s", claims.UserID)

	// 6. Buat wrapper stream dengan context baru
	wrappedStream := &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctxWithUser,
	}

	// 7. Teruskan ke handler dengan stream yang sudah memiliki context user
	return handler(srv, wrappedStream)
}

// wrappedServerStream wraps grpc.ServerStream with new context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the wrapper's context
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// PENJELASAN FILE auth_middleware.go:
// File ini berisi middleware untuk validasi JWT token sebelum request diproses
//
// Constant Context Keys:
// - UserIDKey, UserEmailKey, UserRoleKey: Digunakan untuk menyimpan data user di context
// - Setelah token valid, info user disimpan di context untuk diakses handler
//
// Fungsi AuthInterceptor (Unary RPC):
// - Berjalan sebelum setiap request sampai ke handler
// - Cek apakah method termasuk public (Login, Register, GetMakes, dll) -> bypass
// - Ambil token dari header "authorization" dengan format "Bearer <token>"
// - Validate token dengan utils.ValidateToken()
// - Simpan user_id, email, role ke context jika token valid
// - Handler bisa akses dengan ctx.Value(auth.UserIDKey)
//
// Fungsi StreamAuthInterceptor (Stream RPC):
// - Sama seperti AuthInterceptor tapi untuk streaming RPC
// - Digunakan untuk GetNotifications (server-side streaming)
// - Wrap ServerStream dengan context yang berisi user info
//
// Public Methods (tidak perlu token):
// - /carapp.AuthService/Login dan /Register
// - /carapp.NhtsaDataService/GetMakes dan /GetModelsForMake
// - /carapp.MobilService/ListMobil dan /GetMobil
//
// Flow:
// Client Request -> Interceptor -> Cek Public Method -> Validate Token
// -> Extract Claims -> Simpan ke Context -> Pass ke Handler
