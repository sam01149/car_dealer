'use client';

import { useEffect, useState } from 'react';
import { mobilClient } from '@/lib/grpcClient';
import { ListMobilRequest, Mobil } from '@/proto/carapp_pb';
import Link from 'next/link';

export default function HomePage() {
  const [mobils, setMobils] = useState<Mobil[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchMobils = async () => {
      try {
        setLoading(true);
        // Buat request. Kita bisa atur filter di sini jika mau
        // Untuk saat ini, kita ambil default (status = 'tersedia')
        const req = new ListMobilRequest();
        
        // Panggil gRPC (ini rute publik, interceptor akan mengizinkan)
        const response = await new Promise<any>((resolve, reject) => {
          mobilClient.listMobil(req, {}, (err, response) => {
            if (err) {
              console.error('gRPC Error:', err);
              reject(err);
            } else {
              resolve(response);
            }
          });
        });
        
        console.log('=== HOMEPAGE DEBUG ===');
        console.log('Mobils fetched:', response.getMobilsList().length);
        console.log('Mobils:', response.getMobilsList().map((m: Mobil) => ({
          id: m.getId(),
          merk: m.getMerk(),
          model: m.getModel(),
          status: m.getStatus(),
          ownerId: m.getOwnerId(),
          fotoUrl: m.getFotoUrl(),
          hargaRental: m.getHargaRentalPerHari()
        })));
        console.log('First car foto_url:', response.getMobilsList()[0]?.getFotoUrl());
        console.log('=====================');
        
        setMobils(response.getMobilsList());
      } catch (err: any) {
        setError(`Gagal memuat mobil: ${err.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchMobils();
  }, []);

  if (loading) {
    return (
      <main className="container mx-auto p-8">
        <h1 className="text-3xl font-bold text-gray-800">Mencari mobil...</h1>
      </main>
    );
  }

  if (error) {
    return (
      <main className="container mx-auto p-8">
        <h1 className="text-3xl font-bold text-red-500">{error}</h1>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-16 px-8 mb-8">
        <div className="container mx-auto">
          <h1 className="text-5xl font-bold mb-4 animate-fade-in">🚗 Selamat Datang di CarApp</h1>
          <p className="text-xl opacity-90">Platform jual-beli dan rental mobil terpercaya</p>
          <p className="text-lg mt-2 opacity-80">Temukan mobil impian Anda dengan harga terbaik</p>
        </div>
      </div>

      {/* Stats Section */}
      <div className="container mx-auto px-8 mb-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white rounded-xl shadow-lg p-6 text-center transform hover:scale-105 transition-transform">
            <div className="text-4xl mb-2">🏆</div>
            <h3 className="text-2xl font-bold text-gray-800">{mobils.length}+</h3>
            <p className="text-gray-600">Mobil Tersedia</p>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-6 text-center transform hover:scale-105 transition-transform">
            <div className="text-4xl mb-2">✅</div>
            <h3 className="text-2xl font-bold text-gray-800">100%</h3>
            <p className="text-gray-600">Mobil Terverifikasi</p>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-6 text-center transform hover:scale-105 transition-transform">
            <div className="text-4xl mb-2">🔒</div>
            <h3 className="text-2xl font-bold text-gray-800">Aman</h3>
            <p className="text-gray-600">Transaksi Terjamin</p>
          </div>
        </div>
      </div>

      {/* Car Showcase */}
      <div className="container mx-auto px-8 pb-12">
        <div className="flex items-center justify-between mb-8">
          <h2 className="text-4xl font-bold text-gray-800">🚘 Etalase Mobil Premium</h2>
          <div className="flex gap-3">
            <span className="px-4 py-2 bg-blue-100 text-blue-700 rounded-full text-sm font-semibold">Baru</span>
            <span className="px-4 py-2 bg-green-100 text-green-700 rounded-full text-sm font-semibold">Bekas</span>
          </div>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {mobils.length === 0 ? (
            <div className="col-span-full text-center py-12">
              <div className="text-6xl mb-4">🚗</div>
              <p className="text-xl text-gray-600">Belum ada mobil yang dijual.</p>
              <p className="text-gray-500 mt-2">Jadilah yang pertama untuk menjual mobil!</p>
            </div>
          ) : (
            mobils.map((mobil) => (
              <Link
                href={`/mobil/${mobil.getId()}`}
                key={mobil.getId()}
                className="group relative bg-white rounded-2xl shadow-lg overflow-hidden transform hover:scale-105 hover:shadow-2xl transition-all duration-300"
              >
                {/* Badge Status */}
                <div className="absolute top-4 right-4 z-10">
                  <span className={`px-3 py-1 rounded-full text-xs font-bold uppercase ${
                    mobil.getKondisi() === 'baru' 
                      ? 'bg-blue-500 text-white' 
                      : 'bg-green-500 text-white'
                  }`}>
                    {mobil.getKondisi()}
                  </span>
                </div>

                {/* Image Placeholder with Gradient */}
                <div className="h-48 bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center relative overflow-hidden">
                  {mobil.getFotoUrl() ? (
                    <img 
                      src={mobil.getFotoUrl()} 
                      alt={`${mobil.getMerk()} ${mobil.getModel()}`}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                      onError={(e) => {
                        // Fallback jika gambar gagal load
                        e.currentTarget.style.display = 'none';
                        e.currentTarget.parentElement!.innerHTML = '<div class="text-6xl opacity-90">🚗</div>';
                      }}
                    />
                  ) : (
                    <div className="text-6xl opacity-90 group-hover:scale-110 transition-transform">
                      🚗
                    </div>
                  )}
                  <div className="absolute inset-0 bg-black opacity-0 group-hover:opacity-10 transition-opacity"></div>
                </div>

                {/* Content */}
                <div className="p-5">
                  <div className="mb-3">
                    <h3 className="text-xl font-bold text-gray-800 group-hover:text-blue-600 transition-colors line-clamp-1">
                      {mobil.getMerk()} {mobil.getModel()}
                    </h3>
                    <p className="text-sm text-gray-500">{mobil.getTahun()}</p>
                  </div>

                  <div className="mb-4">
                    <p className="text-2xl font-bold text-blue-600">
                      Rp {mobil.getHargaJual().toLocaleString('id-ID')}
                    </p>
                  </div>

                  <div className="flex items-center justify-between text-sm text-gray-600 mb-4">
                    <div className="flex items-center gap-1">
                      <span>📍</span>
                      <span className="truncate">{mobil.getLokasi()}</span>
                    </div>
                  </div>

                  {/* CTA Button */}
                  <button className="w-full bg-gradient-to-r from-blue-500 to-purple-500 text-white py-3 rounded-lg font-semibold group-hover:from-blue-600 group-hover:to-purple-600 transition-all">
                    Lihat Detail →
                  </button>
                </div>

                {/* Hover Overlay Effect */}
                <div className="absolute inset-0 border-4 border-blue-500 opacity-0 group-hover:opacity-100 transition-opacity rounded-2xl pointer-events-none"></div>
              </Link>
            ))
          )}
        </div>
      </div>

      {/* Call to Action Footer */}
      <div className="bg-gradient-to-r from-purple-600 to-blue-600 text-white py-12 mt-12">
        <div className="container mx-auto px-8 text-center">
          <h3 className="text-3xl font-bold mb-4">Punya mobil untuk dijual?</h3>
          <p className="text-lg mb-6 opacity-90">Jual mobil Anda dengan mudah dan cepat di CarApp!</p>
          <Link 
            href="/mobil/jual" 
            className="inline-block bg-white text-blue-600 px-8 py-3 rounded-full font-bold text-lg hover:bg-gray-100 transition-colors shadow-lg"
          >
            Jual Mobil Sekarang →
          </Link>
        </div>
      </div>
    </main>
  );
}

