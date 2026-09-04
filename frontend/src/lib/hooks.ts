import { useState, useEffect, useCallback } from 'react';
import { apiGet } from './api';
import type { Wallet } from '../types';

export function useWallet() {
  const [wallet, setWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchWallet = useCallback(async () => {
    try {
      setLoading(true);
      const data = await apiGet<Wallet>('/wallet');
      setWallet(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to load wallet');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchWallet();
  }, [fetchWallet]);

  return { wallet, loading, error, refetch: fetchWallet };
}