# 🧪 CarApp gRPC API Test Script
# PowerShell script untuk testing API dengan mudah

param(
    [string]$Action = "menu",
    [string]$Token = ""
)

$SERVER = "localhost:9090"
$GREEN = "Green"
$YELLOW = "Yellow"
$RED = "Red"
$CYAN = "Cyan"

function Show-Menu {
    Write-Host "`n========================================" -ForegroundColor $CYAN
    Write-Host "🚗 CarApp gRPC API Testing Menu" -ForegroundColor $CYAN
    Write-Host "========================================`n" -ForegroundColor $CYAN
    
    Write-Host "Authentication:" -ForegroundColor $YELLOW
    Write-Host "  1. Register New User"
    Write-Host "  2. Login User"
    Write-Host "  3. Login as Dealer (Admin)`n"
    
    Write-Host "Mobil Service:" -ForegroundColor $YELLOW
    Write-Host "  4. List All Mobil (Public)"
    Write-Host "  5. Get Mobil Detail"
    Write-Host "  6. Create New Mobil (Auth Required)"
    Write-Host "  7. Get Rolls Royce Detail`n"
    
    Write-Host "NHTSA Service:" -ForegroundColor $YELLOW
    Write-Host "  8. Get All Car Brands"
    Write-Host "  9. Get Models for Brand`n"
    
    Write-Host "Transaction Service:" -ForegroundColor $YELLOW
    Write-Host "  10. Buy Mobil (Auth Required)`n"
    
    Write-Host "Dashboard Service:" -ForegroundColor $YELLOW
    Write-Host "  11. Get Dashboard Summary (Auth Required)`n"
    
    Write-Host "Notification Service:" -ForegroundColor $YELLOW
    Write-Host "  12. Get Notifications Stream (Auth Required)`n"
    
    Write-Host "Utilities:" -ForegroundColor $YELLOW
    Write-Host "  13. List All Services"
    Write-Host "  14. Check Server Status"
    Write-Host "  0. Exit`n"
    
    $choice = Read-Host "Pilih menu (0-14)"
    return $choice
}

function Test-ServerConnection {
    Write-Host "`n🔍 Checking server connection..." -ForegroundColor $CYAN
    try {
        $result = grpcurl -plaintext $SERVER list 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Server is running at $SERVER" -ForegroundColor $GREEN
            return $true
        } else {
            Write-Host "❌ Server not responding!" -ForegroundColor $RED
            Write-Host "Please start the backend: go run main.go" -ForegroundColor $YELLOW
            return $false
        }
    } catch {
        Write-Host "❌ grpcurl not found! Install with: choco install grpcurl" -ForegroundColor $RED
        return $false
    }
}

function Register-User {
    Write-Host "`n📝 Register New User" -ForegroundColor $CYAN
    $name = Read-Host "Name"
    $email = Read-Host "Email"
    $password = Read-Host "Password" -AsSecureString
    $passwordText = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
        [Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))
    $phone = Read-Host "Phone (+628...)"
    
    $json = @{
        name = $name
        email = $email
        password = $passwordText
        phone = $phone
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Sending request..." -ForegroundColor $YELLOW
    grpcurl -plaintext -d $json $SERVER carapp.AuthService/Register
}

function Login-User {
    Write-Host "`n🔐 Login User" -ForegroundColor $CYAN
    $email = Read-Host "Email"
    $password = Read-Host "Password" -AsSecureString
    $passwordText = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
        [Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))
    
    $json = @{
        email = $email
        password = $passwordText
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Sending request..." -ForegroundColor $YELLOW
    $response = grpcurl -plaintext -d $json $SERVER carapp.AuthService/Login | ConvertFrom-Json
    
    if ($response.token) {
        Write-Host "`n✅ Login successful!" -ForegroundColor $GREEN
        Write-Host "📋 Token copied to clipboard!" -ForegroundColor $GREEN
        $response.token | Set-Clipboard
        Write-Host "`nUser Info:" -ForegroundColor $CYAN
        Write-Host "  Name: $($response.user.name)"
        Write-Host "  Email: $($response.user.email)"
        Write-Host "  Role: $($response.user.role)"
        Write-Host "  ID: $($response.user.id)"
    }
}

function Login-Dealer {
    Write-Host "`n🔐 Login as Dealer (Admin)" -ForegroundColor $CYAN
    
    $json = @{
        email = "dealer@carapp.com"
        password = "dealer123"
    } | ConvertTo-Json -Compress
    
    Write-Host "🚀 Logging in as dealer..." -ForegroundColor $YELLOW
    $response = grpcurl -plaintext -d $json $SERVER carapp.AuthService/Login | ConvertFrom-Json
    
    if ($response.token) {
        Write-Host "`n✅ Login successful!" -ForegroundColor $GREEN
        Write-Host "📋 Token copied to clipboard!" -ForegroundColor $GREEN
        $response.token | Set-Clipboard
        Write-Host "`nDealer Info:" -ForegroundColor $CYAN
        Write-Host "  Name: $($response.user.name)"
        Write-Host "  Email: $($response.user.email)"
        Write-Host "  Role: $($response.user.role)"
    }
}

function List-Mobil {
    Write-Host "`n🚗 Listing All Mobil (Public)" -ForegroundColor $CYAN
    $page = Read-Host "Page (default 1)"
    if ([string]::IsNullOrEmpty($page)) { $page = 1 }
    $limit = Read-Host "Limit (default 10)"
    if ([string]::IsNullOrEmpty($limit)) { $limit = 10 }
    
    $json = @{
        page = [int]$page
        limit = [int]$limit
        filter_status = "tersedia"
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Fetching mobil list..." -ForegroundColor $YELLOW
    grpcurl -plaintext -d $json $SERVER carapp.MobilService/ListMobil
}

function Get-MobilDetail {
    Write-Host "`n🔍 Get Mobil Detail" -ForegroundColor $CYAN
    $mobilId = Read-Host "Mobil ID"
    
    $json = @{
        mobil_id = $mobilId
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Fetching mobil detail..." -ForegroundColor $YELLOW
    grpcurl -plaintext -d $json $SERVER carapp.MobilService/GetMobil
}

function Create-Mobil {
    Write-Host "`n➕ Create New Mobil (Auth Required)" -ForegroundColor $CYAN
    Write-Host "⚠️  Make sure you have logged in and token is in clipboard!" -ForegroundColor $YELLOW
    
    $token = Get-Clipboard
    if ([string]::IsNullOrEmpty($token)) {
        Write-Host "❌ No token found in clipboard! Please login first." -ForegroundColor $RED
        return
    }
    
    Write-Host "`nEnter mobil details:" -ForegroundColor $CYAN
    $merk = Read-Host "Merk (e.g., Toyota)"
    $model = Read-Host "Model (e.g., Avanza)"
    $tahun = Read-Host "Tahun (e.g., 2023)"
    $kondisi = Read-Host "Kondisi (baru/bekas)"
    $deskripsi = Read-Host "Deskripsi"
    $hargaJual = Read-Host "Harga Jual (IDR)"
    $fotoUrl = Read-Host "Foto URL"
    $lokasi = Read-Host "Lokasi"
    
    $json = @{
        merk = $merk
        model = $model
        tahun = [int]$tahun
        kondisi = $kondisi
        deskripsi = $deskripsi
        harga_jual = [double]$hargaJual
        foto_url = $fotoUrl
        lokasi = $lokasi
        harga_rental_per_hari = 0
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Creating mobil..." -ForegroundColor $YELLOW
    grpcurl -plaintext -H "Authorization: Bearer $token" -d $json $SERVER carapp.MobilService/CreateMobil
}

function Get-RollsRoyceDetail {
    Write-Host "`n🎩 Get Rolls Royce Detail (Example)" -ForegroundColor $CYAN
    Write-Host "This will find and show details of a Rolls Royce mobil" -ForegroundColor $YELLOW
    
    # First, list mobil to find Rolls Royce
    Write-Host "`n🔍 Finding Rolls Royce..." -ForegroundColor $YELLOW
    $json = '{"page":1,"limit":50}'
    $response = grpcurl -plaintext -d $json $SERVER carapp.MobilService/ListMobil | ConvertFrom-Json
    
    $rollsRoyce = $response.mobils | Where-Object { $_.merk -like "*Rolls*" } | Select-Object -First 1
    
    if ($rollsRoyce) {
        Write-Host "✅ Found Rolls Royce: $($rollsRoyce.merk) $($rollsRoyce.model)" -ForegroundColor $GREEN
        Write-Host "`n📋 Details:" -ForegroundColor $CYAN
        $rollsRoyce | ConvertTo-Json -Depth 10
    } else {
        Write-Host "❌ No Rolls Royce found in database" -ForegroundColor $RED
    }
}

function Get-AllBrands {
    Write-Host "`n🏭 Get All Car Brands (NHTSA)" -ForegroundColor $CYAN
    Write-Host "🚀 Fetching brands from cache/API..." -ForegroundColor $YELLOW
    grpcurl -plaintext -d '{}' $SERVER carapp.NhtsaDataService/GetMakes
}

function Get-ModelsForBrand {
    Write-Host "`n🔍 Get Models for Brand" -ForegroundColor $CYAN
    Write-Host "Popular brands:" -ForegroundColor $YELLOW
    Write-Host "  440 - Audi"
    Write-Host "  441 - BMW"
    Write-Host "  445 - Ferrari"
    Write-Host "  493 - Rolls-Royce"
    Write-Host "  460 - Lamborghini"
    
    $brandId = Read-Host "`nBrand ID"
    
    $json = @{
        brand_id = $brandId
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Fetching models..." -ForegroundColor $YELLOW
    grpcurl -plaintext -d $json $SERVER carapp.NhtsaDataService/GetModelsForMake
}

function Buy-Mobil {
    Write-Host "`n💰 Buy Mobil (Auth Required)" -ForegroundColor $CYAN
    Write-Host "⚠️  Make sure you have logged in and token is in clipboard!" -ForegroundColor $YELLOW
    
    $token = Get-Clipboard
    if ([string]::IsNullOrEmpty($token)) {
        Write-Host "❌ No token found in clipboard! Please login first." -ForegroundColor $RED
        return
    }
    
    $mobilId = Read-Host "Mobil ID to buy"
    
    $json = @{
        mobil_id = $mobilId
    } | ConvertTo-Json -Compress
    
    Write-Host "`n🚀 Processing purchase..." -ForegroundColor $YELLOW
    grpcurl -plaintext -H "Authorization: Bearer $token" -d $json $SERVER carapp.TransaksiService/BuyMobil
}

function Get-Dashboard {
    Write-Host "`n📊 Get Dashboard Summary (Auth Required)" -ForegroundColor $CYAN
    Write-Host "⚠️  Make sure you have logged in and token is in clipboard!" -ForegroundColor $YELLOW
    
    $token = Get-Clipboard
    if ([string]::IsNullOrEmpty($token)) {
        Write-Host "❌ No token found in clipboard! Please login first." -ForegroundColor $RED
        return
    }
    
    Write-Host "`n🚀 Fetching dashboard data..." -ForegroundColor $YELLOW
    grpcurl -plaintext -H "Authorization: Bearer $token" -d '{}' $SERVER carapp.DashboardService/GetDashboard
}

function Get-Notifications {
    Write-Host "`n🔔 Get Notifications Stream (Auth Required)" -ForegroundColor $CYAN
    Write-Host "⚠️  This is a streaming endpoint - press Ctrl+C to stop" -ForegroundColor $YELLOW
    Write-Host "⚠️  Make sure you have logged in and token is in clipboard!" -ForegroundColor $YELLOW
    
    $token = Get-Clipboard
    if ([string]::IsNullOrEmpty($token)) {
        Write-Host "❌ No token found in clipboard! Please login first." -ForegroundColor $RED
        return
    }
    
    Write-Host "`n🚀 Starting notification stream..." -ForegroundColor $YELLOW
    grpcurl -plaintext -H "Authorization: Bearer $token" -d '{}' $SERVER carapp.NotifikasiService/GetNotifications
}

function List-Services {
    Write-Host "`n📋 List All gRPC Services" -ForegroundColor $CYAN
    Write-Host "🚀 Fetching services..." -ForegroundColor $YELLOW
    grpcurl -plaintext $SERVER list
    
    Write-Host "`n💡 To see methods, use: grpcurl -plaintext $SERVER list carapp.ServiceName" -ForegroundColor $YELLOW
}

# Main execution
Clear-Host

if (-not (Test-ServerConnection)) {
    exit
}

do {
    $choice = Show-Menu
    
    switch ($choice) {
        "1" { Register-User }
        "2" { Login-User }
        "3" { Login-Dealer }
        "4" { List-Mobil }
        "5" { Get-MobilDetail }
        "6" { Create-Mobil }
        "7" { Get-RollsRoyceDetail }
        "8" { Get-AllBrands }
        "9" { Get-ModelsForBrand }
        "10" { Buy-Mobil }
        "11" { Get-Dashboard }
        "12" { Get-Notifications }
        "13" { List-Services }
        "14" { Test-ServerConnection }
        "0" { 
            Write-Host "`n👋 Goodbye!" -ForegroundColor $CYAN
            exit 
        }
        default { 
            Write-Host "`n❌ Invalid choice!" -ForegroundColor $RED 
        }
    }
    
    Write-Host "`nPress any key to continue..." -ForegroundColor $YELLOW
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    Clear-Host
    
} while ($true)
