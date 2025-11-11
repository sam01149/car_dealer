# Script untuk menambahkan komentar ke semua file Go

# Notifikasi Service
$notifikasiService = @"

// PENJELASAN FILE notifikasi_service.go:
// File ini menangani streaming notifikasi ke client (server-side streaming)
//
// Fungsi GetNotifications (Server Streaming RPC):
// - Mengambil user_id dari context (diisi oleh auth middleware)
// - Query 50 notifikasi terbaru user dari database
// - Kirim notifikasi satu per satu ke client via stream
// - Client akan menerima notifikasi secara real-time
//
// Flow:
// 1. Client buka stream connection
// 2. Server validate token di middleware
// 3. Server query notifikasi dari DB (ORDER BY created_at DESC LIMIT 50)
// 4. Loop setiap notifikasi dan kirim via stream.Send()
// 5. Jika error atau selesai, tutup stream
//
// Database:
// - read_at bisa NULL (notifikasi belum dibaca)
// - Menggunakan sql.NullTime untuk handle NULL value
"@

# Dashboard Service
$dashboardService = @"

// PENJELASAN FILE dashboard_service.go:
// File ini menyediakan summary data untuk dashboard user
//
// Fungsi GetDashboard:
// - Mengambil user_id dari context (sudah tervalidasi di middleware)
// - Jalankan 4 query secara PARALEL menggunakan goroutine & WaitGroup:
//   1. Total mobil milik user (COUNT dari mobils WHERE owner_id)
//   2. Transaksi aktif (COUNT transaksi_jual status='diproses' + transaksi_rental status='aktif')
//   3. Pendapatan terakhir (SUM total dari transaksi_jual WHERE penjual_id dan status='selesai')
//   4. Notifikasi baru (COUNT notifikasi WHERE user_id dan read_at IS NULL)
// - Gunakan errChan untuk capture error dari goroutine
// - Return DashboardSummary dengan semua data aggregate
//
// Keuntungan parallel query:
// - Lebih cepat daripada query sequential (4 query jadi 1x waktu query terlama)
// - Efisien untuk dashboard yang butuh banyak data
"@

# NHTSA Client
$nhtsaClient = @"

// PENJELASAN FILE nhtsa_client.go:
// File ini menangani komunikasi dengan NHTSA API eksternal
//
// Constant:
// - nhtsaBaseURL: Base URL NHTSA API (https://vpic.nhtsa.dot.gov/api/vehicles)
// - httpClient: HTTP client dengan timeout 10 detik
//
// Struct NhtsaMake & NhtsaModel:
// - Untuk parsing JSON response dari API
// - MakeID/ModelID: ID numerik dari NHTSA
// - MakeName/ModelName: Nama merek/model mobil
//
// Fungsi FetchAllMakes:
// - Request GET ke /getallmakes?format=json
// - Return semua merek mobil yang ada di database NHTSA
// - Parse JSON response ke slice []NhtsaMake
//
// Fungsi FetchModelsForMakeID:
// - Request GET ke /GetModelsForMakeId/{makeID}?format=json
// - Return semua model untuk merek tertentu
// - Parse JSON response ke slice []NhtsaModel
//
// Error handling:
// - Cek HTTP status code (harus 200 OK)
// - Decode JSON response
// - Log jumlah data yang berhasil diambil
"@

# NHTSA Service
$nhtsaService = @"

// PENJELASAN FILE nhtsa_service.go:
// File ini mengelola cache data NHTSA di database lokal
//
// Constant cacheTTL = 24 jam:
// - Data merek/model dari NHTSA di-cache selama 24 jam
// - Setelah 24 jam, cache dianggap expired dan fetch ulang dari API
//
// Fungsi GetMakes:
// - Coba ambil dari cache DB (tabel nhtsa_makes_cache)
// - Jika cache expired atau kosong -> fetch dari NHTSA API
// - Filter hanya 50 merek terkenal (Toyota, Honda, BMW, dll)
// - Simpan ke cache DB di background (goroutine)
// - Return list merek ke client
//
// Fungsi GetModelsForMake:
// - Coba ambil dari cache DB (tabel nhtsa_models_cache)
// - Jika cache expired -> fetch dari NHTSA API
// - Simpan ke cache DB di background
// - Return list model untuk merek tertentu
//
// Helper Functions:
// - getMakesFromCache: Query cache dari DB, cek apakah masih valid
// - saveMakesToCache: Insert/update cache dengan UPSERT (ON CONFLICT)
// - getModelsFromCache: Query model cache untuk brand_id tertentu
// - saveModelsToCache: Insert/update model cache
// - filterPopularMakes: Filter hanya merek terkenal dari response API
//
// Keuntungan caching:
// - Mengurangi beban ke NHTSA API
// - Response lebih cepat untuk user
// - Hemat bandwidth
"@

Write-Host "Menambahkan komentar ke file-file..." -ForegroundColor Cyan

Add-Content -Path "internal\notifikasi\notifikasi_service.go" -Value $notifikasiService
Write-Host "✓ notifikasi_service.go" -ForegroundColor Green

Add-Content -Path "internal\dashboard\dashboard_service.go" -Value $dashboardService
Write-Host "✓ dashboard_service.go" -ForegroundColor Green

Add-Content -Path "internal\NHTSA\nhtsa_client.go" -Value $nhtsaClient
Write-Host "✓ nhtsa_client.go" -ForegroundColor Green

Add-Content -Path "internal\NHTSA\nhtsa_service\nhtsa_service.go" -Value $nhtsaService
Write-Host "✓ nhtsa_service.go" -ForegroundColor Green

Write-Host "`nSelesai! Semua komentar berhasil ditambahkan." -ForegroundColor Green
