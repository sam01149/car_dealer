'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { notifikasiClient, addAuthMetadata } from '@/lib/grpcClient';
import { GetNotificationsRequest, Notifikasi } from '@/proto/carapp_pb';

export default function NotifikasiPage() {
  const router = useRouter();
  const { user, isLoading } = useAuth();
  const [notifikasis, setNotifikasis] = useState<Notifikasi[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login?redirect=/notifikasi');
      return;
    }

    if (!user || isLoading) return; // Tunggu sampai user fully loaded

    const fetchNotifikasi = async () => {
      try {
        setLoading(true);
        console.log('🔔 Fetching notifications for user:', user.getId());
        
        // Dapatkan token dari localStorage
        const token = localStorage.getItem('authToken');
        console.log('🔑 Token available:', token ? 'YES' : 'NO');
        
        if (!token) {
          setError('Token tidak ditemukan. Silakan login kembali.');
          setLoading(false);
          return;
        }
        
        const req = new GetNotificationsRequest();
        
        // Untuk streaming gRPC-Web, metadata harus berupa object dengan key lowercase
        const metadata = {
          'authorization': `Bearer ${token}`
        };
        
        console.log('📡 Metadata prepared:', metadata);
        console.log('📡 Starting stream...');

        const stream = notifikasiClient.getNotifications(req, metadata);
        const notifs: Notifikasi[] = [];

        stream.on('data', (notif: Notifikasi) => {
          console.log('📩 Received notification:', notif.getPesan());
          notifs.push(notif);
        });

        stream.on('end', () => {
          console.log('✅ Stream ended. Total notifications:', notifs.length);
          setNotifikasis(notifs);
          setLoading(false);
        });

        stream.on('error', (err: any) => {
          console.error('❌ Stream Error:', err);
          console.error('Error details:', err.message, err.code);
          console.error('Full error object:', err);
          
          // Cek apakah ini masalah auth
          if (err.message && err.message.includes('UserID dari token')) {
            setError(`Gagal memuat notifikasi: ${err.message}\n\n⚠️ Backend mungkin belum di-restart. Pastikan backend sudah berjalan dengan code terbaru (go run main.go)`);
          } else {
            setError(`Gagal memuat notifikasi: ${err.message}`);
          }
          setLoading(false);
        });
      } catch (err: any) {
        console.error('❌ Fetch Error:', err);
        setError(`Error: ${err.message}`);
        setLoading(false);
      }
    };

    fetchNotifikasi();
  }, [user, isLoading, router]);

  if (loading || isLoading) {
    return (
      <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 p-8">
        <div className="container mx-auto">
          <h1 className="text-3xl font-bold text-gray-800 mb-6">🔔 Notifikasi</h1>
          <p className="text-gray-600">Memuat notifikasi...</p>
        </div>
      </main>
    );
  }

  if (error) {
    return (
      <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 p-8">
        <div className="container mx-auto max-w-4xl">
          <h1 className="text-3xl font-bold text-gray-800 mb-6">🔔 Notifikasi</h1>
          <div className="bg-red-50 border-2 border-red-300 rounded-lg p-6">
            <h2 className="text-xl font-bold text-red-800 mb-2">❌ Error</h2>
            <p className="text-red-600 mb-4">{error}</p>
            <div className="bg-white rounded p-4 mb-4">
              <p className="text-sm text-gray-700 mb-2"><strong>Troubleshooting:</strong></p>
              <ul className="list-disc list-inside text-sm text-gray-600 space-y-1">
                <li>Pastikan Anda sudah login</li>
                <li>Coba logout dan login kembali</li>
                <li>Pastikan backend server berjalan di port 9090</li>
                <li>Cek console browser untuk detail error</li>
              </ul>
            </div>
            <button
              onClick={() => router.push('/dashboard')}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Kembali ke Dashboard
            </button>
          </div>
        </div>
      </main>
    );
  }

  // Group notifications by type
  const getIconForType = (tipe: string) => {
    switch (tipe) {
      case 'beli': return '🛒';
      case 'jual': return '💰';
      case 'rental': return '🔑';
      default: return '📬';
    }
  };

  const getColorForType = (tipe: string) => {
    switch (tipe) {
      case 'beli': return 'from-blue-500 to-blue-600';
      case 'jual': return 'from-green-500 to-green-600';
      case 'rental': return 'from-purple-500 to-purple-600';
      default: return 'from-gray-500 to-gray-600';
    }
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 p-8">
      <div className="container mx-auto max-w-4xl">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">🔔 Notifikasi</h1>
          <p className="text-gray-600">Riwayat aktivitas dan transaksi Anda</p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div className="bg-white rounded-xl shadow-lg p-4 text-center">
            <div className="text-2xl mb-1">📊</div>
            <p className="text-2xl font-bold text-gray-800">{notifikasis.length}</p>
            <p className="text-sm text-gray-600">Total Notifikasi</p>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-4 text-center">
            <div className="text-2xl mb-1">🛒</div>
            <p className="text-2xl font-bold text-blue-600">
              {notifikasis.filter(n => n.getTipe() === 'beli').length}
            </p>
            <p className="text-sm text-gray-600">Pembelian</p>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-4 text-center">
            <div className="text-2xl mb-1">💰</div>
            <p className="text-2xl font-bold text-green-600">
              {notifikasis.filter(n => n.getTipe() === 'jual').length}
            </p>
            <p className="text-sm text-gray-600">Penjualan</p>
          </div>
        </div>

        {/* Notifications List */}
        <div className="space-y-4">
          {notifikasis.length === 0 ? (
            <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
              <div className="text-6xl mb-4">📭</div>
              <h3 className="text-2xl font-bold text-gray-800 mb-2">Belum Ada Notifikasi</h3>
              <p className="text-gray-600">Notifikasi Anda akan muncul di sini setelah melakukan transaksi</p>
            </div>
          ) : (
            notifikasis.map((notif) => {
              const createdAt = notif.getCreatedAt();
              const date = createdAt ? new Date(createdAt.getSeconds() * 1000) : new Date();
              
              return (
                <div
                  key={notif.getId()}
                  className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:scale-[1.02]"
                >
                  <div className="flex">
                    {/* Icon Section */}
                    <div className={`w-20 bg-gradient-to-br ${getColorForType(notif.getTipe())} flex items-center justify-center`}>
                      <span className="text-4xl">{getIconForType(notif.getTipe())}</span>
                    </div>
                    
                    {/* Content Section */}
                    <div className="flex-1 p-5">
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase text-white bg-gradient-to-r ${getColorForType(notif.getTipe())}`}>
                            {notif.getTipe()}
                          </span>
                          {notif.getPriority() === 'urgent' && (
                            <span className="px-2 py-1 rounded-full text-xs font-bold bg-red-500 text-white">
                              🔥 URGENT
                            </span>
                          )}
                        </div>
                        <span className="text-sm text-gray-500">
                          {date.toLocaleDateString('id-ID', { 
                            day: '2-digit', 
                            month: 'short', 
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </span>
                      </div>
                      
                      <p className="text-gray-800 leading-relaxed">
                        {notif.getPesan()}
                      </p>
                    </div>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Back Button */}
        <div className="mt-8 text-center">
          <button
            onClick={() => router.push('/dashboard')}
            className="px-6 py-3 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-lg font-semibold hover:from-blue-600 hover:to-purple-600 transition-all shadow-lg hover:shadow-xl"
          >
            ← Kembali ke Dashboard
          </button>
        </div>
      </div>
    </main>
  );
}
