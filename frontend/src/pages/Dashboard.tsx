import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiGet, apiPost } from '../lib/api';
import { useWallet } from '../lib/hooks';
import { useToast } from '../context/ToastContext';
import type { PortfolioSummary, WatchlistItem, Order, Stock } from '../types';
import { formatCurrency, formatPercent } from '../lib/format';
import Modal from '../components/ui/Modal';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import Spinner from '../components/ui/Spinner';
import { depositSchema, formatValidationErrors } from '../lib/validation';

const simulateDemoPrice = (stock: WatchlistItem['stock']) => {
  if (!stock) return stock;

  const previousClose = stock.close || stock.ltp || 100;
  const phase = Date.now() / 1000;
  const wave = Math.sin(phase / 5 + stock.id) * Math.max(previousClose * 0.008, 0.7);
  const drift = Math.cos(phase / 7 + stock.id * 0.8) * Math.max(previousClose * 0.004, 0.45);
  const next = Math.max(previousClose + wave + drift, 1);

  return {
    ...stock,
    ltp: Number(next.toFixed(2)),
    high: Math.max(stock.high || next, Number(next.toFixed(2))),
    low: Math.min(stock.low || next, Number(next.toFixed(2))),
  };
};

const simulateDemoStock = (stock: Stock) => {
  if (!stock) return stock;
  const previousClose = stock.close || stock.ltp || 100;
  const phase = Date.now() / 1000;
  const wave = Math.sin(phase / 5 + stock.id) * Math.max(previousClose * 0.008, 0.7);
  const drift = Math.cos(phase / 7 + stock.id * 0.8) * Math.max(previousClose * 0.004, 0.45);
  const next = Math.max(previousClose + wave + drift, 1);

  return {
    ...stock,
    ltp: Number(next.toFixed(2)),
    high: Math.max(stock.high || next, Number(next.toFixed(2))),
    low: Math.min(stock.low || next, Number(next.toFixed(2))),
  } as Stock;
};

export default function Dashboard() {
  const [summary, setSummary] = useState<PortfolioSummary | null>(null);
  const [watchlistsPreview, setWatchlistsPreview] = useState<Array<{ watchlist: { id: number; name: string }; items: WatchlistItem[] }>>([]);
  const [topStocks, setTopStocks] = useState<Stock[]>([]);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [showDeposit, setShowDeposit] = useState(false);
  const [showWithdraw, setShowWithdraw] = useState(false);
  const [amount, setAmount] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const { wallet, refetch: refetchWallet } = useWallet();
  const { showToast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    const load = async () => {
      try {
        const [summaryData, watchlists, stocksResp, ordersData] = await Promise.all([
          apiGet<PortfolioSummary>('/portfolio/summary'),
          apiGet<{ id: number }[]>('/watchlists'),
          apiGet<any>('/stocks?page=1&limit=3'),
          apiGet<Order[]>('/orders'),
        ]);

        setSummary(summaryData);
        if (watchlists.length > 0) {
          // fetch up to first 3 watchlists and top 4 stocks for each for preview
          const trimmed = watchlists.slice(0, 3);
          const previews = await Promise.all(trimmed.map(async (wl) => {
            try {
              const items = await apiGet<WatchlistItem[]>(`/watchlists/${wl.id}/stocks`);
              return { watchlist: { id: wl.id, name: (wl as any).name || `Watchlist ${wl.id}` }, items: (items || []).slice(0, 4) };
            } catch (e) {
              return { watchlist: { id: wl.id, name: (wl as any).name || `Watchlist ${wl.id}` }, items: [] };
            }
          }));
          setWatchlistsPreview(previews);
        }

        // `/stocks` returns a paging object { items: Stock[] }
        const items: Stock[] = Array.isArray(stocksResp?.items) ? stocksResp.items : [];

        // apply fallback ltp when API returns zero/empty prices so UI shows fluctuations
        const withLtp = items.map(s => ({
          ...s,
          ltp: s.ltp && s.ltp > 0 ? s.ltp : (s.close && s.close > 0 ? s.close : 100 + (s.id % 50)),
        } as Stock));

        setTopStocks(withLtp);

        // ensure recent orders include stock and fallback ltp
        const orders = (ordersData || []).slice(0, 3).map(o => {
          if (!o.stock) return o;
          const stock = { ...o.stock } as Stock;
          stock.ltp = stock.ltp && stock.ltp > 0 ? stock.ltp : (stock.close && stock.close > 0 ? stock.close : 100 + (stock.id % 50));
          return { ...o, stock } as Order;
        });

        setRecentOrders(orders);
      } catch (err: any) {
        setTimeout(() => showToast(err.message || 'Failed to load dashboard', 'error'));
      } finally {
        setLoading(false);
      }
    };
    load();
    // refresh when watchlists change elsewhere
    const onWatchlistChange = async () => {
      try {
        const lists = await apiGet<any[]>('/watchlists');
        if (lists.length > 0) {
          const trimmed = lists.slice(0, 3);
          const previews = await Promise.all(trimmed.map(async (wl) => {
            try {
              const items = await apiGet<WatchlistItem[]>(`/watchlists/${wl.id}/stocks`);
              return { watchlist: { id: wl.id, name: wl.name || `Watchlist ${wl.id}` }, items: (items || []).slice(0, 4) };
            } catch (e) {
              return { watchlist: { id: wl.id, name: wl.name || `Watchlist ${wl.id}` }, items: [] };
            }
          }));
          setWatchlistsPreview(previews);
        }
      } catch (e) {
        console.error('onWatchlistChange error', e);
      }
    };

    window.addEventListener('watchlists:changed', onWatchlistChange);

    return () => window.removeEventListener('watchlists:changed', onWatchlistChange);
  }, [showToast]);

  useEffect(() => {
    if (!watchlistsPreview || watchlistsPreview.length === 0) return;
    const intervalId = window.setInterval(() => {
      try {
        setWatchlistsPreview(prev => {
          if (!Array.isArray(prev)) return prev;
          return prev.map(w => {
            try {
              const newItems = (w.items || []).map(item => {
                try {
                  return { ...item, stock: simulateDemoPrice(item.stock) };
                } catch (e) {
                  console.error('simulateDemoPrice item error', e);
                  return item;
                }
              });
              return { ...w, items: newItems };
            } catch (e) {
              console.error('watchlist group tick error', e);
              return w;
            }
          });
        });
      } catch (e) {
        console.error('watchlistsPreview tick error', e);
      }
    }, 4000);

    return () => window.clearInterval(intervalId);
  }, [watchlistsPreview?.length]);

  // demo ticks for top stocks
  useEffect(() => {
    if (!topStocks || topStocks.length === 0) return;
    const id = window.setInterval(() => {
      try {
        setTopStocks(prev => {
          if (!Array.isArray(prev)) {
            console.warn('topStocks prev is not array', prev);
            return prev as any;
          }
          return prev.map(s => {
            try {
              return simulateDemoStock(s);
            } catch (e) {
              console.error('simulateDemoStock error', e);
              return s;
            }
          });
        });
      } catch (e) {
        console.error('topStocks tick error', e);
      }
    }, 4000);
    return () => window.clearInterval(id);
  }, [topStocks?.length]);

  // recentOrders: keep static (no demo ticks)

  // global error handlers to surface runtime issues during development
  useEffect(() => {
    const onError = (ev: ErrorEvent) => {
      console.error('Global error:', ev.error || ev.message);
      setTimeout(() => showToast(`Dashboard runtime error: ${ev.message}`, 'error'));
    };
    const onRejection = (ev: PromiseRejectionEvent) => {
      console.error('Unhandled rejection:', ev.reason);
      setTimeout(() => showToast(`Unhandled promise rejection: ${String(ev.reason)}`, 'error'));
    };
    window.addEventListener('error', onError);
    window.addEventListener('unhandledrejection', onRejection);
    return () => {
      window.removeEventListener('error', onError);
      window.removeEventListener('unhandledrejection', onRejection);
    };
  }, [showToast]);

  const handleTransaction = async (type: 'deposit' | 'withdraw') => {
    const schema = type === 'deposit' ? depositSchema : depositSchema;
    const result = schema.safeParse({ amount: Number(amount) });
    if (!result.success) {
      setErrors(formatValidationErrors(result.error));
      return;
    }
    setErrors({});
    setSubmitting(true);
    try {
      await apiPost(`/wallet/${type}`, { amount: Number(amount) });
      showToast(`${type === 'deposit' ? 'Deposit' : 'Withdrawal'} successful`, 'success');
      setAmount('');
      setShowDeposit(false);
      setShowWithdraw(false);
      refetchWallet();
    } catch (err: any) {
      showToast(err.message || `${type} failed`, 'error');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size={32} />
      </div>
    );
  }

  const cards = [
    { label: 'Total Value', value: summary ? formatCurrency(summary.total_value) : '$0.00', change: summary ? `+${formatPercent(summary.total_pnl_percentage)}` : '', positive: true },
    { label: "Today's P&L", value: summary ? formatCurrency(summary.total_unrealized_pnl) : '$0.00', change: 'All-time', positive: true },
    { label: 'Available Balance', value: wallet ? formatCurrency(wallet.balance) : '$0.00', change: 'Ready to invest', positive: false },
    { label: 'Holdings', value: summary ? summary.holdings_count.toString() : '0', change: 'Assets', positive: false },
  ];

  return (
    <div className="space-y-8">
      {/* Greeting */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <div className="flex items-center space-x-2 mb-1">
            <span className="w-2 h-2 rounded-full bg-tertiary animate-pulse" />
            <span className="font-label-caps text-label-caps text-tertiary">Market Open</span>
          </div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Good morning, Trader</h2>
        </div>
        <div className="flex space-x-3">
          <Button variant="secondary" onClick={() => { setShowWithdraw(true); setAmount(''); setErrors({}); }}>
            Withdraw
          </Button>
          <Button onClick={() => { setShowDeposit(true); setAmount(''); setErrors({}); }}>
            Add Funds
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-gutter">
        {cards.map(card => (
          <div key={card.label} className="bg-surface-dark border border-border-muted rounded-xl p-6 relative overflow-hidden group hover:border-electric-crimson/50 transition-colors">
            <h3 className="font-label-caps text-label-caps text-on-surface-variant mb-2">{card.label}</h3>
            <div className="font-data-tabular text-[28px] leading-tight font-medium text-on-surface mb-1">{card.value}</div>
            <div className={`flex items-center font-data-tabular text-sm ${card.positive ? 'text-tertiary' : 'text-on-surface-variant'}`}>
              {card.change}
            </div>
            <div className="absolute -bottom-10 -right-10 w-24 h-24 bg-tertiary/10 blur-2xl rounded-full" />
          </div>
        ))}
      </div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-gutter">
        {/* Left: Market Overview + Recent Orders */}
        <div className="lg:col-span-2 space-y-gutter">
          <section>
            <div className="flex items-center justify-between mb-4">
              <h3 className="font-headline-md text-headline-md text-on-surface">Market Overview</h3>
              <button onClick={() => navigate('/markets')} className="text-on-surface-variant hover:text-electric-crimson text-sm flex items-center transition-colors">
                View All <span className="material-symbols-outlined text-[18px] ml-1">chevron_right</span>
              </button>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              {(topStocks && topStocks.length > 0 ? topStocks : []).map(stock => {
                const change = stock && stock.close > 0 ? `${formatPercent(((stock.ltp || 0) - (stock.close || 0)) / (stock.close || 1) * 100)}` : '0%';
                const positive = stock ? (stock.ltp - (stock.close || 0)) >= 0 : true;
                return (
                  <div key={stock.id} className="bg-surface-dark border border-border-muted rounded-lg p-4 relative overflow-hidden">
                    <div className="flex justify-between items-start mb-2">
                      <span className="font-label-caps text-label-caps text-on-surface">{stock.display_name || stock.symbol}</span>
                      <span className={`font-data-tabular text-sm ${positive ? 'text-tertiary' : 'text-electric-crimson'}`}>{change}</span>
                    </div>
                    <div className="font-data-tabular text-lg text-on-surface">{stock ? formatCurrency(stock.ltp || 0) : '-'}</div>
                    <div className="mt-4 h-8 w-full border-b-2 relative" style={{ borderColor: positive ? 'rgba(103,221,151,0.5)' : 'rgba(255,0,94,0.5)' }}>
                      <div className="absolute bottom-0 left-0 w-full h-full" style={{ background: `linear-gradient(to top, ${positive ? 'rgba(103,221,151,0.1)' : 'rgba(255,0,94,0.1)'}, transparent)` }} />
                    </div>
                  </div>
                );
              })}
            </div>
          </section>

          <section className="bg-surface-dark border border-border-muted rounded-xl p-6">
            <div className="flex items-center justify-between mb-6">
              <h3 className="font-headline-md text-headline-md text-on-surface">Recent Orders</h3>
              <button onClick={() => navigate('/orders')} className="material-symbols-outlined text-on-surface-variant hover:text-on-surface transition-colors">more_horiz</button>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-border-muted">
                    <th className="pb-3 font-label-caps text-label-caps text-on-surface-variant font-normal">Symbol</th>
                    <th className="pb-3 font-label-caps text-label-caps text-on-surface-variant font-normal">Type</th>
                    <th className="pb-3 font-label-caps text-label-caps text-on-surface-variant font-normal text-right">Price</th>
                    <th className="pb-3 font-label-caps text-label-caps text-on-surface-variant font-normal text-right">Shares</th>
                    <th className="pb-3 font-label-caps text-label-caps text-on-surface-variant font-normal text-right">Status</th>
                  </tr>
                </thead>
                <tbody className="font-data-tabular text-sm">
                  {(recentOrders && recentOrders.length > 0 ? recentOrders : []).map(o => {
                    const price = o.price ?? o.stock?.ltp ?? 0;
                    const isBuy = o.side === 'BUY';
                    return (
                      <tr key={o.id} className="border-b border-border-muted/50 hover:bg-surface-light transition-colors">
                        <td className="py-4">
                          <div className="flex flex-col items-start max-w-[220px]">
                            <span className="font-label-caps text-label-caps text-on-surface mb-1">{o.stock?.symbol || '—'}</span>
                            <span className="text-xs text-on-surface-variant truncate">{o.stock?.company_name || o.stock?.display_name || ''}</span>
                          </div>
                        </td>
                        <td className="py-4"><span className={`px-2 py-1 rounded text-xs ${isBuy ? 'bg-tertiary/10 text-tertiary' : 'bg-error-pure/10 text-error-pure'}`}>{isBuy ? 'Buy' : 'Sell'}</span></td>
                        <td className="py-4 text-right text-on-surface">{formatCurrency(price)}</td>
                        <td className="py-4 text-right text-on-surface-variant">{o.quantity}</td>
                        <td className="py-4 text-right"><span className={o.status === 'FILLED' ? 'text-tertiary' : 'text-on-surface-variant'}>{o.status === 'FILLED' ? 'Filled' : o.status === 'PENDING' ? 'Pending' : o.status}</span></td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </section>
        </div>

        {/* Right: Watchlist Preview */}
        <div className="space-y-gutter">
          <section className="bg-surface-dark border border-border-muted rounded-xl p-6 h-full flex flex-col">
            <div className="flex items-center justify-between mb-6">
              <h3 className="font-headline-md text-headline-md text-on-surface">Watchlist</h3>
              <button onClick={() => navigate('/watchlist')} className="material-symbols-outlined text-on-surface-variant hover:text-electric-crimson transition-colors">add</button>
            </div>
            <div className="flex-1 overflow-y-auto pr-2 space-y-4">
              {(!watchlistsPreview || watchlistsPreview.length === 0) ? (
                <p className="text-on-surface-variant text-sm">No watchlist items yet.</p>
              ) : (
                watchlistsPreview.map(wp => (
                  <div key={wp.watchlist.id} className="space-y-2">
                    <div className="text-sm font-medium text-on-surface mb-1">{wp.watchlist.name}</div>
                    {(wp.items || []).map(item => {
                      const stock = item.stock;
                      const change = stock ? (stock.close > 0 ? ((stock.ltp - stock.close) / stock.close) * 100 : 0) : 0;
                      const positive = change >= 0;
                      return (
                        <div key={item.id} className="flex items-center justify-between p-3 rounded-lg hover:bg-surface-container-low transition-colors cursor-pointer group border border-transparent hover:border-border-muted" onClick={() => stock && navigate(`/stocks/${stock.id}`)}>
                          <div className="flex flex-col">
                            <span className="font-label-caps text-label-caps text-on-surface mb-1">{stock?.symbol}</span>
                            <span className="text-xs text-on-surface-variant">{stock?.company_name}</span>
                          </div>
                          <div className="flex flex-col items-end">
                            <span className="font-data-tabular text-sm text-on-surface">{stock ? formatCurrency(stock.ltp) : '-'}</span>
                            <span className={`font-data-tabular text-xs ${positive ? 'text-tertiary' : 'text-electric-crimson'}`}>{formatPercent(change)}</span>
                          </div>
                        </div>
                      );
                    })}
                    <div className="h-px bg-border-muted/30" />
                  </div>
                ))
              )}
            </div>
            <button onClick={() => navigate('/watchlist')} className="w-full mt-4 py-2 border border-border-muted rounded-lg text-on-surface-variant font-label-caps text-label-caps hover:text-on-surface hover:bg-surface-container-low transition-all">
              View All Watchlists
            </button>
          </section>
        </div>
      </div>

      {/* Deposit Modal */}
      <Modal isOpen={showDeposit} onClose={() => setShowDeposit(false)} title="Add Funds">
        <div className="space-y-4">
          <Input
            id="depositAmount"
            label="Amount"
            type="number"
            min="0"
            step="0.01"
            placeholder="0.00"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            error={errors.amount}
            icon="attach_money"
          />
          <Button onClick={() => handleTransaction('deposit')} disabled={submitting} className="w-full">
            {submitting ? <Spinner size={20} /> : 'Deposit'}
          </Button>
        </div>
      </Modal>

      {/* Withdraw Modal */}
      <Modal isOpen={showWithdraw} onClose={() => setShowWithdraw(false)} title="Withdraw Funds">
        <div className="space-y-4">
          <Input
            id="withdrawAmount"
            label="Amount"
            type="number"
            min="0"
            step="0.01"
            placeholder="0.00"
            value={amount}
            onChange={e => setAmount(e.target.value)}
            error={errors.amount}
            icon="money_off"
          />
          <Button onClick={() => handleTransaction('withdraw')} disabled={submitting} className="w-full">
            {submitting ? <Spinner size={20} /> : 'Withdraw'}
          </Button>
        </div>
      </Modal>
    </div>
  );
}