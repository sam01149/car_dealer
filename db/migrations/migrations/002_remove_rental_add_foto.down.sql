-- Rollback: Hapus kolom foto_url
ALTER TABLE mobils DROP COLUMN IF EXISTS foto_url;

-- Rollback: Tambah kembali kolom harga_rental_per_hari
ALTER TABLE mobils ADD COLUMN harga_rental_per_hari NUMERIC;

-- Rollback: Buat kembali tabel transaksi_rental
CREATE TABLE IF NOT EXISTS transaksi_rental (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mobil_id UUID REFERENCES mobils(id),
    pemilik_id UUID REFERENCES users(id),
    penyewa_id UUID REFERENCES users(id),
    tanggal_mulai DATE,
    tanggal_selesai DATE,
    total NUMERIC,
    status TEXT DEFAULT 'aktif',
    denda_per_hari NUMERIC,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
