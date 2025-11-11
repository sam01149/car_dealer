# Script untuk restart backend dan frontend

Write-Host "🔄 Restarting CarApp..." -ForegroundColor Cyan
Write-Host ""

# Kill proses yang sedang berjalan
Write-Host "⏹️  Stopping running processes..." -ForegroundColor Yellow
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process -Name "node" -ErrorAction SilentlyContinue | Where-Object {$_.Path -like "*car_dealer_frontend*"} | Stop-Process -Force

Start-Sleep -Seconds 2

Write-Host "✅ Processes stopped" -ForegroundColor Green
Write-Host ""

# Start backend di background
Write-Host "🚀 Starting Backend (Go gRPC Server)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot'; go run main.go"

Start-Sleep -Seconds 3

# Start frontend di background  
Write-Host "🚀 Starting Frontend (Next.js)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot\car_dealer_frontend'; npm run dev"

Write-Host ""
Write-Host "✅ SELESAI! Aplikasi sedang starting..." -ForegroundColor Green
Write-Host ""
Write-Host "📝 Info:" -ForegroundColor Yellow
Write-Host "   - Backend: http://localhost:9090" -ForegroundColor White
Write-Host "   - Frontend: http://localhost:3000" -ForegroundColor White
Write-Host ""
Write-Host "⏳ Tunggu 10-15 detik, lalu buka http://localhost:3000" -ForegroundColor Cyan
Write-Host ""
