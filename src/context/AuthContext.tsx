'use client';

import { User } from '@/proto/carapp_pb';
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

// Tipe data yang akan kita simpan di context
interface AuthContextType {
  token: string | null;
  user: User | null;
  isLoading: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
}

// Buat Context-nya
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Buat Provider (pembungkus)
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true); // Untuk cek local storage

  // Saat pertama kali load, cek apakah token sudah ada di local storage
  useEffect(() => {
    const storedToken = localStorage.getItem('authToken');
    const storedUser = localStorage.getItem('authUser');
    if (storedToken && storedUser) {
      setToken(storedToken);
      // Kita perlu parse user dari string JSON kembali ke Tipe 'User'
      const userObj = JSON.parse(storedUser);
      const user = new User();
      user.setId(userObj.id);
      user.setName(userObj.name);
      user.setEmail(userObj.email);
      user.setPhone(userObj.phone);
      user.setRole(userObj.role);
      setUser(user);
    }
    setIsLoading(false);
  }, []);

  const login = (newToken: string, newUser: User) => {
    setToken(newToken);
    setUser(newUser);
    // Simpan di local storage agar tidak hilang saat refresh
    localStorage.setItem('authToken', newToken);
    // Kita simpan sebagai JSON string
    localStorage.setItem('authUser', JSON.stringify(newUser.toObject()));
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('authToken');
    localStorage.removeItem('authUser');
  };

  return (
    <AuthContext.Provider value={{ token, user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Buat custom hook agar mudah digunakan
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
