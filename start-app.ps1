# Script untuk menjalankan aplikasi Car Dealer
param(
    [string]$Action = "start"
)

$ProjectRoot = "c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc"
$FrontendPath = Join-Path $ProjectRoot "car_dealer_frontend"

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "  Car Dealer App Runner" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

function Start-Backend {
    Write-Host "🚀 Starting Go Backend Server..." -ForegroundColor Yellow
    Write-Host "Location: $ProjectRoot" -ForegroundColor Gray
    Write-Host ""
    
    # Kill existing process
    Write-Host "Checking for existing Go processes..." -ForegroundColor Gray
    Get-Process | Where-Object { $_.ProcessName -eq 'go' -or $_.ProcessName -eq 'main' } | Stop-Process -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 1
    
    # Start new process
    Set-Location $ProjectRoot
    Write-Host "Starting server..." -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop the backend server" -ForegroundColor Yellow
    Write-Host ""
    go run main.go
}

function Start-Frontend {
    Write-Host "🚀 Starting Next.js Frontend..." -ForegroundColor Yellow
    Write-Host "Location: $FrontendPath" -ForegroundColor Gray
    Write-Host ""
    
    Set-Location $FrontendPath
    Write-Host "Starting server..." -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop the frontend server" -ForegroundColor Yellow
    Write-Host ""
    npm run dev
}

function Stop-Servers {
    Write-Host "🛑 Stopping all servers..." -ForegroundColor Yellow
    
    # Stop Go
    Write-Host "Stopping Go backend..." -ForegroundColor Gray
    Get-Process | Where-Object { $_.ProcessName -eq 'go' -or $_.ProcessName -eq 'main' } | Stop-Process -Force -ErrorAction SilentlyContinue
    
    # Stop Node (port 3000)
    Write-Host "Stopping Next.js frontend..." -ForegroundColor Gray
    Get-NetTCPConnection -LocalPort 3000 -ErrorAction SilentlyContinue | ForEach-Object { 
        Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
    }
    
    # Stop port 9090
    Get-NetTCPConnection -LocalPort 9090 -ErrorAction SilentlyContinue | ForEach-Object { 
        Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host "✅ All servers stopped!" -ForegroundColor Green
}

function Show-Status {
    Write-Host "📊 Server Status:" -ForegroundColor Yellow
    Write-Host ""
    
    # Check port 9090 (Backend)
    $backend = Get-NetTCPConnection -LocalPort 9090 -ErrorAction SilentlyContinue
    if ($backend) {
        Write-Host "✅ Backend (Go gRPC):     Running on port 9090" -ForegroundColor Green
    } else {
        Write-Host "❌ Backend (Go gRPC):     Not running" -ForegroundColor Red
    }
    
    # Check port 3000 (Frontend)
    $frontend = Get-NetTCPConnection -LocalPort 3000 -ErrorAction SilentlyContinue
    if ($frontend) {
        Write-Host "✅ Frontend (Next.js):    Running on port 3000" -ForegroundColor Green
        Write-Host "   URL: http://localhost:3000" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Frontend (Next.js):    Not running" -ForegroundColor Red
    }
    
    Write-Host ""
    
    # Show process details
    $processes = Get-Process | Where-Object { 
        $_.ProcessName -eq 'go' -or 
        $_.ProcessName -eq 'main' -or 
        $_.ProcessName -eq 'node'
    } | Select-Object Id, ProcessName, @{Name="Memory(MB)";Expression={[math]::Round($_.WorkingSet / 1MB, 2)}}
    
    if ($processes) {
        Write-Host "📋 Related Processes:" -ForegroundColor Yellow
        $processes | Format-Table -AutoSize
    }
}

switch ($Action) {
    "backend" {
        Start-Backend
    }
    
    "frontend" {
        Start-Frontend
    }
    
    "stop" {
        Stop-Servers
    }
    
    "status" {
        Show-Status
    }
    
    "restart" {
        Write-Host "♻️  Restarting servers..." -ForegroundColor Yellow
        Stop-Servers
        Start-Sleep -Seconds 2
        Write-Host ""
        Write-Host "Please run:" -ForegroundColor Cyan
        Write-Host "  Terminal 1: .\start-app.ps1 -Action backend" -ForegroundColor White
        Write-Host "  Terminal 2: .\start-app.ps1 -Action frontend" -ForegroundColor White
    }
    
    default {
        Write-Host "Usage: .\start-app.ps1 -Action <action>" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Available actions:" -ForegroundColor Cyan
        Write-Host "  backend   - Start Go backend server (port 9090)" -ForegroundColor White
        Write-Host "  frontend  - Start Next.js frontend (port 3000)" -ForegroundColor White
        Write-Host "  stop      - Stop all servers" -ForegroundColor White
        Write-Host "  status    - Check server status" -ForegroundColor White
        Write-Host "  restart   - Restart all servers" -ForegroundColor White
        Write-Host ""
        Write-Host "Quick Start:" -ForegroundColor Cyan
        Write-Host "  1. Open Terminal 1:" -ForegroundColor Gray
        Write-Host "     .\start-app.ps1 -Action backend" -ForegroundColor White
        Write-Host ""
        Write-Host "  2. Open Terminal 2:" -ForegroundColor Gray
        Write-Host "     .\start-app.ps1 -Action frontend" -ForegroundColor White
        Write-Host ""
        Write-Host "  3. Open Browser:" -ForegroundColor Gray
        Write-Host "     http://localhost:3000" -ForegroundColor White
        Write-Host ""
        Write-Host "Check status:" -ForegroundColor Cyan
        Write-Host "  .\start-app.ps1 -Action status" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
