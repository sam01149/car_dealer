# 🚀 Quick Start Guide - Car Dealer App

## 📋 Prerequisites Checklist
- ✅ PostgreSQL running (port 5432)
- ✅ Database `car_db` created
- ✅ Go 1.25+ installed
- ✅ Node.js 18+ installed
- ✅ Marketcheck API key configured

---

## 🎬 First Time Setup

### 1. Clone & Setup (Already Done)
```powershell
cd c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc
```

### 2. Run Database Seeder (One Time)
```powershell
# This will create dealer user and add 47 real cars
go run ./cmd/seeder/main.go
```

**Expected Output:**
```
✅ SELESAI! Berhasil menyimpan 47 mobil baru
```

### 3. Start Backend Server
```powershell
# Terminal 1
go run main.go
```

**Expected Output:**
```
2025/11/09 15:45:00 Berhasil terhubung ke database PostgreSQL!
2025/11/09 15:45:00 Server gRPC-Web (HTTP) berjalan di 0.0.0.0:9090...
```

### 4. Start Frontend Server
```powershell
# Terminal 2
cd car_dealer_frontend
npm run dev
```

**Expected Output:**
```
- ready started server on 0.0.0.0:3000, url: http://localhost:3000
```

### 5. Open Browser
```
http://localhost:3000
```

**You should see 47 cars on the homepage!** 🎉

---

## 🧪 Quick Testing Scenarios

### Scenario 1: Browse Cars (No Login Required)
1. Open http://localhost:3000
2. ✅ See 47 cars listed
3. Click any car
4. ✅ See detailed information
5. ❌ No buy/rental buttons (not logged in)

### Scenario 2: Register & Buy Car
1. Click "Register" → Create account
   - Email: `buyer1@test.com`
   - Password: `test123`
2. Go to homepage
3. Click any car (NOT owned by you)
4. ✅ **"Beli Mobil Ini" button appears!**
5. Click "Beli Mobil Ini"
6. ✅ Success message appears
7. ✅ Auto-redirect to dashboard
8. ✅ See transaction in dashboard

### Scenario 3: Rental Car
1. Login as buyer
2. Browse cars
3. Click "Rental Mobil Ini"
4. Select start & end date
5. Click "Konfirmasi Rental"
6. ✅ Rental successful!

### Scenario 4: Sell Your Own Car
1. Login
2. Click "Jual Mobil" in navbar
3. Fill form:
   - Merk: Toyota
   - Model: Avanza
   - Tahun: 2023
   - Kondisi: baru
   - Harga Jual: 250000000
   - Harga Rental: 500000
   - Lokasi: Jakarta
   - Deskripsi: Mobil keluarga nyaman
4. Submit
5. ✅ Car added to your inventory
6. Go to dashboard
7. ✅ See your car in "Mobil yang Dijual"

### Scenario 5: Login as Dealer
1. Logout
2. Login with:
   - Email: `dealer@carapp.com`
   - Password: `dealer123`
3. Go to dashboard
4. ✅ See all 47 cars in inventory!

---

## 🛠️ Useful Commands

### Check Server Status
```powershell
.\start-app.ps1 -Action status
```

### Stop All Servers
```powershell
.\start-app.ps1 -Action stop
```

### Check Database
```powershell
# See all available cars
.\test-db.ps1 -Action available

# See all data
.\test-db.ps1 -Action check

# See transactions
.\test-db.ps1 -Action transactions
```

### Add More Cars
```powershell
# Run seeder again (duplicates will be skipped)
go run ./cmd/seeder/main.go
```

### Reset Cars (Danger!)
```powershell
# Delete all cars from dealer
$env:PGPASSWORD='123456'
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d car_db -c "DELETE FROM mobils WHERE owner_id = (SELECT id FROM users WHERE email = 'dealer@carapp.com');"

# Re-run seeder
go run ./cmd/seeder/main.go
```

---

## 🔍 Console Debugging

### Frontend (Browser Console - F12)
```javascript
// Homepage
=== HOMEPAGE DEBUG ===
Mobils fetched: 47

// Detail Page
=== DEBUG INFO ===
User: <user-id> or NOT LOGGED IN
Owner ID: <owner-id>
isOwner: false
Status: tersedia
isAvailable: true
Show Buy Button: true
Show Rental Button: true
```

### Backend (Terminal)
```
// When buying a car:
Received request: POST /carapp.TransaksiService/BuyMobil from 127.0.0.1:xxxxx
[INFO] Transaksi BuyMobil berhasil untuk user xxx
[INFO] Notifikasi dibuat untuk penjual xxx
[INFO] Notifikasi dibuat untuk pembeli xxx
```

---

## 📊 Database Quick Reference

### Users Table
```sql
-- See all users
SELECT id, email, name, role FROM users;

-- Dealer user
Email: dealer@carapp.com
Password: dealer123
```

### Mobils Table
```sql
-- Count available cars
SELECT COUNT(*) FROM mobils WHERE status='tersedia';

-- See dealer's cars
SELECT merk, model, tahun, harga_jual 
FROM mobils 
WHERE owner_id = (SELECT id FROM users WHERE email='dealer@carapp.com');
```

### Transaksis Table
```sql
-- See all transactions
SELECT 
  t.tipe_transaksi,
  m.merk || ' ' || m.model as mobil,
  t.total_harga,
  t.status
FROM transaksis t
JOIN mobils m ON t.mobil_id = m.id;
```

---

## 🚨 Common Issues & Solutions

| Problem | Solution |
|---------|----------|
| **Port 3000 in use** | `Get-NetTCPConnection -LocalPort 3000 \| ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }` |
| **Port 9090 in use** | `Get-NetTCPConnection -LocalPort 9090 \| ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }` |
| **No cars showing** | Run seeder: `go run ./cmd/seeder/main.go` |
| **Buy button not showing** | Make sure you're logged in & not the car owner |
| **CORS error** | Check backend is running on port 9090 |
| **Unauthorized error** | Logout and login again (token expired) |
| **Database connection error** | Check PostgreSQL is running & credentials in `.env` |

---

## 📱 User Accounts for Testing

### Dealer Account (Pre-created by Seeder)
- **Email:** dealer@carapp.com
- **Password:** dealer123
- **Role:** admin
- **Owns:** 47 cars from Marketcheck

### Create Test Buyer Accounts
```
Buyer 1:
- Email: buyer1@test.com
- Password: test123

Buyer 2:
- Email: buyer2@test.com
- Password: test123
```

---

## 🎯 Feature Checklist

### ✅ Completed Features
- [x] User Registration & Login
- [x] JWT Authentication
- [x] Browse Cars (Public)
- [x] Car Details Page
- [x] Buy Car
- [x] Rental Car (with date selection)
- [x] Sell Your Own Car
- [x] Dashboard (My Cars, Transactions, Notifications)
- [x] Real-time Notifications (on transaction)
- [x] Inventory Seeder (47 real cars from Marketcheck)
- [x] Transaction History
- [x] gRPC-Web Backend
- [x] Next.js Frontend with TypeScript
- [x] PostgreSQL Database
- [x] AuthContext (auto token injection)

### 🚀 Optional Enhancements (Not Required)
- [ ] Pagination (for 100+ cars)
- [ ] Search & Filter
- [ ] Car Photos (from Marketcheck API)
- [ ] USD → IDR conversion
- [ ] User Profile Page
- [ ] Reviews & Ratings
- [ ] Favorite/Wishlist
- [ ] Email Notifications
- [ ] Admin Panel

---

## 📚 Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Next.js Frontend (Port 3000)             │
│  - React Components                                          │
│  - gRPC-Web Client                                           │
│  - AuthContext (JWT Token Management)                        │
│  - TypeScript                                                │
└────────────────────────┬────────────────────────────────────┘
                         │ gRPC-Web (HTTP)
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                Go gRPC-Web Server (Port 9090)               │
│  - Auth Service (Register, Login)                            │
│  - Mobil Service (List, Get, Create)                         │
│  - Transaksi Service (Buy, Rent)                             │
│  - Notifikasi Service (Create, List)                         │
│  - Dashboard Service (Summary)                               │
│  - JWT Middleware                                            │
│  - CORS Handler                                              │
└────────────────────────┬────────────────────────────────────┘
                         │ SQL
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL Database (Port 5432)                │
│  - users                                                     │
│  - mobils (47 cars from Marketcheck)                         │
│  - transaksis                                                │
│  - notifikasis                                               │
└─────────────────────────────────────────────────────────────┘

External APIs:
  ┌─────────────────────┐
  │  Marketcheck API    │ → Inventory Seeder
  └─────────────────────┘
```

---

## 🎓 Learning Resources

### Proto Files
- `proto/carapp.proto` - gRPC service definitions

### Key Backend Files
- `main.go` - Server entry point
- `internal/auth/` - Authentication logic
- `internal/mobil/` - Car business logic
- `internal/transaksi/` - Transaction logic
- `cmd/seeder/main.go` - Inventory seeder

### Key Frontend Files
- `app/page.tsx` - Homepage (car list)
- `app/mobil/[id]/page.tsx` - Car detail page
- `app/dashboard/page.tsx` - User dashboard
- `src/lib/grpcClient.ts` - gRPC client setup
- `src/context/AuthContext.tsx` - Auth state management

---

## 💡 Tips & Tricks

### Development Workflow
```powershell
# Always run in this order:
1. Start PostgreSQL
2. Run seeder (first time only)
3. Start Go backend
4. Start Next.js frontend
5. Open browser

# To reset everything:
1. Stop all servers
2. Reset database
3. Run seeder
4. Restart servers
```

### Debugging
```powershell
# Backend logs
go run main.go | Tee-Object -FilePath logs-backend.txt

# Frontend logs  
npm run dev | Tee-Object -FilePath logs-frontend.txt

# Database queries
$env:PGPASSWORD='123456'
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d car_db
```

### Performance
```sql
-- Add indexes for better performance (already done in migrations)
CREATE INDEX idx_mobils_status ON mobils(status);
CREATE INDEX idx_mobils_owner ON mobils(owner_id);
CREATE INDEX idx_transaksis_buyer ON transaksis(buyer_id);
```

---

**🎉 You're all set! Happy coding!** 🚀

For troubleshooting, see: `TROUBLESHOOTING.md`
For Langkah 14 details, see: `LANGKAH_14_SELESAI.md`
