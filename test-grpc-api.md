# 🧪 Testing gRPC API - CarApp

## Prerequisites
- Backend running di `localhost:9090`
- Install grpcurl: `choco install grpcurl` atau download dari https://github.com/fullstorydev/grpcurl/releases

## 🔍 List All Services & Methods

```bash
# List semua services
grpcurl -plaintext localhost:9090 list

# List methods di AuthService
grpcurl -plaintext localhost:9090 list carapp.AuthService

# List methods di MobilService  
grpcurl -plaintext localhost:9090 list carapp.MobilService

# Describe service detail
grpcurl -plaintext localhost:9090 describe carapp.AuthService
```

---

## 1️⃣ AuthService - Login & Register

### Register New User
```bash
grpcurl -plaintext -d '{
  "name": "Test User",
  "email": "testuser@example.com",
  "password": "password123",
  "phone": "+6281234567890"
}' localhost:9090 carapp.AuthService/Register
```

**Expected Response:**
```json
{
  "user": {
    "id": "uuid-123...",
    "name": "Test User",
    "email": "testuser@example.com",
    "phone": "+6281234567890",
    "role": "client"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login User
```bash
grpcurl -plaintext -d '{
  "email": "testuser@example.com",
  "password": "password123"
}' localhost:9090 carapp.AuthService/Login
```

### Login as Dealer (Admin)
```bash
grpcurl -plaintext -d '{
  "email": "dealer@carapp.com",
  "password": "dealer123"
}' localhost:9090 carapp.AuthService/Login
```

**Simpan token dari response untuk request berikutnya!**

---

## 2️⃣ MobilService - CRUD Mobil

### List Mobil (Public - No Auth)
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "limit": 10,
  "filter_status": "tersedia"
}' localhost:9090 carapp.MobilService/ListMobil
```

### Get Detail Mobil (Public - No Auth)
```bash
# Ganti dengan ID mobil yang ada
grpcurl -plaintext -d '{
  "mobil_id": "uuid-mobil-123"
}' localhost:9090 carapp.MobilService/GetMobil
```

### Create Mobil (Auth Required)
```bash
# Ganti YOUR_JWT_TOKEN dengan token dari login
grpcurl -plaintext \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "merk": "Toyota",
    "model": "Avanza",
    "tahun": 2023,
    "kondisi": "baru",
    "deskripsi": "Mobil keluarga nyaman, siap pakai",
    "harga_jual": 250000000,
    "foto_url": "https://example.com/avanza.jpg",
    "lokasi": "Jakarta Selatan",
    "harga_rental_per_hari": 500000
  }' localhost:9090 carapp.MobilService/CreateMobil
```

---

## 3️⃣ NhtsaDataService - Data Merek & Model

### Get All Makes (Brands)
```bash
grpcurl -plaintext -d '{}' localhost:9090 carapp.NhtsaDataService/GetMakes
```

**Response akan berisi list merek mobil:**
```json
{
  "makes": [
    {"brandId": "440", "name": "Audi"},
    {"brandId": "441", "name": "BMW"},
    {"brandId": "445", "name": "Ferrari"},
    ...
  ]
}
```

### Get Models for Make (Brand)
```bash
# Contoh: Get models untuk Ferrari (brand_id: 445)
grpcurl -plaintext -d '{
  "brand_id": "445"
}' localhost:9090 carapp.NhtsaDataService/GetModelsForMake
```

**Response:**
```json
{
  "models": [
    {"modelId": "1234", "brandId": "445", "name": "488 GTB"},
    {"modelId": "1235", "brandId": "445", "name": "F8 Tributo"},
    {"modelId": "1236", "brandId": "445", "name": "SF90 Stradale"}
  ]
}
```

---

## 4️⃣ TransaksiService - Beli Mobil

### Buy Mobil (Auth Required)
```bash
grpcurl -plaintext \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "mobil_id": "uuid-mobil-yang-mau-dibeli"
  }' localhost:9090 carapp.TransaksiService/BuyMobil
```

**Expected Response:**
```json
{
  "id": "transaksi-uuid-123",
  "mobilId": "mobil-uuid",
  "penjualId": "seller-uuid",
  "pembeliId": "buyer-uuid",
  "total": 250000000,
  "status": "selesai"
}
```

---

## 5️⃣ DashboardService - User Dashboard

### Get Dashboard Summary (Auth Required)
```bash
grpcurl -plaintext \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{}' localhost:9090 carapp.DashboardService/GetDashboard
```

**Response:**
```json
{
  "totalMobilAnda": 5,
  "transaksiAktif": 2,
  "pendapatanTerakhir": 500000000,
  "notifikasiBaru": 3
}
```

---

## 6️⃣ NotifikasiService - Real-time Notifications (Stream)

### Get Notifications Stream (Auth Required)
```bash
grpcurl -plaintext \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{}' localhost:9090 carapp.NotifikasiService/GetNotifications
```

**Response (Streaming):**
```json
{
  "id": "notif-1",
  "userId": "user-uuid",
  "tipe": "jual",
  "pesan": "Mobil Anda berhasil terjual!",
  "priority": "high",
  "createdAt": "2025-11-11T10:00:00Z"
}
{
  "id": "notif-2",
  "userId": "user-uuid",
  "tipe": "beli",
  "pesan": "Pembelian berhasil!",
  "priority": "medium",
  "createdAt": "2025-11-11T10:05:00Z"
}
...
```

**Note:** Streaming akan tetap terbuka dan mengirim notifikasi baru secara real-time!

---

## 🎯 Test Scenarios untuk Presentasi

### Scenario 1: Register → Login → Create Mobil
```bash
# 1. Register
grpcurl -plaintext -d '{"name":"Demo User","email":"demo@test.com","password":"demo123","phone":"+6281234567890"}' localhost:9090 carapp.AuthService/Register

# 2. Login (simpan token)
grpcurl -plaintext -d '{"email":"demo@test.com","password":"demo123"}' localhost:9090 carapp.AuthService/Login

# 3. Create Mobil (ganti TOKEN)
grpcurl -plaintext -H "Authorization: Bearer TOKEN" -d '{"merk":"Honda","model":"Civic","tahun":2023,"kondisi":"baru","deskripsi":"Sport sedan","harga_jual":450000000,"foto_url":"https://example.com/civic.jpg","lokasi":"Jakarta"}' localhost:9090 carapp.MobilService/CreateMobil
```

### Scenario 2: List Mobil → Get Detail → Buy
```bash
# 1. List mobil tersedia
grpcurl -plaintext -d '{"page":1,"limit":5}' localhost:9090 carapp.MobilService/ListMobil

# 2. Get detail (copy ID dari list)
grpcurl -plaintext -d '{"mobil_id":"COPY_ID_HERE"}' localhost:9090 carapp.MobilService/GetMobil

# 3. Login sebagai buyer
grpcurl -plaintext -d '{"email":"buyer@test.com","password":"buyer123"}' localhost:9090 carapp.AuthService/Login

# 4. Buy mobil (ganti TOKEN dan mobil_id)
grpcurl -plaintext -H "Authorization: Bearer TOKEN" -d '{"mobil_id":"MOBIL_ID"}' localhost:9090 carapp.TransaksiService/BuyMobil
```

### Scenario 3: NHTSA Data Integration
```bash
# 1. Get all brands
grpcurl -plaintext -d '{}' localhost:9090 carapp.NhtsaDataService/GetMakes

# 2. Get models for Rolls-Royce (brand_id: 493)
grpcurl -plaintext -d '{"brand_id":"493"}' localhost:9090 carapp.NhtsaDataService/GetModelsForMake
```

---

## 📊 Monitoring & Debugging

### Check Server Reflection (List available services)
```bash
grpcurl -plaintext localhost:9090 list
```

### Describe Message Structure
```bash
grpcurl -plaintext localhost:9090 describe carapp.Mobil
grpcurl -plaintext localhost:9090 describe carapp.CreateMobilRequest
```

### Test with Verbose Output
```bash
grpcurl -plaintext -v -d '{}' localhost:9090 carapp.DashboardService/GetDashboard
```

---

## 🚨 Common Errors & Solutions

### Error: "connection refused"
- **Solution:** Pastikan backend running di `localhost:9090`
- Check: `netstat -ano | findstr :9090`

### Error: "Unauthenticated"
- **Solution:** Tambahkan header Authorization dengan token JWT
- Format: `-H "Authorization: Bearer YOUR_TOKEN"`

### Error: "NotFound" atau "PermissionDenied"
- **Solution:** Check apakah user punya akses ke resource (e.g., owner mobil)

### Error: "InvalidArgument"
- **Solution:** Periksa format data request (required fields, types)

---

## 💡 Tips untuk Presentasi

1. **Demo Live:** Jalankan test scenarios di terminal sambil explain
2. **Bandingkan dengan REST:** Show ukuran data (JSON vs Protobuf)
3. **Streaming Demo:** Show real-time notifications
4. **Performance:** Explain kenapa gRPC lebih cepat (binary, HTTP/2)
5. **Cache NHTSA:** Explain bagaimana cache mengurangi API calls

---

## 📦 Alternative: Using Postman

Postman juga support gRPC! Import proto file:
1. Open Postman
2. New → gRPC Request
3. Import proto file: `proto/carapp.proto`
4. Server URL: `localhost:9090`
5. Test methods dengan UI yang lebih user-friendly

---

## 🎓 Konsep Penting untuk Dijelaskan

### 1. Protocol Buffers
- Binary format, lebih efisien dari JSON
- Strongly typed (compile-time validation)
- Backward/forward compatible

### 2. HTTP/2
- gRPC pakai HTTP/2 (multiplexing, streaming)
- REST API pakai HTTP/1.1 (request-response only)

### 3. Streaming
- **Unary:** 1 request → 1 response (seperti REST)
- **Server Streaming:** 1 request → many responses (NotifikasiService)
- **Client Streaming:** many requests → 1 response
- **Bidirectional:** many requests ↔ many responses

### 4. Middleware/Interceptor
- Authentication (JWT validation)
- Logging
- Error handling
- CORS untuk gRPC-Web

### 5. Code Generation
- Proto file → Generate code otomatis (Go, TypeScript, Python, dll)
- Type-safe, no manual parsing
