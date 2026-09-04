import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import type { User } from '../types';
import { apiPost, getAuthToken, setAuthToken, clearAuth, getAuthUser, setAuthUser } from '../lib/api';
import type { AuthResponse } from '../types';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => getAuthUser());
  const [token, setToken] = useState<string | null>(() => getAuthToken());

  useEffect(() => {
    if (!token) {
      clearAuth();
      setUser(null);
    }
  }, [token]);

  const login = async (email: string, password: string) => {
    const resp = await apiPost<AuthResponse>('/auth/login', { email, password });
    setAuthToken(resp.token);
    setAuthUser(resp.user);
    setToken(resp.token);
    setUser(resp.user);
  };

  const register = async (name: string, email: string, password: string) => {
    const resp = await apiPost<AuthResponse>('/auth/register', { name, email, password });
    setAuthToken(resp.token);
    setAuthUser(resp.user);
    setToken(resp.token);
    setUser(resp.user);
  };

  const logout = () => {
    clearAuth();
    setToken(null);
    setUser(null);
    window.location.href = '/login';
  };

  return (
    <AuthContext.Provider value={{ user, token, isAuthenticated: !!token, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}