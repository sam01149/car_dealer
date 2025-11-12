package main

import (
	"log"
	"net/http"
	"os"

	"carapp.com/m/internal/auth"
	"carapp.com/m/internal/dashboard"
	"carapp.com/m/internal/db"
	"carapp.com/m/internal/mobil"
	"carapp.com/m/internal/nhtsa/nhtsa_service"
	"carapp.com/m/internal/notifikasi"
	"carapp.com/m/internal/transaksi"
	pb "carapp.com/m/proto"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan...")
	}

	// 2. Koneksi ke Database
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	// 3. Buat server gRPC dengan UnaryInterceptor dan StreamInterceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor),
		grpc.StreamInterceptor(auth.StreamAuthInterceptor),
	)

	// 4. Register Services
	authServer := auth.NewAuthService(dbConn)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	mobilServer := mobil.NewMobilService(dbConn)
	pb.RegisterMobilServiceServer(grpcServer, mobilServer)

	nhtsaServer := nhtsa_service.NewNhtsaDataService(dbConn)
	pb.RegisterNhtsaDataServiceServer(grpcServer, nhtsaServer)

	transaksiServer := transaksi.NewTransaksiService(dbConn)
	pb.RegisterTransaksiServiceServer(grpcServer, transaksiServer)

	notifikasiServer := notifikasi.NewNotifikasiService(dbConn)
	pb.RegisterNotifikasiServiceServer(grpcServer, notifikasiServer)

	dashboardServer := dashboard.NewDashboardService(dbConn)
	pb.RegisterDashboardServiceServer(grpcServer, dashboardServer)

	reflection.Register(grpcServer)

	// 5. Buat wrapper gRPC-Web
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	// 6. Buat Handler HTTP dengan CORS
	// File server untuk uploads
	fileServer := http.FileServer(http.Dir("./uploads"))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://10.0.7.129:3000", "http://127.0.0.1:3000", "http://172.18.208.1:3000"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "X-Grpc-Web", "X-User-Agent", "Authorization", "authorization"},
		ExposedHeaders:   []string{"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Serve static files dari /uploads/
		if len(r.URL.Path) > 9 && r.URL.Path[:9] == "/uploads/" {
			log.Printf("Serving static file: %s", r.URL.Path)
			http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, r)
			return
		}

		if wrappedGrpc.IsAcceptableGrpcCorsRequest(r) || wrappedGrpc.IsGrpcWebRequest(r) {
			wrappedGrpc.ServeHTTP(w, r)
			return
		}
		log.Printf("Not a gRPC-Web request: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	}))

	// 7. Jalankan Server HTTP
	grpcPort := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcPort == "" {
		grpcPort = "0.0.0.0:9090"
	}

	httpServer := &http.Server{
		Addr:    grpcPort,
		Handler: corsHandler,
	}

	log.Printf("Server gRPC-Web (HTTP) berjalan di %s...", grpcPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}

// PENJELASAN FILE main.go:
// File ini adalah entry point aplikasi backend Car Dealer gRPC server
// Fungsi utama:
// - Load konfigurasi dari .env (DB_SOURCE, JWT_SECRET_KEY, dll)
// - Buat koneksi ke database PostgreSQL
// - Setup gRPC server dengan middleware autentikasi (UnaryInterceptor & StreamInterceptor)
// - Daftarkan semua service: Auth, Mobil, NHTSA, Transaksi, Notifikasi, Dashboard
// - Wrap gRPC dengan grpc-web agar bisa diakses browser
// - Setup CORS untuk mengizinkan frontend mengakses API
// - Jalankan HTTP server di port 9090 (default)
