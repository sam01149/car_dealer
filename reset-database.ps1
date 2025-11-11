# Script untuk reset database (seperti migrate:fresh --seed di Laravel)
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  DATABASE RESET & SEED" -ForegroundColor Yellow
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# 1. Load .env untuk ambil DB_SOURCE
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^DB_SOURCE=(.*)$') {
            $env:DB_SOURCE = $matches[1] -replace '"', ''
        }
    }
}

if (-not $env:DB_SOURCE) {
    Write-Host "❌ DB_SOURCE tidak ditemukan di .env" -ForegroundColor Red
    exit 1
}

# Extract info dari connection string
if ($env:DB_SOURCE -match 'postgresql://([^:]+):([^@]+)@([^:]+):(\d+)/([^\?]+)') {
    $dbUser = $matches[1]
    $dbPass = $matches[2]
    $dbHost = $matches[3]
    $dbPort = $matches[4]
    $dbName = $matches[5]
    
    Write-Host "📊 Database Info:" -ForegroundColor White
    Write-Host "   Host: $dbHost" -ForegroundColor Gray
    Write-Host "   Port: $dbPort" -ForegroundColor Gray
    Write-Host "   DB: $dbName" -ForegroundColor Gray
    Write-Host "   User: $dbUser" -ForegroundColor Gray
    Write-Host ""
}

# 2. Konfirmasi dari user
Write-Host "⚠️  PERINGATAN: Ini akan menghapus SEMUA data!" -ForegroundColor Red
$confirm = Read-Host "Ketik 'YES' untuk lanjut"

if ($confirm -ne "YES") {
    Write-Host "❌ Reset dibatalkan" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "🗑️  Step 1: Menghapus semua data..." -ForegroundColor Yellow

# 3. Truncate semua tabel
$truncateSQL = @"
TRUNCATE TABLE notifikasi CASCADE;
TRUNCATE TABLE transaksi_rental CASCADE;
TRUNCATE TABLE transaksi_jual CASCADE;
TRUNCATE TABLE mobils CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE nhtsa_models_cache CASCADE;
TRUNCATE TABLE nhtsa_makes_cache CASCADE;
"@

# Simpan ke temp file
$truncateSQL | Out-File -FilePath "temp_truncate.sql" -Encoding UTF8

# Jalankan truncate
$env:PGPASSWORD = $dbPass
& psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f temp_truncate.sql 2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✅ Data berhasil dihapus" -ForegroundColor Green
} else {
    Write-Host "   ❌ Gagal menghapus data" -ForegroundColor Red
    Write-Host "   💡 Pastikan PostgreSQL running dan kredensial benar" -ForegroundColor Yellow
    Remove-Item temp_truncate.sql -ErrorAction SilentlyContinue
    exit 1
}

Remove-Item temp_truncate.sql -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "🌱 Step 2: Menjalankan seeder..." -ForegroundColor Yellow

# 4. Jalankan seeder
go run cmd/seeder/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host "✅ RESET & SEED SELESAI!" -ForegroundColor Green
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "📝 Data yang tersedia:" -ForegroundColor White
    Write-Host "   • User Dealer (Email: dealer@carapp.com, Pass: dealer123)" -ForegroundColor Gray
    Write-Host "   • 50 mobil dari Marketcheck API" -ForegroundColor Gray
    Write-Host ""
    Write-Host "🚀 Sekarang bisa test aplikasi dari awal!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Seeder gagal dijalankan" -ForegroundColor Red
}
