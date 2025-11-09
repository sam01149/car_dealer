'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/src/context/AuthContext';
import { useRouter } from 'next/navigation';
import { dashboardClient, addAuthMetadata } from '@/lib/grpcClient';
import { DashboardSummary } from '@/proto/carapp_pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export default function DashboardPage() {
  const { user, logout, isLoading } = useAuth();
  const router = useRouter();
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    // Jika loading selesai dan user tidak ada, tendang ke login
    if (!isLoading && !user) {
      router.push('/login');
      return;
    }

    // Jika user ada, panggil dashboard
    if (user) {
      const getDashboardData = async () => {
        try {
          // Buat request Kosong (Empty)
          const req = new Empty();
          
          // Buat metadata dengan token authorization
          const metadata = addAuthMetadata({});
          console.log('Sending metadata:', metadata);
          console.log('Token from localStorage:', localStorage.getItem('authToken'));
          
          // Panggil gRPC (dengan metadata yang sudah ada token)
          const response = await new Promise<DashboardSummary>((resolve, reject) => {
            dashboardClient.getDashboard(req, metadata, (err, response) => {
              if (err) {
                console.error('gRPC Error:', err);
                console.error('Error code:', err.code);
                console.error('Error message:', err.message);
                reject(err);
              } else {
                resolve(response!);
              }
            });
          });
          
          setSummary(response);

        } catch (err: any) {
          setError(`Gagal memuat dashboard: ${err.message}`);
          console.error(err);
        }
      };

      getDashboardData();
    }
  }, [user, isLoading, router]);

  // Tampilkan loading saat mengecek auth atau fetch data
  if (isLoading || !summary) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center p-24">
        <h1 className="text-3xl font-bold text-gray-800">Loading Dashboard...</h1>
      </main>
    );
  }

  // Tampilkan error jika ada
  if (error) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center p-24">
        <h1 className="text-3xl font-bold text-red-500">{error}</h1>
      </main>
    );
  }

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8 bg-gray-100">
      <div className="w-full max-w-6xl">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>
              <p className="text-gray-600 mt-2">
                Selamat datang, <span className="font-semibold text-green-600">{user?.getName()}</span>!
              </p>
            </div>
            <button
              onClick={handleLogout}
              className="px-6 py-2 bg-red-600 text-white rounded-lg font-semibold hover:bg-red-700 transition-colors duration-200 shadow-md hover:shadow-lg"
            >
              Logout
            </button>
          </div>
        </div>

        {/* Dashboard Summary - Data dari gRPC */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">Ringkasan Dashboard</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-blue-50 rounded-lg border-2 border-blue-200">
              <h2 className="text-xl font-bold text-gray-800">Total Mobil Anda</h2>
              <p className="text-4xl font-bold text-blue-600">{summary.getTotalMobilAnda()}</p>
            </div>
            <div className="p-6 bg-purple-50 rounded-lg border-2 border-purple-200">
              <h2 className="text-xl font-bold text-gray-800">Notifikasi Baru</h2>
              <p className="text-4xl font-bold text-purple-600">{summary.getNotifikasiBaru()}</p>
            </div>
            <div className="p-6 bg-green-50 rounded-lg border-2 border-green-200">
              <h2 className="text-xl font-bold text-gray-800">Transaksi Aktif</h2>
              <p className="text-4xl font-bold text-green-600">{summary.getTransaksiAktif()}</p>
            </div>
            <div className="p-6 bg-orange-50 rounded-lg border-2 border-orange-200">
              <h2 className="text-xl font-bold text-gray-800">Total Pendapatan (Jual)</h2>
              <p className="text-4xl font-bold text-orange-600">
                Rp {summary.getPendapatanTerakhir().toLocaleString('id-ID')}
              </p>
            </div>
          </div>
        </div>

        {/* User Info Card */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">Informasi Akun</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">Nama Lengkap</p>
              <p className="text-lg font-semibold text-gray-800">{user?.getName()}</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">Email</p>
              <p className="text-lg font-semibold text-gray-800">{user?.getEmail()}</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">No. Telepon</p>
              <p className="text-lg font-semibold text-gray-800">{user?.getPhone()}</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">Alamat</p>
              <p className="text-lg font-semibold text-gray-800">-</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">Role</p>
              <p className="text-lg font-semibold text-gray-800 capitalize">{user?.getRole()}</p>
            </div>
            <div className="p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-gray-600">User ID</p>
              <p className="text-lg font-semibold text-gray-800">{user?.getId()}</p>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">Menu Cepat</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <button 
              onClick={() => router.push('/')}
              className="p-6 bg-blue-50 border-2 border-blue-200 rounded-lg hover:bg-blue-100 transition-colors duration-200 text-left"
            >
              <div className="text-2xl mb-2">🚗</div>
              <h3 className="text-lg font-semibold text-gray-800">Lihat Mobil</h3>
              <p className="text-sm text-gray-600">Jelajahi koleksi mobil</p>
            </button>
            <button 
              onClick={() => router.push('/mobil/jual')}
              className="p-6 bg-green-50 border-2 border-green-200 rounded-lg hover:bg-green-100 transition-colors duration-200 text-left"
            >
              <div className="text-2xl mb-2">💵</div>
              <h3 className="text-lg font-semibold text-gray-800">Jual Mobil</h3>
              <p className="text-sm text-gray-600">Pasang iklan mobil Anda</p>
            </button>
            <button 
              onClick={() => router.push('/notifikasi')}
              className="p-6 bg-purple-50 border-2 border-purple-200 rounded-lg hover:bg-purple-100 transition-colors duration-200 text-left"
            >
              <div className="text-2xl mb-2">🔔</div>
              <h3 className="text-lg font-semibold text-gray-800">Notifikasi</h3>
              <p className="text-sm text-gray-600">Lihat riwayat transaksi</p>
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
