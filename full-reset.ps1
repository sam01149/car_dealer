# Full reset database (DROP dan CREATE ulang)
Write-Host "⚠️  FULL DATABASE RESET" -ForegroundColor Red
Write-Host ""

# Load DB info
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^DB_SOURCE=(.*)$') {
            $dbSource = $matches[1] -replace '"', ''
            
            if ($dbSource -match 'postgresql://([^:]+):([^@]+)@([^:]+):(\d+)/([^\?]+)') {
                $env:PGUSER = $matches[1]
                $env:PGPASSWORD = $matches[2]
                $env:PGHOST = $matches[3]
                $env:PGPORT = $matches[4]
                $dbName = $matches[5]
            }
        }
    }
}

$confirm = Read-Host "Ketik nama database '$dbName' untuk konfirmasi"

if ($confirm -ne $dbName) {
    Write-Host "❌ Dibatalkan" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "🗑️  Dropping database..." -ForegroundColor Yellow

# Drop database
& psql -d postgres -c "DROP DATABASE IF EXISTS $dbName;"
Write-Host "✅ Database dropped" -ForegroundColor Green

# Create database
& psql -d postgres -c "CREATE DATABASE $dbName;"
Write-Host "✅ Database created" -ForegroundColor Green

Write-Host ""
Write-Host "📋 Running migrations..." -ForegroundColor Yellow

# Run migration
& psql -d $dbName -f "db/migrations/migrations/001_init_schema.up.sql"
Write-Host "✅ Migration completed" -ForegroundColor Green

Write-Host ""
Write-Host "🌱 Running seeder..." -ForegroundColor Yellow

# Run seeder
go run cmd/seeder/main.go

Write-Host ""
Write-Host "✅ FULL RESET COMPLETE!" -ForegroundColor Green
