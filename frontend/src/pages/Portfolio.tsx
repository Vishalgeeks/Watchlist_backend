import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiGet, apiPost } from '../lib/api';
import type { PortfolioSummary, Portfolio } from '../types';
import { formatCurrency, formatPercent } from '../lib/format';
import Spinner from '../components/ui/Spinner';
import OrderModal from '../components/orders/OrderModal';
import { useWallet } from '../lib/hooks';
import type { StockResponse } from '../types';

export default function Portfolio() {
  const [summary, setSummary] = useState<PortfolioSummary | null>(null);
  const [holdings, setHoldings] = useState<Portfolio[]>([]);
  const [showDeposit, setShowDeposit] = useState(false);
  const [amount, setAmount] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedStock, setSelectedStock] = useState<StockResponse | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const { wallet } = useWallet();
  const navigate = useNavigate();

  const handleTransaction = async () => {
    const amt = Number(amount);
    if (!amt || amt <= 0) {
      setErrors({ amount: 'Amount must be greater than 0' });
      return;
    }
    setErrors({});
    setSubmitting(true);
    try {
      await apiGet('/wallet');
      await apiPost('/wallet/deposit', { amount: amt, description: 'Portfolio deposit' });
      setAmount('');
      setShowDeposit(false);
    } catch (err: any) {
      // show simple alert if toast not available
      alert(err?.message || 'Deposit failed');
    } finally {
      setSubmitting(false);
    }
  };

  useEffect(() => {
    const load = async () => {
      try {
        const [summaryData, holdingsData] = await Promise.all([
          apiGet<PortfolioSummary>('/portfolio/summary'),
          apiGet<Portfolio[]>('/portfolio'),
        ]);
        setSummary(summaryData);
        setHoldings(holdingsData);
      } catch (err: any) {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  // Demo tick for holdings to keep portfolio values alive in demo mode
  useEffect(() => {
    if (!holdings || holdings.length === 0) return;

    const simulateDemoPrice = (s: any) => {
      if (!s) return s;
      const previousClose = s.close || s.ltp || 100;
      const phase = Date.now() / 1000;
      const wave = Math.sin(phase / 5 + s.id) * Math.max(previousClose * 0.008, 0.7);
      const drift = Math.cos(phase / 7 + s.id * 0.8) * Math.max(previousClose * 0.004, 0.45);
      const next = Math.max(previousClose + wave + drift, 1);
      return {
        ...s,
        ltp: Number(next.toFixed(2)),
        high: Math.max(s.high || next, Number(next.toFixed(2))),
        low: Math.min(s.low || next, Number(next.toFixed(2))),
      };
    };

    const idInt = setInterval(() => {
      setHoldings(prev => prev.map(h => ({ ...h, stock: simulateDemoPrice(h.stock) })));
    }, 4000);

    return () => clearInterval(idInt);
  }, [holdings]);

  if (loading) {
    return <div className="flex items-center justify-center h-64"><Spinner size={32} /></div>;
  }

  const filtered = holdings.filter(h =>
    !searchQuery ||
    h.stock?.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
    h.stock?.company_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const cards = [
    { label: 'Total Value', value: summary ? formatCurrency(summary.total_value) : '$0.00' },
    { label: 'Total Invested', value: summary ? formatCurrency(summary.total_cost_basis) : '$0.00' },
    { label: 'Available Cash', value: wallet ? formatCurrency(wallet.balance) : '$0.00' },
    { label: 'All-Time P&L', value: summary ? formatCurrency(summary.total_unrealized_pnl) : '$0.00', positive: (summary?.total_unrealized_pnl ?? 0) >= 0 },
  ];

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <header className="flex justify-between items-end">
        <div>
          <h2 className="font-headline-md-mobile md:font-headline-md text-headline-md-mobile md:text-headline-md text-on-surface font-semibold">Portfolio Overview</h2>
          <p className="font-body-md text-body-md text-on-surface-variant mt-1">Real-time performance and holdings</p>
        </div>
        <div className="hidden sm:flex gap-3">
          <button className="bg-surface-light border border-border-muted text-on-surface hover:bg-surface-container-low px-4 py-2 rounded flex items-center gap-2 transition-colors">
            <span className="material-symbols-outlined text-sm">download</span>
            <span className="font-label-caps text-label-caps">Export</span>
          </button>
          <button onClick={() => { setShowDeposit(true); setAmount(''); setErrors({}); }} className="bg-electric-crimson text-white hover:bg-primary-container px-4 py-2 rounded flex items-center gap-2 transition-colors">
            <span className="material-symbols-outlined text-sm">add</span>
            <span className="font-label-caps text-label-caps">Deposit</span>
          </button>
        </div>
      </header>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map(card => (
          <div key={card.label} className="bg-surface-dark border border-border-muted rounded-lg p-5 flex flex-col gap-2 relative overflow-hidden group">
            <span className="font-label-caps text-label-caps text-on-surface-variant">{card.label}</span>
            <div className="font-headline-md text-headline-md tabular-nums tracking-tight text-on-surface">{card.value}</div>
            {card.positive !== undefined && (
              <div className={`flex items-center gap-1 text-xs ${card.positive ? 'text-tertiary' : 'text-error-pure'}`}>
                <span className="font-data-tabular text-data-tabular">{card.positive ? 'Profitable' : 'Loss'}</span>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Performance Chart */}
      <section className="bg-surface-dark border border-border-muted rounded-lg p-6 flex flex-col gap-6">
        <div className="flex justify-between items-center">
          <h3 className="font-label-caps text-label-caps text-on-surface">1M Performance</h3>
          <div className="flex gap-2">
            {['1D', '1W', '1M', '1Y', 'ALL'].map((t, i) => (
              <button key={t} className={`px-3 py-1 rounded font-label-caps text-label-caps transition-colors ${i === 2 ? 'bg-surface-container-high text-electric-crimson border border-border-muted' : 'bg-surface text-on-surface-variant hover:text-on-surface'}`}>
                {t}
              </button>
            ))}
          </div>
        </div>
        <div className="w-full h-64 bg-surface-container relative rounded border border-border-muted overflow-hidden flex items-end">
          <svg className="w-full h-full text-tertiary stroke-current opacity-80" preserveAspectRatio="none" viewBox="0 0 100 100">
            <path d="M0,80 Q10,75 20,60 T40,50 T60,30 T80,20 T100,10 L100,100 L0,100 Z" fill="url(#port-grad)" stroke="none" />
            <path d="M0,80 Q10,75 20,60 T40,50 T60,30 T80,20 T100,10" fill="none" strokeWidth="2" />
            <defs>
              <linearGradient id="port-grad" x1="0%" x2="0%" y1="0%" y2="100%">
                <stop offset="0%" stopColor="currentColor" stopOpacity="0.2" />
                <stop offset="100%" stopColor="currentColor" stopOpacity="0" />
              </linearGradient>
            </defs>
          </svg>
        </div>
      </section>

      {/* Holdings Table */}
      <section className="bg-surface-dark border border-border-muted rounded-lg overflow-x-auto">
        <div className="p-4 border-b border-border-muted flex justify-between items-center bg-surface-container/50">
          <h3 className="font-label-caps text-label-caps text-on-surface">Current Holdings</h3>
          <div className="relative">
            <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm">search</span>
            <input
              className="bg-surface border border-border-muted rounded pl-9 pr-3 py-1.5 text-sm text-on-surface focus:outline-none focus:border-electric-crimson placeholder:text-on-surface-variant/50 w-48"
              placeholder="Search asset..."
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
            />
          </div>
        </div>
        <table className="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-border-muted bg-surface/50 font-label-caps text-label-caps text-on-surface-variant">
              <th className="py-3 px-4 font-normal">Asset</th>
              <th className="py-3 px-4 font-normal text-right">Quantity</th>
              <th className="py-3 px-4 font-normal text-right">Avg Price</th>
              <th className="py-3 px-4 font-normal text-right">Current Price</th>
              <th className="py-3 px-4 font-normal text-right">Invested</th>
              <th className="py-3 px-4 font-normal text-right">Value</th>
              <th className="py-3 px-4 font-normal text-right">P&L</th>
              <th className="py-3 px-4 font-normal text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="font-data-tabular text-data-tabular divide-y divide-border-muted/50">
            {filtered.length === 0 ? (
              <tr><td colSpan={8} className="py-8 text-center text-on-surface-variant">No holdings found.</td></tr>
            ) : (
              filtered.map(h => {
                const currentPrice = h.stock?.ltp ?? 0;
                const invested = h.avg_price * h.quantity;
                const value = h.current_value || currentPrice * h.quantity;
                const pnl = h.unrealized_pnl || (currentPrice - h.avg_price) * h.quantity;
                const pnlPct = invested > 0 ? (pnl / invested) * 100 : 0;
                const positive = pnl >= 0;
                return (
                  <tr key={h.id} className="hover:bg-surface-light transition-colors group">
                    <td className="py-3 px-4 flex items-center gap-3">
                      <div className="w-8 h-8 rounded bg-surface border border-border-muted flex items-center justify-center font-bold text-on-surface">{h.stock?.symbol.slice(0, 4)}</div>
                      <div>
                        <div className="font-semibold text-on-surface">{h.stock?.company_name}</div>
                        <div className="text-xs text-on-surface-variant">{h.stock?.exchange}</div>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-right">{h.quantity}</td>
                    <td className="py-3 px-4 text-right">{formatCurrency(h.avg_price)}</td>
                    <td className="py-3 px-4 text-right">{formatCurrency(currentPrice)}</td>
                    <td className="py-3 px-4 text-right">{formatCurrency(invested)}</td>
                    <td className="py-3 px-4 text-right text-on-surface font-semibold">{formatCurrency(value)}</td>
                    <td className="py-3 px-4 text-right">
                      <div className={positive ? 'text-tertiary' : 'text-error-pure'}>{positive ? '+' : ''}{formatCurrency(pnl)}</div>
                      <div className={`text-xs ${positive ? 'text-tertiary' : 'text-error-pure'}`}>{formatPercent(pnlPct)}</div>
                    </td>
                    <td className="py-3 px-4 text-right">
                      <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button className="px-2 py-1 bg-surface border border-border-muted rounded text-xs hover:text-electric-crimson" onClick={() => navigate(`/stocks/${h.stock_id}`)}>View</button>
                        <button
                          className="px-2 py-1 bg-surface border border-border-muted rounded text-xs hover:text-electric-crimson"
                          onClick={() => {
                            if (h.stock) {
                              setSelectedStock(h.stock);
                              setModalOpen(true);
                            }
                          }}
                        >
                          Trade
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </section>

      <OrderModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        stock={selectedStock}
        walletBalance={wallet?.balance ?? 0}
      />

      {/* Deposit Modal */}
      {showDeposit && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setShowDeposit(false)} />
          <div className="bg-surface-dark border border-border-muted rounded-lg p-6 z-10 w-full max-w-md">
            <h3 className="font-headline-md text-headline-md mb-4">Deposit Funds</h3>
            <div className="space-y-4">
              <input type="number" step="0.01" min="0" value={amount} onChange={e => setAmount(e.target.value)} className="w-full bg-surface border border-border-muted rounded px-3 py-2" placeholder="Amount" />
              {errors.amount && <div className="text-error-pure text-sm">{errors.amount}</div>}
              <div className="flex justify-end gap-2">
                <button className="px-4 py-2 rounded border" onClick={() => setShowDeposit(false)}>Cancel</button>
                <button className="px-4 py-2 rounded bg-electric-crimson text-white" onClick={() => handleTransaction()} disabled={submitting}>{submitting ? '...' : 'Deposit'}</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}