-- User
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'client', -- client/admin
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Mobil
CREATE TABLE mobils (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID REFERENCES users(id),
    merk TEXT,
    model TEXT,
    tahun INT,
    kondisi TEXT, -- baru/bekas
    deskripsi TEXT,
    harga_jual NUMERIC,
    harga_rental_per_hari NUMERIC,
    lokasi TEXT,
    status TEXT DEFAULT 'tersedia', -- tersedia/terjual/dirental
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- TransaksiJual
CREATE TABLE transaksi_jual (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mobil_id UUID REFERENCES mobils(id),
    penjual_id UUID REFERENCES users(id),
    pembeli_id UUID REFERENCES users(id),
    total NUMERIC,
    status TEXT DEFAULT 'diproses', -- diproses/selesai/dibatalkan
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- TransaksiRental
CREATE TABLE transaksi_rental (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mobil_id UUID REFERENCES mobils(id),
    pemilik_id UUID REFERENCES users(id),
    penyewa_id UUID REFERENCES users(id),
    tanggal_mulai DATE,
    tanggal_selesai DATE,
    total NUMERIC,
    status TEXT DEFAULT 'aktif', -- aktif/selesai/dibatalkan
    denda_per_hari NUMERIC,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Notifikasi
CREATE TABLE notifikasi (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    tipe TEXT, -- jual/beli/rental/info
    pesan TEXT,
    priority TEXT DEFAULT 'normal',
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Cache Merk/Model
CREATE TABLE brand_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id TEXT UNIQUE,
    name TEXT,
    raw JSONB,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE model_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id TEXT UNIQUE,
    brand_id TEXT,
    name TEXT,
    raw JSONB,
    updated_at TIMESTAMP DEFAULT NOW()
);