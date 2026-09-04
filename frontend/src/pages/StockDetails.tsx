import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { apiGet } from '../lib/api';
import { useWallet } from '../lib/hooks';
import type { Stock } from '../types';
import { formatCurrency, formatPercent } from '../lib/format';
import OrderModal from '../components/orders/OrderModal';
import Spinner from '../components/ui/Spinner';
import Button from '../components/ui/Button';

export default function StockDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [stock, setStock] = useState<Stock | null>(null);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [orderSide, setOrderSide] = useState<'BUY' | 'SELL'>('BUY');
  const { wallet } = useWallet();

  useEffect(() => {
    const load = async () => {
      try {
        const data = await apiGet<Stock>(`/stocks/${id}`);
        setStock(data);
      } catch (err: any) {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [id]);

  // Demo price simulation for standalone stock details (keeps UI lively in demo mode)
  useEffect(() => {
    if (!stock) return;

    const simulateDemoPrice = (s: Stock) => {
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
      } as Stock;
    };

    const idInt = setInterval(() => {
      setStock(prev => (prev ? simulateDemoPrice(prev) : prev));
    }, 4000);

    return () => clearInterval(idInt);
  }, [stock]);

  if (loading) {
    return <div className="flex items-center justify-center h-64"><Spinner size={32} /></div>;
  }

  if (!stock) {
    return (
      <div className="text-center py-16">
        <p className="text-on-surface-variant mb-4">Stock not found.</p>
        <Button onClick={() => navigate('/markets')}>Back to Markets</Button>
      </div>
    );
  }

  const change = stock.ltp - stock.close;
  const changePct = stock.close > 0 ? (change / stock.close) * 100 : 0;
  const positive = change >= 0;

  const stats = [
    { label: 'Open', value: stock.open.toFixed(2) },
    { label: 'Prev Close', value: stock.close.toFixed(2) },
    { label: 'High', value: stock.high.toFixed(2) },
    { label: 'Low', value: stock.low.toFixed(2) },
    { label: 'Volume', value: stock.vol.toLocaleString() },
    { label: 'OI', value: stock.oi.toLocaleString() },
  ];

  const openModal = (side: 'BUY' | 'SELL') => {
    setOrderSide(side);
    setModalOpen(true);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b border-border-muted">
        <div className="flex items-center gap-6">
          <div className="w-16 h-16 rounded-xl bg-surface-container flex items-center justify-center border border-border-muted shadow-lg shrink-0">
            <span className="font-headline-md text-headline-md font-black text-on-surface">{stock.symbol.slice(0, 2)}</span>
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h2 className="font-display-lg text-display-lg text-on-surface m-0 leading-none">{stock.symbol}</h2>
              <span className="px-2 py-1 bg-surface-container-low text-on-surface-variant font-label-caps text-label-caps rounded border border-border-muted uppercase">{stock.exchange}</span>
            </div>
            <p className="font-body-md text-body-md text-on-surface-variant mt-1">{stock.company_name}</p>
          </div>
        </div>
        <div className="flex flex-col md:items-end gap-1">
          <div className="flex items-end gap-3">
            <span className="font-display-lg text-display-lg text-on-surface font-tabular-nums leading-none tracking-tight">{formatCurrency(stock.ltp)}</span>
            <span className={`font-data-tabular text-data-tabular flex items-center bg-tertiary/10 px-2 py-1 rounded border border-tertiary/20 mb-1 ${positive ? 'text-tertiary' : 'text-electric-crimson'}`}>
              <span className="material-symbols-outlined text-sm mr-1">{positive ? 'arrow_upward' : 'arrow_downward'}</span>
              {formatPercent(changePct)}
            </span>
          </div>
          <p className="font-label-caps text-label-caps text-on-surface-variant">Market Open</p>
        </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Left: Chart & Stats */}
        <div className="lg:col-span-8 xl:col-span-9 flex flex-col gap-6">
          {/* Chart */}
          <div className="bg-surface-dark rounded-xl border border-border-muted flex flex-col overflow-hidden">
            <div className="flex items-center justify-between p-4 border-b border-border-muted bg-surface-container-lowest/50">
              <div className="flex gap-2">
                {['1D', '1W', '1M', '3M', '1Y', 'ALL'].map((t, i) => (
                  <button key={t} className={`px-3 py-1 rounded font-label-caps text-label-caps transition-colors ${i === 2 ? 'bg-primary/10 text-electric-crimson border border-electric-crimson/30' : 'text-on-surface-variant hover:bg-surface-container'}`}>
                    {t}
                  </button>
                ))}
              </div>
            </div>
            <div className="h-[400px] w-full relative p-4 chart-grid flex flex-col justify-between">
              <svg className="w-full h-full absolute inset-0 z-0" preserveAspectRatio="none" viewBox="0 0 1000 400">
                  {(() => {
                    // build a simple dynamic sparkline path based on current ltp
                    const pts = [] as number[];
                    const base = stock.ltp || stock.close || 100;
                    for (let i = 0; i <= 10; i++) {
                      const t = Date.now() / 1000 + i;
                      const y = 350 - (Math.sin(t / 2 + base / 100 + i) * Math.min(base * 0.6, 120));
                      pts.push(Math.max(30, Math.min(370, y)));
                    }
                    const pathD = pts.map((y, i) => `${i === 0 ? 'M' : 'L'}${(i / 10) * 1000} ${y}`).join(' ');
                    const fillPath = `${pathD} L 1000 400 L 0 400 Z`;
                    return (
                      <>
                        <path d={pathD} fill="none" stroke={positive ? '#67dd97' : '#FF005E'} strokeWidth="2" vectorEffect="non-scaling-stroke" />
                        <path d={fillPath} fill={positive ? 'url(#grad-g)' : 'url(#grad-r)'} opacity="0.18" />
                      </>
                    );
                  })()}
                <defs>
                  <linearGradient id="grad-g" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="#67dd97" stopOpacity="1" />
                    <stop offset="100%" stopColor="#67dd97" stopOpacity="0" />
                  </linearGradient>
                  <linearGradient id="grad-r" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="#FF005E" stopOpacity="1" />
                    <stop offset="100%" stopColor="#FF005E" stopOpacity="0" />
                  </linearGradient>
                </defs>
              </svg>
            </div>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-4">
            {stats.map(s => (
              <div key={s.label} className="bg-surface-container rounded-lg border border-border-muted p-4 relative overflow-hidden group hover:border-surface-bright transition-colors">
                <h4 className="font-label-caps text-label-caps text-on-surface-variant mb-1">{s.label}</h4>
                <p className="font-data-tabular text-data-tabular text-on-surface text-lg">{s.value}</p>
              </div>
            ))}
          </div>

          {/* About */}
          <div className="bg-surface-dark rounded-xl border border-border-muted p-6">
            <h3 className="font-headline-md text-headline-md text-on-surface mb-4">About {stock.company_name}</h3>
            <p className="font-body-md text-body-md text-on-surface-variant leading-relaxed">
              {stock.description || `${stock.company_name} (${stock.symbol}) trades on the ${stock.exchange} exchange under instrument type ${stock.instrument_type}.`}
            </p>
            <div className="mt-4 flex flex-wrap gap-4">
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-on-surface-variant text-sm">tag</span>
                <span className="font-label-caps text-label-caps text-on-surface-variant">ISIN: {stock.isin}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-on-surface-variant text-sm">straighten</span>
                <span className="font-label-caps text-label-caps text-on-surface-variant">Lot Size: {stock.lot_size}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="material-symbols-outlined text-on-surface-variant text-sm">tune</span>
                <span className="font-label-caps text-label-caps text-on-surface-variant">Tick: {stock.tick_size}</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right: Trade Panel */}
        <div className="lg:col-span-4 xl:col-span-3">
          <div className="sticky top-gutter bg-surface-container-high rounded-xl border border-border-muted p-6 flex flex-col gap-6 shadow-2xl">
            <div className="absolute top-0 right-0 w-32 h-32 bg-electric-crimson/5 rounded-full blur-3xl pointer-events-none" />
            <h3 className="font-label-caps text-label-caps text-on-surface-variant mb-4 uppercase tracking-widest border-b border-border-muted pb-2">Trade {stock.symbol}</h3>
            <div className="flex items-center justify-between">
              <span className="font-body-md text-body-md text-on-surface">Available Cash</span>
              <span className="font-data-tabular text-data-tabular text-on-surface">{wallet ? formatCurrency(wallet.balance) : '$0.00'}</span>
            </div>
            <div className="flex gap-3">
              <Button onClick={() => openModal('BUY')} className="flex-1 bg-tertiary/20 text-tertiary border border-tertiary/40 hover:bg-tertiary/30">
                <span className="material-symbols-outlined text-[18px]">add_shopping_cart</span> Buy
              </Button>
              <Button onClick={() => openModal('SELL')} variant="secondary" className="flex-1 border-error-pure/40 text-error-pure hover:bg-error-pure/10">
                <span className="material-symbols-outlined text-[18px]">sell</span> Sell
              </Button>
            </div>
          </div>
        </div>
      </div>

      <OrderModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        stock={stock}
        walletBalance={wallet?.balance ?? 0}
        defaultSide={orderSide}
      />
    </div>
  );
}