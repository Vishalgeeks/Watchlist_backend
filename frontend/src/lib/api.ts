import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import type { ApiResponse, User } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: `${API_URL}/api`,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('auth_token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response: AxiosResponse<ApiResponse<unknown>>) => response,
  (error: AxiosError<ApiResponse<unknown>>) => {
    const url = error.config?.url ?? '';
    if (error.response?.status === 401 && !url.includes('/auth/')) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export async function apiGet<T>(url: string): Promise<T> {
  const { data } = await api.get<ApiResponse<T>>(url);
  if (!data.success) throw new Error(data.message);
  return data.data as T;
}

export async function apiPost<T>(url: string, body: unknown): Promise<T> {
  const { data } = await api.post<ApiResponse<T>>(url, body);
  if (!data.success) throw new Error(data.message);
  return data.data as T;
}

export async function apiPut<T>(url: string, body?: unknown): Promise<T> {
  const { data } = await api.put<ApiResponse<T>>(url, body);
  if (!data.success) throw new Error(data.message);
  return data.data as T;
}

export async function apiDelete<T>(url: string): Promise<T> {
  const { data } = await api.delete<ApiResponse<T>>(url);
  if (!data.success) throw new Error(data.message);
  return data.data as T;
}

export function getAuthToken(): string | null {
  return localStorage.getItem('auth_token');
}

export function setAuthToken(token: string): void {
  localStorage.setItem('auth_token', token);
}

export function clearAuth(): void {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('auth_user');
}

export function getAuthUser(): User | null {
  const stored = localStorage.getItem('auth_user');
  return stored ? JSON.parse(stored) : null;
}

export function setAuthUser(user: User): void {
  localStorage.setItem('auth_user', JSON.stringify(user));
}