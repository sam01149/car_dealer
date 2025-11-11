'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation'; // <-- Impor useRouter
import { mobilClient, transaksiClient, addAuthMetadata } from '@/lib/grpcClient'; // <-- Impor transaksiClient & addAuthMetadata
import { GetMobilRequest, Mobil } from '@/proto/carapp_pb';
import { BuyMobilRequest, RentMobilRequest } from '@/proto/carapp_pb'; // <-- Impor Tipe Request
import { useAuth } from '@/context/AuthContext';

export default function MobilDetailPage() {
  const params = useParams(); // { id: 'uuid-mobil-abc' }
  const router = useRouter(); // Hook untuk redirect
  const { user } = useAuth(); // Ambil data user
  
  const [mobil, setMobil] = useState<Mobil | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // State baru untuk Transaksi
  const [txLoading, setTxLoading] = useState(false);
  const [txError, setTxError] = useState('');
  const [txSuccess, setTxSuccess] = useState('');

  // State untuk Pembayaran
  const [showPayment, setShowPayment] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');
  const [totalPrice, setTotalPrice] = useState(0);
  const [kembalian, setKembalian] = useState(0);

  const id = Array.isArray(params.id) ? params.id[0] : params.id;

  // Efek untuk mengambil data mobil (tetap sama)
  useEffect(() => {
    if (!id) return; // Tunggu sampai ID tersedia

    const fetchMobil = async () => {
      try {
        setLoading(true);
        const req = new GetMobilRequest();
        req.setMobilId(id);

        const response = await new Promise<Mobil>((resolve, reject) => {
          mobilClient.getMobil(req, {}, (err, response) => {
            if (err) {
              console.error('gRPC Error:', err);
              reject(err);
            } else {
              resolve(response!);
            }
          });
        });
        
        setMobil(response);
      } catch (err: any) {
        setError(`Gagal memuat mobil: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchMobil();
  }, [id]); // Jalankan ulang jika ID berubah

  // --- HANDLER BARU: Beli Mobil ---
  const handleBeli = async () => {
    if (!user) {
      router.push('/login');
      return;
    }
    if (!mobil || !id) return;

    // Tampilkan form pembayaran
    setTotalPrice(mobil.getHargaJual());
    setShowPayment(true);
    setPaymentAmount('');
    setKembalian(0);
  };

  // Handler Proses Pembayaran Beli
  const handleProcessPayment = async () => {
    const amount = parseFloat(paymentAmount);
    
    if (isNaN(amount) || amount <= 0) {
      setTxError('Masukkan jumlah pembayaran yang valid!');
      return;
    }

    if (amount < totalPrice) {
      setTxError(`Pembayaran kurang! Anda perlu membayar Rp ${totalPrice.toLocaleString('id-ID')}, tapi hanya membayar Rp ${amount.toLocaleString('id-ID')}`);
      return;
    }

    const change = amount - totalPrice;
    setKembalian(change);

    setTxLoading(true);
    setTxError('');
    setTxSuccess('');

    try {
      const req = new BuyMobilRequest();
      req.setMobilId(id!);

      // Buat metadata dengan token
      const metadata = addAuthMetadata({});

      // Panggil gRPC (dengan metadata yang sudah ada token)
      const response = await new Promise<any>((resolve, reject) => {
        transaksiClient.buyMobil(req, metadata, (err, response) => {
          if (err) {
            console.error('gRPC Error:', err);
            reject(err);
          } else {
            resolve(response);
          }
        });
      });

      setTxSuccess(`✅ Pembelian berhasil! (ID Transaksi: ${response.getId()})${change > 0 ? `\n💰 Kembalian Anda: Rp ${change.toLocaleString('id-ID')}` : ''}\n\nAnda akan diarahkan ke dashboard...`);
      setShowPayment(false);
      
      // Redirect ke dashboard setelah 4 detik
      setTimeout(() => {
        router.push('/dashboard');
      }, 4000);

    } catch (err: any) {
      setTxError(`Gagal membeli: ${err.message}`);
    } finally {
      setTxLoading(false);
    }
  };

  if (loading) {
    return <main className="container mx-auto p-8"><h1 className="text-2xl">Loading mobil...</h1></main>;
  }

  if (error || !mobil) {
    return <main className="container mx-auto p-8"><h1 className="text-2xl text-red-500">{error || 'Mobil tidak ditemukan.'}</h1></main>;
  }

  // Cek apakah user ini adalah pemiliknya
  const isOwner = user?.getId() === mobil.getOwnerId();
  // Mobil hanya bisa dibeli/dirental jika:
  // 1. Statusnya 'tersedia'
  // 2. Transaksi pembelian/rental belum sukses (txSuccess masih kosong)
  const isAvailable = mobil.getStatus() === 'tersedia' && !txSuccess;

  // Debug logs (comment out in production)
  console.log('=== DEBUG INFO ===');
  console.log('User:', user ? user.getId() : 'NOT LOGGED IN');
  console.log('Owner ID:', mobil.getOwnerId());
  console.log('isOwner:', isOwner);
  console.log('Status:', mobil.getStatus());
  console.log('isAvailable:', isAvailable);
  console.log('Harga Jual:', mobil.getHargaJual());
  console.log('Show Buy Button:', !isOwner && isAvailable && mobil.getHargaJual() > 0);
  console.log('txSuccess:', txSuccess);
  console.log('================');

  return (
    <main className="container mx-auto p-8">
      <div className="bg-white text-black p-8 rounded-lg shadow-lg">
        {/* Foto Mobil */}
        {mobil.getFotoUrl() && (
          <div className="mb-6 rounded-lg overflow-hidden">
            <img 
              src={mobil.getFotoUrl()} 
              alt={`${mobil.getMerk()} ${mobil.getModel()}`}
              className="w-full h-96 object-cover"
              onError={(e) => {
                e.currentTarget.style.display = 'none';
              }}
            />
          </div>
        )}

        {/* ... (Info Detail Mobil, tetap sama) ... */}
        <h1 className="text-5xl font-bold mb-4">
          {mobil.getTahun()} {mobil.getMerk()} {mobil.getModel()}
        </h1>
        
        <p className="text-3xl font-semibold text-green-700 mb-6">
          Rp {mobil.getHargaJual().toLocaleString('id-ID')}
        </p>
        
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div><strong>Kondisi:</strong> <span className="capitalize">{mobil.getKondisi()}</span></div>
          <div><strong>Lokasi:</strong> {mobil.getLokasi()}</div>
          <div><strong>Status:</strong> <span className="capitalize font-bold">{mobil.getStatus()}</span></div>
          <div><strong>Penjual:</strong> {isOwner ? "Anda" : (mobil.getOwnerName() || `User ${mobil.getOwnerId().substring(0, 8)}...`)}</div>
        </div>
        
        <p className="text-gray-700 mb-8">{mobil.getDeskripsi()}</p>

        {/* --- Area Notifikasi Transaksi --- */}
        <div className="my-4">
          {txLoading && <p className="text-blue-600">Memproses transaksi Anda...</p>}
          {txError && <p className="text-red-600">{txError}</p>}
          {txSuccess && <p className="text-green-600">{txSuccess}</p>}
        </div>

        {/* --- Tombol Aksi --- */}
        <div className="flex flex-col gap-4">
          {/* Tombol Beli */}
          {!isOwner && isAvailable && mobil.getHargaJual() > 0 && (
            <button
              onClick={handleBeli}
              disabled={txLoading} // Nonaktifkan tombol saat loading
              className="bg-blue-600 text-white px-6 py-3 rounded-lg font-bold text-lg hover:bg-blue-700 disabled:bg-gray-400"
            >
              {txLoading ? 'Memproses...' : 'Beli Mobil Ini'}
            </button>
          )}

          {isOwner && (
            <p className="text-lg italic text-gray-600">Ini adalah mobil Anda.</p>
          )}
        </div>

        {/* --- Modal Pembayaran --- */}
        {showPayment && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-8 rounded-lg shadow-2xl max-w-md w-full">
              <h2 className="text-2xl font-bold mb-4">💳 Pembayaran Pembelian</h2>
              
              <div className="mb-4">
                <p className="text-gray-700 mb-2">
                  <strong>Total yang harus dibayar:</strong>
                </p>
                <p className="text-3xl font-bold text-blue-600">
                  Rp {totalPrice.toLocaleString('id-ID')}
                </p>
              </div>

              <div className="mb-4">
                <label className="block text-gray-700 font-semibold mb-2">
                  Jumlah Pembayaran:
                </label>
                <input
                  type="number"
                  value={paymentAmount}
                  onChange={(e) => setPaymentAmount(e.target.value)}
                  placeholder="Masukkan jumlah pembayaran"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {kembalian > 0 && (
                <div className="mb-4 p-4 bg-green-100 border border-green-300 rounded-lg">
                  <p className="text-green-800 font-semibold">
                    💰 Kembalian Anda: <span className="text-2xl">Rp {kembalian.toLocaleString('id-ID')}</span>
                  </p>
                </div>
              )}

              {txError && (
                <div className="mb-4 p-4 bg-red-100 border border-red-300 rounded-lg">
                  <p className="text-red-800">{txError}</p>
                </div>
              )}

              <div className="flex gap-3">
                <button
                  onClick={handleProcessPayment}
                  disabled={txLoading}
                  className="flex-1 bg-blue-600 text-white py-3 rounded-lg font-bold hover:bg-blue-700 disabled:bg-gray-400"
                >
                  {txLoading ? 'Memproses...' : 'Proses Pembayaran'}
                </button>
                <button
                  onClick={() => {
                    setShowPayment(false);
                    setPaymentAmount('');
                    setKembalian(0);
                    setTxError('');
                  }}
                  disabled={txLoading}
                  className="flex-1 bg-gray-500 text-white py-3 rounded-lg font-bold hover:bg-gray-600 disabled:bg-gray-400"
                >
                  Batal
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
