package main

import (
	"log"
	"net"
	"os"

	"carapp.com/m/internal/NHTSA/nhtsa_service"
	"carapp.com/m/internal/auth"
	"carapp.com/m/internal/dashboard"
	"carapp.com/m/internal/db"
	"carapp.com/m/internal/mobil"
	"carapp.com/m/internal/notifikasi"
	"carapp.com/m/internal/transaksi"
	pb "carapp.com/m/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// Import package baru untuk interceptor
)

// PERTAMA, install dependency baru:
// go get github.com/grpc-ecosystem/go-grpc-middleware

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan...")
	}

	// 2. Koneksi ke Database
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	// 3. Setup Server gRPC
	grpcPort := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcPort == "" {
		grpcPort = "0.0.0.0:9090"
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Gagal listen: %v", err)
	}

	// Terapkan Unary Interceptor (Middleware)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor),
	)

	// 4. Register Services
	// 4. Register Services
	authServer := auth.NewAuthService(dbConn)
	pb.RegisterAuthServiceServer(s, authServer)

	mobilServer := mobil.NewMobilService(dbConn)
	pb.RegisterMobilServiceServer(s, mobilServer)

	nhtsaServer := nhtsa_service.NewNhtsaDataService(dbConn)
	pb.RegisterNhtsaDataServiceServer(s, nhtsaServer)

	transaksiServer := transaksi.NewTransaksiService(dbConn)
	pb.RegisterTransaksiServiceServer(s, transaksiServer)

	// --- DAFTARKAN SERVICE BARU KITA ---
	notifikasiServer := notifikasi.NewNotifikasiService(dbConn)
	pb.RegisterNotifikasiServiceServer(s, notifikasiServer)
	// ------------------------------------
	dashboardServer := dashboard.NewDashboardService(dbConn)
	pb.RegisterDashboardServiceServer(s, dashboardServer)

	reflection.Register(s)

	log.Printf("Server gRPC berjalan di %s...", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
