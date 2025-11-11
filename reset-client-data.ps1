# Script untuk reset semua data client dari database
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  RESET DATA CLIENT" -ForegroundColor Yellow
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "⚠️  PERINGATAN: Ini akan menghapus semua data client!" -ForegroundColor Red
Write-Host "   - Semua user client (kecuali dealer)" -ForegroundColor Yellow
Write-Host "   - Semua mobil dari client" -ForegroundColor Yellow
Write-Host "   - Semua transaksi" -ForegroundColor Yellow
Write-Host "   - Semua notifikasi client" -ForegroundColor Yellow
Write-Host ""
Write-Host "   Yang TIDAK dihapus:" -ForegroundColor Green
Write-Host "   ✓ User dealer (dealer@carapp.com)" -ForegroundColor Gray
Write-Host "   ✓ 42 mobil dealer dengan foto" -ForegroundColor Gray
Write-Host ""

$confirm = Read-Host "Ketik 'RESET' untuk lanjut"

if ($confirm -ne "RESET") {
    Write-Host "❌ Reset dibatalkan" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "🔄 Menjalankan reset..." -ForegroundColor Yellow

# Jalankan reset script
go run cmd/reset-client/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host "✅ RESET BERHASIL!" -ForegroundColor Green
    Write-Host "============================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "🎯 Database sekarang bersih dan siap untuk testing!" -ForegroundColor Green
    Write-Host ""
    Write-Host "💡 Langkah selanjutnya:" -ForegroundColor White
    Write-Host "   1. Buka http://localhost:3000" -ForegroundColor Gray
    Write-Host "   2. Register user baru" -ForegroundColor Gray
    Write-Host "   3. Test fitur jual mobil" -ForegroundColor Gray
    Write-Host "   4. Test fitur beli mobil" -ForegroundColor Gray
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "❌ Reset gagal" -ForegroundColor Red
    exit 1
}
