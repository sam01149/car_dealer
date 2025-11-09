'use client'; // <-- PENTING! Ini menandakan Komponen Klien

import { useState } from 'react';
import { RegisterRequest } from '@/proto/carapp_pb'; // <-- Impor Tipe Request
import { authClient } from '@/lib/grpcClient'; // <-- Impor Klien gRPC kita
import { useAuth } from '@/src/context/AuthContext';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function RegisterPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phone, setPhone] = useState('');

  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  const { login } = useAuth();
  const router = useRouter();

  const handleRegister = async () => {
    setError('');
    setSuccess('');

    // Validasi input
    if (!name || !email || !password) {
      setError('Nama, Email, dan Password wajib diisi!');
      return;
    }

    // 1. Buat request gRPC (menggunakan Tipe dari file .pb.js)
    const req = new RegisterRequest();
    req.setName(name);
    req.setEmail(email);
    req.setPassword(password);
    req.setPhone(phone);

    try {
      // 2. Panggil gRPC client (ini adalah async call)
      // Metadata kosong untuk public endpoint
      const metadata = {};
      
      const response = await new Promise((resolve, reject) => {
        authClient.register(req, metadata, (err, response) => {
          if (err) {
            console.error('gRPC Error:', err);
            console.error('Error code:', err.code);
            console.error('Error message:', err.message);
            reject(err);
          } else {
            resolve(response);
          }
        });
      });

      // 3. Tangani respons
      const user = (response as any).getUser();
      const token = (response as any).getToken();

      if (user) {
        setSuccess(
          `Registrasi berhasil! Halo, ${user.getName()}!`
        );
        
        // Login otomatis setelah registrasi
        login(token, user);
        
        // Reset form
        setName('');
        setEmail('');
        setPassword('');
        setPhone('');
        
        // Redirect ke dashboard setelah 2 detik
        setTimeout(() => {
          router.push('/dashboard');
        }, 2000);
      }
    } catch (err: any) {
      // 4. Tangani error gRPC
      console.error('Error detail:', err);
      const errorMessage = err.message || 'Terjadi kesalahan saat menghubungi server';
      setError(`Registrasi gagal: ${errorMessage}`);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gray-100">
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-3xl font-bold mb-6 text-center text-gray-800">Registrasi Akun</h1>
        <div className="flex flex-col gap-4 w-full">
          {/* Tampilkan pesan error atau sukses */}
          {error && (
            <div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}
          {success && (
            <div className="p-3 bg-green-100 border border-green-400 text-green-700 rounded">
              {success}
            </div>
          )}

          <input
            type="text"
            placeholder="Nama Lengkap"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="text"
            placeholder="Telepon"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <button
            onClick={handleRegister}
            className="p-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700 transition-colors duration-200 shadow-md hover:shadow-lg"
          >
            Daftar
          </button>
          
          <div className="text-center mt-4">
            <p className="text-gray-600">
              Sudah punya akun?{' '}
              <Link href="/login" className="text-blue-600 hover:text-blue-800 font-semibold">
                Login di sini
              </Link>
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
