'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/src/context/AuthContext';
import { useRouter } from 'next/navigation';
import { mobilClient, nhtsaClient, addAuthMetadata } from '@/lib/grpcClient';
import { CreateMobilRequest, Make, GetModelsForMakeRequest } from '@/proto/carapp_pb';
import { GetMakesRequest } from '@/proto/carapp_pb';

export default function JualMobilPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  // State untuk data dari NHTSA API
  const [makes, setMakes] = useState<Make[]>([]);
  const [selectedMake, setSelectedMake] = useState('');
  const [selectedMakeId, setSelectedMakeId] = useState('');
  const [models, setModels] = useState<any[]>([]);
  const [loadingModels, setLoadingModels] = useState(false);
  
  // State untuk search merek
  const [searchMake, setSearchMake] = useState('');
  const [showMakeDropdown, setShowMakeDropdown] = useState(false);
  const [filteredMakes, setFilteredMakes] = useState<Make[]>([]);
  
  // State untuk form
  const [model, setModel] = useState('');
  const [tahun, setTahun] = useState(2020);
  const [kondisi, setKondisi] = useState('baru');
  const [hargaJual, setHargaJual] = useState(0);
  const [lokasi, setLokasi] = useState('');
  const [deskripsi, setDeskripsi] = useState('');
  const [fotoUrl, setFotoUrl] = useState('');

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // 1. Cek Autentikasi
  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login?redirect=/mobil/jual');
    }
  }, [user, isLoading, router]);

  // 2. Ambil data Merek Mobil (NHTSA)
  useEffect(() => {
    const fetchMakes = async () => {
      try {
        const req = new GetMakesRequest();
        const response = await new Promise<any>((resolve, reject) => {
          nhtsaClient.getMakes(req, {}, (err, response) => {
            if (err) {
              console.error('gRPC Error:', err);
              reject(err);
            } else {
              resolve(response);
            }
          });
        });
        setMakes(response.getMakesList());
      } catch (err: any) {
        setError(`Gagal memuat merek mobil: ${err.message}`);
      }
    };
    fetchMakes();
  }, []);

  // 2b. Ambil data Model berdasarkan Merek yang dipilih
  useEffect(() => {
    if (!selectedMakeId) {
      setModels([]);
      setModel('');
      return;
    }

    const fetchModels = async () => {
      try {
        setLoadingModels(true);
        const req = new GetModelsForMakeRequest();
        req.setBrandId(selectedMakeId);
        
        const response = await new Promise<any>((resolve, reject) => {
          nhtsaClient.getModelsForMake(req, {}, (err, response) => {
            if (err) {
              console.error('gRPC Error:', err);
              reject(err);
            } else {
              resolve(response);
            }
          });
        });
        
        setModels(response.getModelsList());
      } catch (err: any) {
        console.error('Gagal memuat model:', err);
        setModels([]);
      } finally {
        setLoadingModels(false);
      }
    };

    fetchModels();
  }, [selectedMakeId]);

  // Filter merek berdasarkan search
  useEffect(() => {
    if (searchMake.trim() === '') {
      setFilteredMakes([]);
      setShowMakeDropdown(false);
    } else {
      const filtered = makes.filter(make => 
        make.getName().toLowerCase().includes(searchMake.toLowerCase())
      );
      setFilteredMakes(filtered);
      setShowMakeDropdown(true);
    }
  }, [searchMake, makes]);

  // 3. Handle Submit Form
  const handleSubmit = async () => {
    setError('');
    setSuccess('');
    
    if (!selectedMake || !model || tahun <= 1900 || hargaJual <= 0) {
      setError('Data tidak valid. Mohon cek Merek, Model, Tahun, dan Harga.');
      return;
    }

    try {
      const req = new CreateMobilRequest();
      req.setMerk(selectedMake);
      req.setModel(model);
      req.setTahun(tahun);
      req.setKondisi(kondisi);
      req.setHargaJual(hargaJual);
      req.setLokasi(lokasi);
      req.setDeskripsi(deskripsi);
      req.setFotoUrl(fotoUrl);
      // 'harga_rental_per_hari' bisa di-skip (default 0)

      // Buat metadata dengan token
      const metadata = addAuthMetadata({});

      // Panggil gRPC (dengan metadata yang sudah ada token)
      const response = await new Promise<any>((resolve, reject) => {
        mobilClient.createMobil(req, metadata, (err, response) => {
          if (err) {
            console.error('gRPC Error:', err);
            reject(err);
          } else {
            resolve(response);
          }
        });
      });
      
      setSuccess(`Mobil ${response.getMerk()} berhasil dibuat!`);
      // Arahkan ke halaman detail mobil yang baru dibuat
      setTimeout(() => {
        router.push(`/mobil/${response.getId()}`);
      }, 2000);

    } catch (err: any) {
      setError(`Gagal membuat mobil: ${err.message}`);
    }
  };
  
  if (isLoading || !user) {
    return <main className="container mx-auto p-8 bg-gray-50 min-h-screen"><h1 className="text-2xl text-gray-800">Loading...</h1></main>;
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-gray-100 py-12">
      <div className="container mx-auto px-4">
        <div className="max-w-2xl mx-auto">
          <h1 className="text-4xl font-bold mb-8 text-gray-800 text-center">🚗 Jual Mobil Anda</h1>
          
          <div className="bg-white rounded-xl shadow-lg p-8">
            {error && (
              <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg">
                <strong>Error:</strong> {error}
              </div>
            )}
            {success && (
              <div className="mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg">
                <strong>Berhasil!</strong> {success}
              </div>
            )}

            <div className="flex flex-col gap-6">
              {/* Search Merek dari API */}
              <div className="relative">
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Merek Mobil <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  placeholder="Ketik untuk mencari merek... (contoh: Toyota, Honda)"
                  value={searchMake}
                  onChange={(e) => setSearchMake(e.target.value)}
                  onFocus={() => {
                    if (searchMake.trim() !== '') {
                      setShowMakeDropdown(true);
                    }
                  }}
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                {selectedMake && (
                  <div className="mt-2 p-2 bg-blue-50 border border-blue-200 rounded-lg flex items-center justify-between">
                    <span className="text-sm text-blue-800">
                      <strong>Terpilih:</strong> {selectedMake}
                    </span>
                    <button
                      onClick={() => {
                        setSelectedMake('');
                        setSelectedMakeId('');
                        setSearchMake('');
                        setShowMakeDropdown(false);
                      }}
                      className="text-red-600 hover:text-red-800 font-bold"
                    >
                      ✕
                    </button>
                  </div>
                )}
                
                {/* Dropdown hasil search */}
                {showMakeDropdown && filteredMakes.length > 0 && (
                  <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {filteredMakes.slice(0, 20).map((make) => (
                      <button
                        key={make.getBrandId()}
                        type="button"
                        onClick={() => {
                          setSelectedMake(make.getName());
                          setSelectedMakeId(make.getBrandId());
                          setSearchMake(make.getName());
                          setShowMakeDropdown(false);
                        }}
                        className="w-full text-left p-3 hover:bg-blue-50 border-b border-gray-100 last:border-b-0 transition-colors"
                      >
                        <span className="text-gray-800">{make.getName()}</span>
                      </button>
                    ))}
                    {filteredMakes.length > 20 && (
                      <div className="p-3 text-sm text-gray-500 text-center border-t border-gray-200">
                        Menampilkan 20 dari {filteredMakes.length} hasil. Ketik lebih spesifik untuk mempersempit pencarian.
                      </div>
                    )}
                  </div>
                )}
                
                {/* Pesan jika tidak ada hasil */}
                {showMakeDropdown && filteredMakes.length === 0 && searchMake.trim() !== '' && (
                  <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg p-3">
                    <p className="text-gray-500 text-sm">Tidak ada merek yang cocok dengan "{searchMake}"</p>
                  </div>
                )}
              </div>
              
              {/* Model */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Model <span className="text-red-500">*</span>
                </label>
                {!selectedMake ? (
                  <div className="w-full p-3 border border-gray-300 rounded-lg bg-gray-100 text-gray-500">
                    Pilih merek terlebih dahulu
                  </div>
                ) : loadingModels ? (
                  <div className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-600">
                    Memuat model...
                  </div>
                ) : models.length > 0 ? (
                  <select
                    value={model}
                    onChange={(e) => setModel(e.target.value)}
                    className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="">-- Pilih Model --</option>
                    {models.map((m) => (
                      <option key={m.getModelId()} value={m.getName()}>
                        {m.getName()}
                      </option>
                    ))}
                  </select>
                ) : (
                  <input 
                    type="text" 
                    placeholder="Model tidak tersedia, ketik manual" 
                    value={model} 
                    onChange={(e) => setModel(e.target.value)} 
                    className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                  />
                )}
              </div>

              {/* Tahun */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Tahun <span className="text-red-500">*</span>
                </label>
                <input 
                  type="number" 
                  placeholder="Contoh: 2020" 
                  value={tahun} 
                  onChange={(e) => setTahun(parseInt(e.target.value))} 
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                />
              </div>

              {/* Harga */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Harga Jual (Rp) <span className="text-red-500">*</span>
                </label>
                <input 
                  type="number" 
                  placeholder="Contoh: 250000000" 
                  value={hargaJual} 
                  onChange={(e) => {
                    const value = e.target.value;
                    // Bulatkan ke integer untuk menghindari floating point precision issue
                    setHargaJual(value ? Math.round(parseFloat(value)) : 0);
                  }}
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                />
              </div>

              {/* Kondisi */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Kondisi
                </label>
                <select 
                  value={kondisi} 
                  onChange={(e) => setKondisi(e.target.value)} 
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="baru">Baru</option>
                  <option value="bekas">Bekas</option>
                </select>
              </div>

              {/* Lokasi */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Lokasi
                </label>
                <input 
                  type="text" 
                  placeholder="Contoh: Jakarta, Surabaya, Bandung" 
                  value={lokasi} 
                  onChange={(e) => setLokasi(e.target.value)} 
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                />
              </div>

              {/* Deskripsi */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Deskripsi <span className="text-red-500">*</span>
                </label>
                <textarea 
                  placeholder="jelaskan singkat tentang mobil anda dan sertakan plat mobil" 
                  value={deskripsi} 
                  onChange={(e) => setDeskripsi(e.target.value)} 
                  required
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                  rows={4} 
                />
              </div>

              {/* Foto URL */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  URL Foto Mobil <span className="text-gray-500">(Opsional)</span>
                </label>
                <input 
                  type="url" 
                  placeholder="https://example.com/foto-mobil.jpg" 
                  value={fotoUrl} 
                  onChange={(e) => setFotoUrl(e.target.value)} 
                  className="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" 
                />
                <p className="mt-2 text-sm text-gray-500">
                  💡 Masukkan URL foto mobil dari internet (contoh: dari Google Drive, Imgur, atau website lain)
                </p>
                
                {/* Preview foto jika URL valid */}
                {fotoUrl && (
                  <div className="mt-4 border border-gray-300 rounded-lg p-4 bg-gray-50">
                    <p className="text-sm font-semibold text-gray-700 mb-2">Preview Foto:</p>
                    <img 
                      src={fotoUrl} 
                      alt="Preview mobil" 
                      className="w-full max-h-64 object-cover rounded-lg"
                      onError={(e) => {
                        e.currentTarget.src = '';
                        e.currentTarget.alt = '❌ URL foto tidak valid';
                        e.currentTarget.className = 'w-full p-4 text-center text-red-500 bg-red-50 rounded-lg';
                      }}
                    />
                  </div>
                )}
              </div>
              
              <button 
                onClick={handleSubmit} 
                className="w-full p-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-bold text-lg transition-colors duration-200 shadow-lg hover:shadow-xl"
              >
                📤 Submit Iklan Mobil
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
