# 💱 Perbaikan Konversi Harga USD → IDR

## ❌ Masalah Sebelumnya

Harga dari Marketcheck API dalam **USD (Dollar Amerika)**, tapi disimpan langsung ke database tanpa konversi:

```
Contoh: Ford Ranger 2023
Harga API: $31,499 USD
Disimpan: Rp 31.499 ❌ (SALAH!)
Seharusnya: Rp 497.684.200 ✅
```

## ✅ Solusi

Menambahkan konversi otomatis di seeder dengan kurs tetap:

```go
// KONVERSI USD KE IDR
const USD_TO_IDR = 15800.0  // 1 USD = Rp 15.800
hargaJualIDR := mobil.Price * USD_TO_IDR
```

### Perubahan di `cmd/seeder/main.go`:

1. **Tambah konversi USD → IDR**
   ```go
   const USD_TO_IDR = 15800.0
   hargaJualIDR := mobil.Price * USD_TO_IDR
   hargaRentalIDR := hargaJualIDR * 0.005  // 0.5% per hari
   ```

2. **Tambah auto-delete mobil lama**
   ```go
   // Reset mobil lama sebelum insert baru
   dbConn.Exec("DELETE FROM mobils WHERE owner_id = $1", dealerUserID)
   ```

## 📊 Contoh Hasil Konversi

| Mobil | Harga USD | Harga IDR (15.800x) |
|-------|-----------|---------------------|
| 2025 Ford Ranger Lariat | $52,660 | Rp 832.028.000 |
| 2021 Ford Ranger Lariat | $31,845 | Rp 503.151.000 |
| 2024 Ford Ranger XLT | $38,461 | Rp 607.684.200 |
| 2021 Ford Ranger XL | $28,995 | Rp 458.121.000 |
| 2023 Ford Ranger Lariat | $34,995 | Rp 552.921.000 |

## 🔄 Cara Mengubah Kurs

Jika ingin mengubah kurs (misalnya 1 USD = Rp 16.000):

1. Edit file: `cmd/seeder/main.go`
2. Cari baris:
   ```go
   const USD_TO_IDR = 15800.0
   ```
3. Ubah menjadi:
   ```go
   const USD_TO_IDR = 16000.0  // Kurs baru
   ```
4. Jalankan ulang seeder:
   ```powershell
   go run ./cmd/seeder/main.go
   ```

## 🚀 Setelah Update

Sekarang ketika Anda buka:
- **Homepage** → Harga tampil: **Rp 497.684.200** ✅
- **Detail Mobil** → Harga tampil: **Rp 497.684.200** ✅
- **Rental** → Harga per hari: **Rp 2.488.421** ✅ (0.5% dari harga jual)

## 💡 Catatan

- Kurs **Rp 15.800** adalah kurs approximate (tidak realtime)
- Jika ingin realtime, bisa integrasikan dengan Exchange Rate API
- Untuk production, pertimbangkan:
  - Simpan harga USD dan IDR terpisah
  - Update kurs otomatis dari API
  - Tampilkan kedua mata uang (USD & IDR)

## 🧪 Verifikasi

Refresh browser Anda dan cek:
```
✅ Homepage: Harga dalam jutaan/ratusan juta Rupiah
✅ Detail Page: Format "Rp 497.684.200"
✅ Rental: Format "Rp 2.488.421 / hari"
```

---

**Update selesai! Harga sekarang sudah benar dalam Rupiah!** 🎉
