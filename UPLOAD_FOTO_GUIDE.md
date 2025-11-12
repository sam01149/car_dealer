# 📸 Upload Foto Mobil - Implementasi Guide

## 🎯 Konsep

**gRPC Client Streaming Upload:**
1. Client buka stream ke server
2. Kirim chunk pertama dengan **metadata** (filename, content-type, size)
3. Kirim chunks berikutnya dengan **binary data** (64KB per chunk)
4. Server terima semua chunks → gabungkan → simpan ke `/uploads/`
5. Server return URL foto yang sudah tersimpan

## ✅ Langkah Implementasi

### 1. Regenerate Proto Files

```powershell
# Pastikan protoc sudah terinstall
protoc --version

# Regenerate proto (Go)
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       proto/carapp.proto

# Regenerate proto (TypeScript untuk frontend)
cd car_dealer_frontend
npx protoc --plugin=protoc-gen-ts_proto=./node_modules/.bin/protoc-gen-ts_proto `
           --ts_proto_out=./src/proto `
           --ts_proto_opt=outputServices=grpc-js,env=browser,useOptionals=messages `
           ../proto/carapp.proto
```

### 2. Install Dependency (jika belum)

```powershell
# Backend (Go)
go get github.com/google/uuid

# Frontend (Next.js) - sudah ada grpc-web
npm install
```

### 3. Buat Folder Uploads

```powershell
mkdir uploads
```

### 4. Test Upload (Manual via Go Client)

Buat file `test_upload.go`:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "carapp.com/m/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const chunkSize = 64 * 1024 // 64KB per chunk

func main() {
	// Connect ke server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMobilServiceClient(conn)

	// Path foto yang mau diupload
	filePath := "test_civic.jpg" // Ganti dengan path foto kamu
	
	// Buka file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Gagal buka file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, _ := file.Stat()

	// Buka stream dengan auth token (jika perlu)
	token := "YOUR_JWT_TOKEN_HERE"
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	
	stream, err := client.UploadFoto(ctx)
	if err != nil {
		log.Fatalf("Gagal buka stream: %v", err)
	}

	// Kirim metadata dulu
	err = stream.Send(&pb.UploadFotoRequest{
		Data: &pb.UploadFotoRequest_Metadata{
			Metadata: &pb.UploadFotoMetadata{
				Filename:    fileInfo.Name(),
				ContentType: "image/jpeg",
				FileSize:    fileInfo.Size(),
			},
		},
	})
	if err != nil {
		log.Fatalf("Gagal kirim metadata: %v", err)
	}
	fmt.Println("✓ Metadata terkirim")

	// Kirim chunks
	buffer := make([]byte, chunkSize)
	chunkCount := 0
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Gagal baca file: %v", err)
		}

		err = stream.Send(&pb.UploadFotoRequest{
			Data: &pb.UploadFotoRequest_Chunk{
				Chunk: buffer[:n],
			},
		})
		if err != nil {
			log.Fatalf("Gagal kirim chunk: %v", err)
		}
		chunkCount++
		fmt.Printf("✓ Chunk %d terkirim (%d bytes)\n", chunkCount, n)
	}

	// Close stream dan terima response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Gagal terima response: %v", err)
	}

	fmt.Printf("\n🎉 Upload berhasil!\n")
	fmt.Printf("   URL: %s\n", resp.Url)
	fmt.Printf("   Message: %s\n", resp.Message)
}
```

Run:
```powershell
go run test_upload.go
```

### 5. Frontend Implementation (Next.js)

Update `car_dealer_frontend/app/mobil/jual/page.tsx`:

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { getGrpcClient } from '@/lib/grpcClient';
import { UploadFotoRequest, UploadFotoMetadata } from '@/proto/carapp_pb';

const CHUNK_SIZE = 64 * 1024; // 64KB

export default function JualMobilPage() {
  const router = useRouter();
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [fotoFile, setFotoFile] = useState<File | null>(null);
  const [fotoUrl, setFotoUrl] = useState('');

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validasi
    if (!file.type.startsWith('image/')) {
      alert('File harus berupa gambar');
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      alert('Ukuran file maksimal 5MB');
      return;
    }

    setFotoFile(file);
    console.log('✓ File dipilih:', file.name, file.size, 'bytes');
  };

  const uploadFoto = async (): Promise<string> => {
    if (!fotoFile || !token) throw new Error('File atau token tidak ada');

    return new Promise((resolve, reject) => {
      const client = getGrpcClient();
      const metadata = { authorization: `Bearer ${token}` };
      
      const stream = client.uploadFoto(metadata, (err, response) => {
        if (err) {
          console.error('Upload error:', err);
          reject(err);
          return;
        }
        console.log('✓ Upload selesai:', response?.toObject());
        resolve(response!.getUrl());
      });

      // Kirim metadata dulu
      const metadataMsg = new UploadFotoMetadata();
      metadataMsg.setFilename(fotoFile.name);
      metadataMsg.setContentType(fotoFile.type);
      metadataMsg.setFileSize(fotoFile.size);

      const firstReq = new UploadFotoRequest();
      firstReq.setMetadata(metadataMsg);
      stream.write(firstReq);

      // Baca file dan kirim chunks
      const reader = new FileReader();
      let offset = 0;

      const readNextChunk = () => {
        if (offset >= fotoFile.size) {
          stream.end(); // Selesai
          return;
        }

        const chunk = fotoFile.slice(offset, offset + CHUNK_SIZE);
        reader.readAsArrayBuffer(chunk);
      };

      reader.onload = (e) => {
        const arrayBuffer = e.target?.result as ArrayBuffer;
        const uint8Array = new Uint8Array(arrayBuffer);

        const req = new UploadFotoRequest();
        req.setChunk(uint8Array);
        stream.write(req);

        offset += arrayBuffer.byteLength;
        setUploadProgress(Math.round((offset / fotoFile.size) * 100));

        readNextChunk(); // Baca chunk berikutnya
      };

      reader.onerror = () => {
        reject(new Error('Gagal membaca file'));
      };

      readNextChunk(); // Mulai baca
    });
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!token) {
      alert('Silakan login terlebih dahulu');
      router.push('/login');
      return;
    }

    setLoading(true);

    try {
      // 1. Upload foto dulu
      let uploadedUrl = '';
      if (fotoFile) {
        console.log('Uploading foto...');
        uploadedUrl = await uploadFoto();
        setFotoUrl(uploadedUrl);
        console.log('✓ Foto URL:', uploadedUrl);
      }

      // 2. Buat mobil dengan foto URL
      const formData = new FormData(e.currentTarget);
      const client = getGrpcClient();
      const metadata = { authorization: `Bearer ${token}` };

      const request = new CreateMobilRequest();
      request.setMerk(formData.get('merk') as string);
      request.setModel(formData.get('model') as string);
      request.setTahun(parseInt(formData.get('tahun') as string));
      request.setKondisi(formData.get('kondisi') as string);
      request.setDeskripsi(formData.get('deskripsi') as string);
      request.setHargaJual(parseFloat(formData.get('harga_jual') as string));
      request.setFotoUrl(uploadedUrl || ''); // Gunakan URL hasil upload
      request.setLokasi(formData.get('lokasi') as string);

      client.createMobil(request, metadata, (err, response) => {
        setLoading(false);
        if (err) {
          console.error('Error:', err);
          alert(`Gagal: ${err.message}`);
          return;
        }
        alert('Mobil berhasil dipasang!');
        router.push('/dashboard');
      });

    } catch (err: any) {
      setLoading(false);
      alert(`Error: ${err.message}`);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-6">Jual Mobil Anda</h1>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        {/* ... field lainnya ... */}

        {/* Upload Foto */}
        <div>
          <label className="block text-sm font-medium mb-1">Foto Mobil</label>
          <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className="w-full border rounded px-3 py-2"
          />
          {fotoFile && (
            <p className="text-sm text-gray-600 mt-1">
              {fotoFile.name} ({(fotoFile.size / 1024).toFixed(2)} KB)
            </p>
          )}
          {uploadProgress > 0 && uploadProgress < 100 && (
            <div className="mt-2">
              <div className="bg-gray-200 rounded-full h-2">
                <div
                  className="bg-blue-600 h-2 rounded-full transition-all"
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <p className="text-sm text-center mt-1">{uploadProgress}%</p>
            </div>
          )}
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {loading ? 'Memproses...' : 'Pasang Iklan'}
        </button>
      </form>
    </div>
  );
}
```

## 🔧 Troubleshooting

### Error: "undefined: pb.MobilService_UploadFotoServer"
**Solusi:** Proto belum di-regenerate. Jalankan:
```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/carapp.proto
```

### Error: "folder uploads tidak ada"
**Solusi:** 
```powershell
mkdir uploads
```

### Frontend tidak bisa upload
**Solusi:** Pastikan grpc-web proxy (Envoy) sudah support streaming. Update `envoy.yaml` jika perlu.

## 📊 Flow Diagram

```
[Frontend/Client]                    [gRPC Server]               [File System]
      |                                     |                          |
      | 1. User pilih file                 |                          |
      |-------------------------------->    |                          |
      |    Stream Start + Metadata         |                          |
      |                                     |                          |
      | 2. Kirim chunks (64KB each)        |                          |
      |-------------------------------->    |                          |
      |         chunk 1                     | Write to file            |
      |-------------------------------->    |---------------------->   |
      |         chunk 2                     | Write to file            |
      |-------------------------------->    |---------------------->   |
      |         chunk N                     | Write to file            |
      |-------------------------------->    |---------------------->   |
      |                                     |                          |
      | 3. Stream end                       |                          |
      |-------------------------------->    |                          |
      |                                     | Close file               |
      |                                     |---------------------->   |
      |                                     |                          |
      | 4. Response (URL)                   |                          |
      |<--------------------------------    |                          |
      |   {url: "/uploads/abc-123.jpg"}    |                          |
```

## ✅ Testing Checklist

- [ ] Server bisa terima metadata
- [ ] Server bisa terima chunks
- [ ] File tersimpan di `/uploads/`
- [ ] Response URL benar
- [ ] Validasi file type (hanya image)
- [ ] Validasi file size (max 5MB)
- [ ] Progress bar di frontend
- [ ] Error handling jika network terputus
- [ ] Foto tampil di list mobil

## 🎯 Untuk Presentasi

**Highlight:**
1. **gRPC Streaming** - efisien untuk upload file besar (tidak perlu Base64)
2. **Chunking** - 64KB per chunk, hemat memory
3. **Progress tracking** - user bisa lihat progres upload
4. **Type-safe** - proto contract untuk metadata & chunks
5. **Validasi server-side** - file type, size, dll

**Demo Script:**
```
1. Pilih foto mobil (max 5MB)
2. Tunjukkan progress bar (real-time chunking)
3. Foto tersimpan di /uploads/ dengan UUID filename
4. CreateMobil pakai URL hasil upload
5. Foto tampil di list mobil
```
