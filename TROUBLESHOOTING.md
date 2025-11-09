# 🔧 Panduan Troubleshooting - Car Dealer App

## Masalah: Daftar Mobil Tidak Muncul & Tombol Beli Tidak Ada

### ✅ Checklist Langkah-langkah:

#### 1. **Pastikan Backend Go Server Berjalan**
```powershell
# Di terminal, jalankan:
cd c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc
go run main.go
```

**Output yang benar:**
```
2025/11/09 14:43:25 Berhasil terhubung ke database PostgreSQL!
2025/11/09 14:43:25 Server gRPC-Web (HTTP) berjalan di 0.0.0.0:9090...
```

**Jika error "port already in use":**
```powershell
# Kill process di port 9090
Get-NetTCPConnection -LocalPort 9090 -ErrorAction SilentlyContinue | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }

# Tunggu 2 detik, lalu jalankan lagi
go run main.go
```

#### 2. **Pastikan Frontend Next.js Berjalan**
```powershell
cd car_dealer_frontend
npm run dev
```

**Output yang benar:**
```
- ready started server on 0.0.0.0:3000, url: http://localhost:3000
- info using webpack 5.x
```

#### 3. **Cek Database - Apakah Ada Mobil?**

Buka terminal baru dan cek database:

```powershell
# Ganti dengan path psql Anda jika berbeda
$env:PGPASSWORD='123456'
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d car_db -c "SELECT id, merk, model, tahun, status, owner_id FROM mobils;"
```

**Jika TIDAK ADA MOBIL**, Anda perlu menambahkan mobil dulu!

#### 4. **Menambahkan Mobil untuk Testing**

**Langkah-langkah:**

1. **Login sebagai User A:**
   - Buka browser: http://localhost:3000/login
   - Login atau register user pertama (misal: `usera@test.com`)

2. **Jual Mobil:**
   - Setelah login, klik **"Jual Mobil"** di Navbar
   - Atau langsung ke: http://localhost:3000/mobil/jual
   - Isi form:
     - **Merk**: Toyota (atau merk lain dari NHTSA)
     - **Model**: Camry
     - **Tahun**: 2023
     - **Kondisi**: baru / bekas
     - **Lokasi**: Jakarta
     - **Harga Jual**: 500000000 (500 juta)
     - **Harga Rental/Hari**: 1000000 (1 juta)
     - **Deskripsi**: Mobil bagus siap pakai
   - Klik **"Jual Mobil"**

3. **Verifikasi di Database:**
```powershell
$env:PGPASSWORD='123456'
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d car_db -c "SELECT id, merk, model, tahun, status, owner_id FROM mobils WHERE status='tersedia';"
```

4. **Logout User A:**
   - Klik "Logout" di Navbar

5. **Login sebagai User B:**
   - Klik "Register" dan buat akun baru (misal: `userb@test.com`)
   - Atau login dengan akun lain yang sudah ada

6. **Cek Homepage:**
   - Buka: http://localhost:3000
   - Seharusnya mobil yang dijual User A muncul!

7. **Cek Detail Mobil:**
   - Klik salah satu mobil
   - Seharusnya ada **Debug Panel kuning** yang menunjukkan:
     - User ID (User B)
     - Owner ID (User A)
     - Is Owner: NO
     - Status: tersedia
     - Is Available: YES
     - Show Buy: YES
     - Show Rental: YES

8. **Test Beli Mobil:**
   - Klik tombol **"Beli Mobil Ini"**
   - Seharusnya muncul pesan: "Pembelian berhasil!"
   - Setelah 3 detik, otomatis redirect ke Dashboard

#### 5. **Debugging dengan Browser Console**

Buka Developer Tools di browser (F12), lalu cek tab **Console**.

**Di Homepage (`/`):**
```
=== HOMEPAGE DEBUG ===
Mobils fetched: 1
Mobils: [{id: "xxx", merk: "Toyota", model: "Camry", status: "tersedia", ownerId: "xxx"}]
=====================
```

**Di Detail Mobil (`/mobil/[id]`):**
```
=== DEBUG INFO ===
User: xxx-user-b-id
Owner ID: xxx-user-a-id
isOwner: false
Status: tersedia
isAvailable: true
Harga Jual: 500000000
Harga Rental: 1000000
Show Buy Button: true
Show Rental Button: true
================
```

#### 6. **Masalah Umum & Solusi**

| Masalah | Penyebab | Solusi |
|---------|----------|--------|
| **Daftar mobil kosong** | Tidak ada mobil dengan status 'tersedia' di DB | Tambahkan mobil via `/mobil/jual` |
| **Tombol beli tidak muncul** | User belum login | Login dulu di `/login` |
| **Tombol beli tidak muncul** | User adalah pemilik mobil | Logout dan login dengan user lain |
| **Tombol beli tidak muncul** | Status mobil bukan 'tersedia' | Cek status di database |
| **"Gagal memuat mobil"** | Backend Go tidak berjalan | Restart `go run main.go` |
| **CORS error** | Port tidak match | Pastikan Next.js di port 3000 |
| **"Unauthorized"** | Token expired / tidak valid | Logout dan login lagi |

#### 7. **Cek Logs Backend**

Saat User B membeli mobil, di terminal Go server seharusnya muncul:

```
2025/11/09 14:45:00 Received request: POST /carapp.TransaksiService/BuyMobil from 127.0.0.1:xxxx
2025/11/09 14:45:00 [INFO] Transaksi BuyMobil berhasil untuk user xxx, mobil xxx
2025/11/09 14:45:00 [INFO] Notifikasi dibuat untuk penjual xxx
2025/11/09 14:45:00 [INFO] Notifikasi dibuat untuk pembeli xxx
```

#### 8. **Reset Database (Jika Diperlukan)**

Jika database berantakan, reset dengan:

```powershell
$env:PGPASSWORD='123456'
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d car_db -c "TRUNCATE TABLE mobils, transaksis, notifikasis RESTART IDENTITY CASCADE;"
```

**⚠️ WARNING:** Ini akan menghapus SEMUA data mobil, transaksi, dan notifikasi!

### 🎯 Skenario Testing Lengkap

```
[User A] --> Register/Login
         --> Jual Mobil (via /mobil/jual)
         --> Logout

[User B] --> Register/Login
         --> Lihat Homepage (/) --> Mobil User A muncul ✅
         --> Klik Detail Mobil --> Debug panel muncul ✅
         --> Tombol "Beli Mobil Ini" muncul ✅
         --> Klik "Beli" --> "Pembelian berhasil!" ✅
         --> Redirect ke Dashboard --> Transaksi muncul ✅

[User A] --> Login kembali
         --> Dashboard --> Ada notifikasi "Mobil Anda terjual!" ✅
```

### 📞 Jika Masih Bermasalah

1. **Screenshot error di browser console**
2. **Copy paste output terminal Go server**
3. **Cek apakah kedua server benar-benar berjalan:**
   ```powershell
   Get-NetTCPConnection -LocalPort 3000,9090 -ErrorAction SilentlyContinue
   ```
   Seharusnya menunjukkan kedua port LISTEN.

---

## 🚀 Quick Start (Dari Awal)

```powershell
# Terminal 1 - Backend
cd c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc
go run main.go

# Terminal 2 - Frontend
cd c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc\car_dealer_frontend
npm run dev

# Browser
# 1. http://localhost:3000/register --> Buat User A
# 2. http://localhost:3000/mobil/jual --> Jual mobil
# 3. Logout
# 4. http://localhost:3000/register --> Buat User B
# 5. http://localhost:3000 --> Lihat mobil User A
# 6. Klik detail --> Lihat debug panel
# 7. Klik "Beli Mobil Ini" --> Test transaksi!
```

---

**Debug panel kuning akan muncul di halaman detail mobil untuk membantu troubleshooting!**
