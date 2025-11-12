# Setup Protocol Buffers Compiler dan Regenerate Proto
Write-Host "=== Setup Proto Compiler ===" -ForegroundColor Cyan

# Check apakah protoc sudah ada
$protocExists = Get-Command protoc -ErrorAction SilentlyContinue

if (-not $protocExists) {
    Write-Host "`n[ERROR] protoc belum terinstall" -ForegroundColor Red
    Write-Host "`nCara install: choco install protoc" -ForegroundColor Yellow
    exit 1
}

Write-Host "[OK] protoc sudah terinstall" -ForegroundColor Green
protoc --version

# Install Go plugins
Write-Host "`n=== Install Go Proto Plugins ===" -ForegroundColor Cyan
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
Write-Host "[OK] Go plugins installed" -ForegroundColor Green

# Regenerate proto files
Write-Host "`n=== Regenerate Proto Files ===" -ForegroundColor Cyan
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/carapp.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "[OK] Proto files regenerated!" -ForegroundColor Green
} else {
    Write-Host "[ERROR] Failed to regenerate proto" -ForegroundColor Red
    exit 1
}

# Install uuid dependency
Write-Host "`n=== Install Go Dependencies ===" -ForegroundColor Cyan
go get github.com/google/uuid
go mod tidy
Write-Host "[OK] Dependencies installed" -ForegroundColor Green

# Create uploads folder
Write-Host "`n=== Create Uploads Folder ===" -ForegroundColor Cyan
if (-not (Test-Path "uploads")) {
    New-Item -ItemType Directory -Path "uploads" | Out-Null
    Write-Host "[OK] Folder 'uploads' created" -ForegroundColor Green
} else {
    Write-Host "[OK] Folder 'uploads' already exists" -ForegroundColor Green
}

Write-Host "`n[SUCCESS] Setup complete!" -ForegroundColor Green
Write-Host "Next: .\restart-all.ps1" -ForegroundColor Yellow
