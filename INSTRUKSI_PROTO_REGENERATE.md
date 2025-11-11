# INSTRUKSI REGENERATE PROTO FILES

## ⚠️ PENTING: Setelah semua perubahan selesai, jalankan langkah berikut:

### 1. Install protoc (jika belum ada)
Download dari: https://github.com/protocolbuffers/protobuf/releases
Atau install via package manager:
- Windows: `choco install protoc`
- Mac: `brew install protobuf`

### 2. Install protoc plugins untuk Go
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Regenerate Go proto files
```bash
cd c:\Users\sam\Documents\File_Coding\golang\car_dealer_grpc
protoc --go_out=. --go-grpc_out=. proto/carapp.proto
```

### 4. Regenerate TypeScript/JavaScript proto files untuk frontend
```bash
cd car_dealer_frontend
npm install -g protoc-gen-grpc-web
protoc -I=../proto ../proto/carapp.proto \
  --js_out=import_style=commonjs:src/proto \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:src/proto
```

## Perubahan Proto yang Dilakukan:

1. ✅ **Hapus `harga_rental_per_hari`** dari message `Mobil`
2. ✅ **Tambah `foto_url`** ke message `Mobil` (field 9)
3. ✅ **Tambah `owner_name`** ke message `Mobil` (field 13)
4. ✅ **Tambah `foto_url`** ke message `CreateMobilRequest` (field 7)
5. ✅ **Hapus `RentMobil` RPC** dari `TransaksiService`
6. ✅ **Hapus `CompleteRental` RPC** dari `TransaksiService`
7. ✅ **Hapus message** `RentMobilRequest`, `CompleteRentalRequest`, `TransaksiRentalResponse`

## Setelah Proto Files Di-generate:

1. Update mobil_service.go:
   - Uncomment `mobil.OwnerName = ownerName` di GetMobil
   - Uncomment `mobil.FotoUrl = fotoUrl.String` di GetMobil dan ListMobil
   - Update CreateMobil untuk menggunakan `req.FotoUrl`

2. Update frontend:
   - Re-import proto files yang baru
   - Gunakan `mobil.getOwnerName()` untuk tampilkan nama penjual
   - Gunakan `mobil.getFotoUrl()` untuk tampilkan foto
   - Implement upload foto dengan kompresi

## Database Migration Sudah Berjalan: ✅

Migration 002 telah menghapus:
- Tabel `transaksi_rental`
- Kolom `harga_rental_per_hari` dari tabel `mobils`

Dan menambahkan:
- Kolom `foto_url` ke tabel `mobils`

## Next Steps (Manual):

1. **Install protoc** di sistem Anda
2. **Regenerate proto files** dengan perintah di atas
3. **Uncomment kode** di mobil_service.go yang sudah disiapkan
4. **Test aplikasi** untuk memastikan semua berjalan
5. **Implement upload foto** dengan kompresi di frontend
