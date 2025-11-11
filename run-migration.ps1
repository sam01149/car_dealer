# Script untuk menjalankan migration 002
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  RUNNING MIGRATION 002" -ForegroundColor Yellow
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

Write-Host "🔄 Running migration 002..." -ForegroundColor Yellow

# Set password environment variable
$env:PGPASSWORD = $dbPass

# Jalankan migration
& psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f "db\migrations\migrations\002_remove_rental_add_foto.up.sql"

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host "✅ MIGRATION SELESAI!" -ForegroundColor Green
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "📝 Perubahan:" -ForegroundColor White
    Write-Host "   • Menghapus tabel transaksi_rental" -ForegroundColor Gray
    Write-Host "   • Menghapus kolom harga_rental_per_hari dari mobils" -ForegroundColor Gray
    Write-Host "   • Menambahkan kolom foto_url ke mobils" -ForegroundColor Gray
    Write-Host "   • Update status mobil dari 'dirental' ke 'tersedia'" -ForegroundColor Gray
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "❌ Migration gagal" -ForegroundColor Red
    Write-Host "💡 Pastikan PostgreSQL running dan kredensial benar" -ForegroundColor Yellow
    exit 1
}
