import { useState, useEffect } from 'react';
import { apiGet, api } from '../lib/api';
import { useWallet } from '../lib/hooks';
import type { StockResponse } from '../types';
import StockTable from '../components/markets/StockTable';
import OrderModal from '../components/orders/OrderModal';
import Spinner from '../components/ui/Spinner';

type Filter = 'All' | 'Gainers' | 'Losers' | 'Most Active';

type PaginatedStockResponse = {
  items: StockResponse[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
};

const simulateDemoPrice = (stock: StockResponse) => {
  const previousClose = stock.close || stock.ltp || 100;
  const phase = Date.now() / 1000;
  const wave = Math.sin(phase / 5 + stock.id) * Math.max(previousClose * 0.008, 0.7);
  const drift = Math.cos(phase / 7 + stock.id * 0.8) * Math.max(previousClose * 0.004, 0.45);
  const next = Math.max(previousClose + wave + drift, 1);

  return Number(next.toFixed(2));
};

export default function Markets() {
  const [stocks, setStocks] = useState<StockResponse[]>([]);
  const [filtered, setFiltered] = useState<StockResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [filter, setFilter] = useState<Filter>('All');
  const [selectedStock, setSelectedStock] = useState<StockResponse | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const { wallet } = useWallet();

  const loadStocks = async (nextPage = 1) => {
    setLoading(true);
    try {
      const data = await apiGet<PaginatedStockResponse>(`/stocks?page=${nextPage}&limit=20`);
      setStocks(data.items);
      setFiltered(data.items);
      setTotalPages(data.total_pages || 1);
      setPage(data.page || 1);
    } catch {
      setStocks([]);
      setFiltered([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStocks(page);
  }, []);

  // Search effect: depends only on `searchQuery` so demo ticks (which mutate `stocks`) don't re-run search
  useEffect(() => {
    const query = searchQuery.trim();
    if (!query) {
      setFiltered(stocks);
      return;
    }

    setSearchLoading(true);
    const controller = new AbortController();
    const timeoutId = window.setTimeout(async () => {
      try {
        const resp = await api.get<unknown>(`/search/stocks?query=${encodeURIComponent(query)}`, { signal: controller.signal });
        const results = (resp.data && (resp.data as any).data) ? (resp.data as any).data as StockResponse[] : [];
        setFiltered(results);
      } catch (err: any) {
        if (err?.name === 'CanceledError' || err?.name === 'AbortError') return;
        setFiltered([]);
      } finally {
        setSearchLoading(false);
      }
    }, 250);

    return () => {
      controller.abort();
      window.clearTimeout(timeoutId);
      setSearchLoading(false);
    };
  }, [searchQuery]);

  // Keep filtered in sync when stocks change only if there's no active search
  useEffect(() => {
    if (!searchQuery.trim()) {
      setFiltered(stocks);
    }
  }, [stocks, searchQuery]);

  useEffect(() => {
    if (stocks.length === 0) return;

    const intervalId = window.setInterval(() => {
      setStocks(prev =>
        prev.map(stock => {
          const nextLtp = simulateDemoPrice(stock);
          const nextHigh = Math.max(stock.high || nextLtp, nextLtp);
          const nextLow = Math.min(stock.low || nextLtp, nextLtp);

          return {
            ...stock,
            ltp: nextLtp,
            high: nextHigh,
            low: nextLow,
          };
        })
      );
    }, 4000);

    return () => window.clearInterval(intervalId);
  }, [stocks.length]);

  const handleTrade = (stock: StockResponse) => {
    setSelectedStock(stock);
    setModalOpen(true);
  };

  const indices = [
    { name: 'S&P 500', value: '5,123.45', change: '+1.24%', positive: true },
    { name: 'NASDAQ', value: '16,234.12', change: '+1.85%', positive: true },
    { name: 'DOW JONES', value: '38,901.04', change: '-0.45%', positive: false },
    { name: 'RUSSELL 2000', value: '2,056.88', change: '+0.92%', positive: true },
  ];

  const filters: Filter[] = ['All', 'Gainers', 'Losers', 'Most Active'];

  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div className="flex items-center gap-3">
          <h2 className="font-headline-md text-headline-md md:font-display-lg md:text-display-lg text-on-surface">Markets</h2>
          <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 text-[10px] font-semibold uppercase tracking-[0.14em] text-emerald-300">
            Demo market mode
          </span>
        </div>
        <div className="relative w-full md:w-96">
          <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant">search</span>
          <input
            className="w-full bg-surface-dark border border-border-muted rounded-full py-3 pl-12 pr-4 text-on-surface focus:border-electric-crimson focus:ring-1 focus:ring-electric-crimson outline-none transition-all font-body-md text-body-md"
            placeholder="Search markets, tickers, or companies..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      {/* Indices */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-gutter">
        {indices.map(idx => (
          <div key={idx.name} className="glass-panel rounded-xl p-4 inner-rim relative overflow-hidden group">
            <h3 className="font-label-caps text-label-caps text-on-surface-variant mb-2">{idx.name}</h3>
            <div className="font-data-tabular text-data-tabular text-xl text-on-surface mb-1">{idx.value}</div>
            <div className={`flex items-center text-xs gap-1 ${idx.positive ? 'text-tertiary' : 'text-electric-crimson'}`}>
              <span className="material-symbols-outlined text-sm">{idx.positive ? 'arrow_upward' : 'arrow_downward'}</span>
              <span className="font-data-tabular text-data-tabular">{idx.change}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Stock Table */}
      <div className="glass-panel rounded-xl flex flex-col inner-rim overflow-hidden">
        <div className="p-4 border-b border-border-muted flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-surface-dark/50">
          <div className="flex gap-2 bg-surface p-1 rounded-lg border border-border-muted overflow-x-auto w-full sm:w-auto">
            {filters.map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-4 py-1.5 rounded-md font-label-caps text-label-caps whitespace-nowrap transition-colors ${
                  filter === f ? 'bg-surface-container-high text-on-surface inner-rim' : 'text-on-surface-variant hover:text-on-surface'
                }`}
              >
                {f}
              </button>
            ))}
          </div>
          <button className="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-border-muted hover:bg-surface-container-low transition-colors text-on-surface-variant text-sm">
            <span className="material-symbols-outlined text-lg">tune</span>
            <span className="font-label-caps text-label-caps">Filters</span>
          </button>
        </div>

        {loading || searchLoading ? (
          <div className="flex flex-col items-center justify-center h-64 gap-3">
            <Spinner size={32} />
            <span className="text-sm text-on-surface-variant">{loading ? 'Loading market data...' : 'Searching stocks...'}</span>
          </div>
        ) : (
          <>
            <StockTable stocks={filtered} onTrade={handleTrade} filter={filter} />
            {searchQuery.trim() === '' && totalPages > 1 && (
              <div className="flex items-center justify-between border-t border-border-muted p-4">
                <button
                  disabled={page === 1}
                  onClick={() => loadStocks(page - 1)}
                  className="px-3 py-2 rounded-lg border border-border-muted text-sm disabled:opacity-40"
                >
                  Previous
                </button>
                <span className="text-sm text-on-surface-variant">Page {page} of {totalPages}</span>
                <button
                  disabled={page >= totalPages}
                  onClick={() => loadStocks(page + 1)}
                  className="px-3 py-2 rounded-lg border border-border-muted text-sm disabled:opacity-40"
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </div>

      <OrderModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        stock={selectedStock}
        walletBalance={wallet?.balance ?? 0}
      />
    </div>
  );
}