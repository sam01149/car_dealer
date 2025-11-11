-- Hapus tabel transaksi_rental karena fitur rental dihapus
DROP TABLE IF EXISTS transaksi_rental;

-- Hapus kolom harga_rental_per_hari dari tabel mobils
ALTER TABLE mobils DROP COLUMN IF EXISTS harga_rental_per_hari;

-- Tambah kolom foto_url untuk menyimpan foto mobil (base64 atau URL)
ALTER TABLE mobils ADD COLUMN foto_url TEXT;

-- Update status mobil yang mungkin masih 'dirental' menjadi 'tersedia'
UPDATE mobils SET status = 'tersedia' WHERE status = 'dirental';
