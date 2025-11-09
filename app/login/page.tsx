'use client';

import { useState } from 'react';
import { LoginRequest } from '@/proto/carapp_pb';
import { authClient } from '@/lib/grpcClient';
import { useAuth } from '@/src/context/AuthContext';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const { login } = useAuth();
  const router = useRouter();

  const handleLogin = async () => {
    setError('');
    setIsLoading(true);

    if (!email || !password) {
      setError('Email dan Password wajib diisi!');
      setIsLoading(false);
      return;
    }

    const req = new LoginRequest();
    req.setEmail(email);
    req.setPassword(password);

    try {
      const response = await new Promise((resolve, reject) => {
        authClient.login(req, {}, (err, response) => {
          if (err) {
            console.error('gRPC Error:', err);
            reject(err);
          } else {
            resolve(response);
          }
        });
      });
      
      const user = (response as any).getUser();
      const token = (response as any).getToken();

      if (user && token) {
        // Panggil fungsi login dari context kita
        login(token, user);
        
        // Arahkan ke halaman utama
        router.push('/dashboard'); 
      }
    } catch (err: any) {
      console.error('Login error:', err);
      setError(`Login gagal: ${err.message || 'Terjadi kesalahan'}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gray-100">
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-3xl font-bold mb-6 text-center text-gray-800">Login</h1>
        <div className="flex flex-col gap-4 w-full">
          {error && (
            <div className="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}

          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent"
            disabled={isLoading}
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleLogin()}
            className="p-3 border border-gray-300 rounded-lg text-gray-800 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent"
            disabled={isLoading}
          />
          <button
            onClick={handleLogin}
            disabled={isLoading}
            className="p-3 bg-green-600 text-white rounded-lg font-semibold hover:bg-green-700 transition-colors duration-200 shadow-md hover:shadow-lg disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {isLoading ? 'Memproses...' : 'Login'}
          </button>
          
          <p className="text-center text-gray-600 mt-4">
            Belum punya akun?{' '}
            <a href="/" className="text-green-600 hover:text-green-700 font-semibold">
              Daftar di sini
            </a>
          </p>
        </div>
      </div>
    </main>
  );
}
