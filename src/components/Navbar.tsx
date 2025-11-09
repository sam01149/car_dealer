'use client';

import { useAuth } from '@/src/context/AuthContext';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function Navbar() {
  const { user, logout, isLoading } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/login'); // Arahkan ke login setelah logout
  };

  if (isLoading) {
    return (
      <nav className="bg-gray-800 text-white p-4">
        <div className="container mx-auto flex justify-between items-center">
          <div className="text-lg font-bold">Car Dealer</div>
          <div>Loading...</div>
        </div>
      </nav>
    );
  }

  return (
    <nav className="bg-gray-800 text-white p-4">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-lg font-bold">
          🚗 Car Dealer
        </Link>
        <div className="flex gap-4 items-center">
          {user ? (
            // Jika SUDAH login
            <>
              <Link href="/dashboard" className="hover:text-gray-300">
                Dashboard
              </Link>
              <Link href="/mobil/jual" className="hover:text-gray-300">
                Jual Mobil
              </Link>
              <Link href="/notifikasi" className="hover:text-gray-300 flex items-center gap-1">
                <span>🔔</span>
                <span>Notifikasi</span>
              </Link>
              <span className="italic">Halo, {user.getName()}</span>
              <button
                onClick={handleLogout}
                className="bg-red-600 px-3 py-1 rounded hover:bg-red-700"
              >
                Logout
              </button>
            </>
          ) : (
            // Jika BELUM login
            <>
              <Link href="/login" className="hover:text-gray-300">
                Login
              </Link>
              <Link
                href="/register"
                className="bg-blue-600 px-3 py-1 rounded hover:bg-blue-700"
              >
                Register
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
