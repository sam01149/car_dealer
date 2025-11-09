# 🚗 Car Dealer gRPC App - Testing Guide

## 📋 Ringkasan Masalah

**Masalah yang dialami:**
- Daftar mobil tidak muncul di homepage
- Tombol "Beli Mobil" tidak muncul di halaman detail

**Kemungkinan penyebab:**
1. ❌ Tidak ada mobil dengan status 'tersedia' di database
2. ❌ User belum login
3. ❌ User yang login adalah pemilik mobil itu sendiri
4. ❌ Backend server tidak berjalan

---

## 🛠️ Tools yang Sudah Dibuat

### 1. **start-app.ps1** - Server Management
```powershell
# Check status
.\start-app.ps1 -Action status

# Start backend (Terminal 1)
.\start-app.ps1 -Action backend

# Start frontend (Terminal 2)
.\start-app.ps1 -Action frontend

# Stop all servers
.\start-app.ps1 -Action stop
```

### 2. **test-db.ps1** - Database Management
```powershell
# Check database
.\test-db.ps1 -Action check

# Show available mobils
.\test-db.ps1 -Action available

# Show all users
.\test-db.ps1 -Action users

# Add test mobils (requires existing users)
.\test-db.ps1 -Action seed
```

### 3. **TROUBLESHOOTING.md** - Panduan Lengkap
Buka file ini untuk panduan lengkap troubleshooting.

---

## 🚀 Quick Start Testing

### Step 1: Cek Status Server
```powershell
.\start-app.ps1 -Action status
```

**Expected output:**
```
✅ Backend (Go gRPC):     Running on port 9090
✅ Frontend (Next.js):    Running on port 3000
```

**Jika salah satu tidak running:**
```powershell
# Terminal 1
.\start-app.ps1 -Action backend

# Terminal 2 (buka terminal baru)
.\start-app.ps1 -Action frontend
```

### Step 2: Cek Database
```powershell
.\test-db.ps1 -Action check
```

**Jika tidak ada mobil:**
1. Buka browser: http://localhost:3000/register
2. Register User A (misal: `usera@test.com`)
3. Setelah register, otomatis redirect ke dashboard
4. Klik "Jual Mobil" di Navbar
5. Isi form dan submit

**ATAU gunakan seed data:**
```powershell
# Pastikan sudah ada minimal 1 user di database
.\test-db.ps1 -Action seed
```

### Step 3: Test Skenario Lengkap

#### A. Sebagai **User A** (Penjual)

1. **Register/Login**
   ```
   URL: http://localhost:3000/register
   Email: usera@test.com
   Password: password123
   ```

2. **Jual Mobil**
   ```
   URL: http://localhost:3000/mobil/jual
   
   Merk: Toyota
   Model: Camry
   Tahun: 2023
   Kondisi: baru
   Lokasi: Jakarta
   Harga Jual: 500000000
   Harga Rental: 1000000
   Deskripsi: Mobil sedan mewah
   ```

3. **Logout**
   - Klik "Logout" di Navbar

#### B. Sebagai **User B** (Pembeli)

1. **Register/Login**
   ```
   URL: http://localhost:3000/register
   Email: userb@test.com
   Password: password123
   ```

2. **Lihat Homepage**
   ```
   URL: http://localhost:3000
   ```
   
   **✅ Expected:** Mobil yang dijual User A muncul di list

3. **Lihat Detail Mobil**
   - Klik salah satu mobil
   
   **✅ Expected:** 
   - Muncul **debug panel kuning** di atas dengan info:
     ```
     User ID: xxx-userb-xxx
     Owner ID: xxx-usera-xxx
     Is Owner: NO
     Status: tersedia
     Is Available: YES
     Show Buy: YES
     Show Rental: YES
     ```
   - Tombol **"Beli Mobil Ini"** muncul
   - Tombol **"Rental Mobil Ini"** muncul

4. **Beli Mobil**
   - Klik tombol "Beli Mobil Ini"
   
   **✅ Expected:**
   - Muncul pesan: "Memproses transaksi Anda..."
   - Lalu: "Pembelian berhasil! (ID Transaksi: xxx) Anda akan diarahkan..."
   - Setelah 3 detik, redirect ke Dashboard
   - Di Dashboard muncul transaksi baru

5. **Test Rental** (Opsional)
   - Ulangi langkah 2-3 dengan mobil lain
   - Klik "Rental Mobil Ini"
   - Pilih tanggal mulai dan selesai
   - Klik "Konfirmasi Rental"
   
   **✅ Expected:** Sama seperti beli, tapi tipe transaksi "rental"

#### C. Kembali ke **User A** (Penjual)

1. **Logout User B**
2. **Login User A**
3. **Cek Dashboard**
   
   **✅ Expected:**
   - Ada notifikasi: "Mobil Anda [Toyota Camry] telah terjual!"
   - Status mobil berubah menjadi "terjual"

---

## 🐛 Debugging

### Debug Panel di Detail Mobil

File `app/mobil/[id]/page.tsx` sudah ditambahkan **debug panel kuning** yang menampilkan:
- User ID saat ini
- Owner ID mobil
- Apakah user adalah owner
- Status mobil
- Apakah mobil available
- Harga jual dan rental
- Apakah tombol Buy/Rental seharusnya muncul

### Browser Console Logs

Buka **Developer Tools (F12)** → Tab **Console**

**Di Homepage:**
```javascript
=== HOMEPAGE DEBUG ===
Mobils fetched: 1
Mobils: [{...}]
=====================
```

**Di Detail Mobil:**
```javascript
=== DEBUG INFO ===
User: xxx
Owner ID: xxx
isOwner: false
Status: tersedia
...
================
```

### Backend Logs

Di terminal Go server, perhatikan log:
```
2025/11/09 14:45:00 Received request: POST /carapp.MobilService/ListMobil
2025/11/09 14:45:00 Received request: POST /carapp.TransaksiService/BuyMobil
```

---

## 📊 Verifikasi Database

```powershell
# Cek mobil tersedia
.\test-db.ps1 -Action available

# Cek transaksi
.\test-db.ps1 -Action transactions

# Cek semua
.\test-db.ps1 -Action check
```

---

## 🔄 Reset & Mulai Ulang

Jika ingin mulai dari awal:

```powershell
# 1. Reset database (HATI-HATI!)
.\test-db.ps1 -Action reset

# 2. Restart servers
.\start-app.ps1 -Action stop
.\start-app.ps1 -Action backend   # Terminal 1
.\start-app.ps1 -Action frontend  # Terminal 2

# 3. Mulai testing dari awal
```

---

## ✅ Checklist Testing

- [ ] Backend server running (port 9090)
- [ ] Frontend server running (port 3000)
- [ ] Database memiliki minimal 2 users
- [ ] Database memiliki minimal 1 mobil dengan status 'tersedia'
- [ ] User A bisa jual mobil
- [ ] User B bisa lihat mobil di homepage
- [ ] User B bisa lihat detail mobil
- [ ] Debug panel kuning muncul di detail mobil
- [ ] Tombol "Beli Mobil" muncul untuk User B
- [ ] User B bisa membeli mobil
- [ ] Transaksi sukses dan redirect ke dashboard
- [ ] User A menerima notifikasi di dashboard
- [ ] Status mobil berubah menjadi "terjual"

---

## 📞 Jika Masih Bermasalah

1. **Screenshot debug panel kuning**
2. **Screenshot browser console (F12 → Console)**
3. **Copy log dari terminal Go server**
4. **Jalankan:**
   ```powershell
   .\start-app.ps1 -Action status
   .\test-db.ps1 -Action check
   ```

Share output tersebut untuk troubleshooting lebih lanjut.

---

## 📚 File Penting

- `TROUBLESHOOTING.md` - Panduan troubleshooting lengkap
- `start-app.ps1` - Script untuk manage server
- `test-db.ps1` - Script untuk manage database
- `app/mobil/[id]/page.tsx` - Halaman detail mobil (dengan debug panel)
- `app/page.tsx` - Homepage (dengan console.log debug)

---

**Langkah 13 sudah selesai! Selamat testing! 🎉**
