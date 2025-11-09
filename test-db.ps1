# Script untuk cek database
param(
    [string]$Action = "check"
)

$env:PGPASSWORD = '123456'
$PSQL = "psql"  # Atau path lengkap: "C:\Program Files\PostgreSQL\16\bin\psql.exe"

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "  Car Dealer DB Testing Tool" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

switch ($Action) {
    "check" {
        Write-Host "📊 Checking database status..." -ForegroundColor Yellow
        Write-Host ""
        
        Write-Host "👥 Users:" -ForegroundColor Green
        & $PSQL -U postgres -d car_db -c "SELECT id, email, created_at FROM users ORDER BY created_at DESC LIMIT 5;"
        Write-Host ""
        
        Write-Host "🚗 Mobils:" -ForegroundColor Green
        & $PSQL -U postgres -d car_db -c "SELECT id, merk, model, tahun, status, harga_jual, harga_rental_per_hari, owner_id FROM mobils ORDER BY created_at DESC LIMIT 5;"
        Write-Host ""
        
        Write-Host "💰 Transaksis:" -ForegroundColor Green
        & $PSQL -U postgres -d car_db -c "SELECT id, tipe_transaksi, status, total_harga, buyer_id, seller_id FROM transaksis ORDER BY created_at DESC LIMIT 5;"
        Write-Host ""
        
        Write-Host "🔔 Notifikasis:" -ForegroundColor Green
        & $PSQL -U postgres -d car_db -c "SELECT id, user_id, pesan, is_read FROM notifikasis ORDER BY created_at DESC LIMIT 5;"
    }
    
    "mobils" {
        Write-Host "🚗 All Mobils with Status:" -ForegroundColor Yellow
        & $PSQL -U postgres -d car_db -c "SELECT id, merk, model, tahun, status, harga_jual, owner_id FROM mobils ORDER BY created_at DESC;"
    }
    
    "available" {
        Write-Host "🚗 Available Mobils (status='tersedia'):" -ForegroundColor Yellow
        & $PSQL -U postgres -d car_db -c "SELECT id, merk, model, tahun, status, harga_jual, owner_id FROM mobils WHERE status='tersedia' ORDER BY created_at DESC;"
    }
    
    "users" {
        Write-Host "👥 All Users:" -ForegroundColor Yellow
        & $PSQL -U postgres -d car_db -c "SELECT id, email, nama, created_at FROM users ORDER BY created_at DESC;"
    }
    
    "transactions" {
        Write-Host "💰 All Transactions:" -ForegroundColor Yellow
        & $PSQL -U postgres -d car_db -c "SELECT t.id, t.tipe_transaksi, t.status, t.total_harga, m.merk, m.model, u1.email as buyer, u2.email as seller FROM transaksis t LEFT JOIN mobils m ON t.mobil_id = m.id LEFT JOIN users u1 ON t.buyer_id = u1.id LEFT JOIN users u2 ON t.seller_id = u2.id ORDER BY t.created_at DESC;"
    }
    
    "reset" {
        Write-Host "⚠️  WARNING: This will delete ALL data!" -ForegroundColor Red
        $confirm = Read-Host "Type 'YES' to confirm"
        if ($confirm -eq "YES") {
            Write-Host "🗑️  Resetting database..." -ForegroundColor Red
            & $PSQL -U postgres -d car_db -c "TRUNCATE TABLE mobils, transaksis, notifikasis RESTART IDENTITY CASCADE;"
            Write-Host "✅ Database reset complete!" -ForegroundColor Green
        } else {
            Write-Host "❌ Reset cancelled." -ForegroundColor Yellow
        }
    }
    
    "seed" {
        Write-Host "🌱 Seeding test data..." -ForegroundColor Yellow
        Write-Host "Note: Make sure you have at least 2 users registered first!" -ForegroundColor Cyan
        Write-Host "This will NOT create users, only mobils." -ForegroundColor Cyan
        Write-Host ""
        
        # Get first user
        $userId = & $PSQL -U postgres -d car_db -t -c "SELECT id FROM users LIMIT 1;"
        $userId = $userId.Trim()
        
        if ($userId) {
            Write-Host "Adding test mobil for user: $userId" -ForegroundColor Green
            & $PSQL -U postgres -d car_db -c @"
INSERT INTO mobils (merk, model, tahun, kondisi, harga_jual, harga_rental_per_hari, lokasi, deskripsi, status, owner_id) 
VALUES 
    ('Toyota', 'Camry', 2023, 'baru', 500000000, 1000000, 'Jakarta', 'Mobil sedan mewah, kondisi sangat baik', 'tersedia', '$userId'),
    ('Honda', 'Civic', 2022, 'bekas', 350000000, 750000, 'Bandung', 'Mobil sport compact, terawat', 'tersedia', '$userId'),
    ('Suzuki', 'Ertiga', 2021, 'bekas', 200000000, 500000, 'Surabaya', 'MPV keluarga, 7 seater', 'tersedia', '$userId');
"@
            Write-Host "✅ Test data added!" -ForegroundColor Green
        } else {
            Write-Host "❌ No users found. Please register users first!" -ForegroundColor Red
        }
    }
    
    default {
        Write-Host "Usage: .\test-db.ps1 -Action <action>" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Available actions:" -ForegroundColor Cyan
        Write-Host "  check        - Show summary of all tables (default)" -ForegroundColor White
        Write-Host "  mobils       - Show all mobils" -ForegroundColor White
        Write-Host "  available    - Show available mobils only" -ForegroundColor White
        Write-Host "  users        - Show all users" -ForegroundColor White
        Write-Host "  transactions - Show all transactions with details" -ForegroundColor White
        Write-Host "  seed         - Add test mobils (requires existing users)" -ForegroundColor White
        Write-Host "  reset        - Reset ALL data (dangerous!)" -ForegroundColor White
        Write-Host ""
        Write-Host "Examples:" -ForegroundColor Cyan
        Write-Host "  .\test-db.ps1" -ForegroundColor Gray
        Write-Host "  .\test-db.ps1 -Action available" -ForegroundColor Gray
        Write-Host "  .\test-db.ps1 -Action seed" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
