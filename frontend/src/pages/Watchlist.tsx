import { useState, useEffect, useRef } from 'react';
import { apiGet, apiPost, apiDelete } from '../lib/api';
import { useToast } from '../context/ToastContext';
import { useWallet } from '../lib/hooks';
import type { Watchlist, WatchlistItem, StockResponse } from '../types';
import { formatCurrency, formatPercent } from '../lib/format';
import Spinner from '../components/ui/Spinner';
import Modal from '../components/ui/Modal';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import OrderModal from '../components/orders/OrderModal';
import { createWatchlistSchema, addStockSchema, formatValidationErrors } from '../lib/validation';

const simulateDemoPrice = (stock: StockResponse) => {
  const previousClose = stock.close || stock.ltp || 100;
  const phase = Date.now() / 1000;
  const wave = Math.sin(phase / 5 + stock.id) * Math.max(previousClose * 0.008, 0.7);
  const drift = Math.cos(phase / 7 + stock.id * 0.8) * Math.max(previousClose * 0.004, 0.45);
  const next = Math.max(previousClose + wave + drift, 1);

  return Number(next.toFixed(2));
};

const applyDemoTick = (stock: StockResponse | undefined) => {
  if (!stock) return stock;

  const nextLtp = simulateDemoPrice(stock);
  const nextHigh = Math.max(stock.high || nextLtp, nextLtp);
  const nextLow = Math.min(stock.low || nextLtp, nextLtp);

  return {
    ...stock,
    ltp: nextLtp,
    high: nextHigh,
    low: nextLow,
  };
};

export default function Watchlist() {
  const [watchlists, setWatchlists] = useState<Watchlist[]>([]);
  const [activeList, setActiveList] = useState<Watchlist | null>(null);
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [searchResults, setSearchResults] = useState<StockResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddList, setShowAddList] = useState(false);
  const [showAddStock, setShowAddStock] = useState(false);
  const [listName, setListName] = useState('');
  const [selectedStock, setSelectedStock] = useState<StockResponse | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [tradeStock, setTradeStock] = useState<StockResponse | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const { wallet } = useWallet();
  const { showToast } = useToast();

  const searchDebounce = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        const lists = await apiGet<Watchlist[]>('/watchlists');
        setWatchlists(lists);
        if (lists.length > 0) {
          setActiveList(lists[0]);
          const itemsData = await apiGet<WatchlistItem[]>(`/watchlists/${lists[0].id}/stocks`);
          setItems(itemsData);
        }
      } catch (err: any) {
        showToast(err.message || 'Failed to load watchlist', 'error');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [showToast]);

  useEffect(() => {
    if (items.length === 0) return;

    const intervalId = window.setInterval(() => {
      setItems(prev => prev.map(item => ({
        ...item,
        stock: applyDemoTick(item.stock),
      })));
    }, 4000);

    return () => window.clearInterval(intervalId);
  }, [items.length]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setSelectedStock(null);
    setErrors({});

    if (searchDebounce.current) {
      clearTimeout(searchDebounce.current);
    }

    if (query.trim().length < 2) {
      setSearchResults([]);
      return;
    }

    searchDebounce.current = setTimeout(async () => {
      try {
        const results = await apiGet<StockResponse[]>(`/search/stocks?query=${encodeURIComponent(query.trim())}`);
        setSearchResults(results);
      } catch {
        setSearchResults([]);
      }
    }, 300);
  };

  const handleCreateList = async () => {
    const result = createWatchlistSchema.safeParse({ name: listName });
    if (!result.success) {
      setErrors(formatValidationErrors(result.error));
      return;
    }
    setErrors({});
    setSubmitting(true);
    try {
      const newList = await apiPost<Watchlist>('/watchlists', { name: listName });
      setWatchlists(prev => [...prev, newList]);
      setActiveList(newList);
      setItems([]);
      setListName('');
      setShowAddList(false);
      showToast('Watchlist created', 'success');
    } catch (err: any) {
      showToast(err.message || 'Failed to create watchlist', 'error');
    } finally {
      setSubmitting(false);
    }
  };

  const handleAddStock = async () => {
    if (!activeList) return;
    const result = addStockSchema.safeParse({ stock_id: selectedStock ? selectedStock.id : 0 });
    if (!result.success) {
      setErrors({ stock_id: 'Please select a stock' });
      return;
    }
    setErrors({});
    setSubmitting(true);
    try {
      await apiPost(`/watchlists/${activeList.id}/stocks`, { stock_id: selectedStock!.id });
      const itemsData = await apiGet<WatchlistItem[]>(`/watchlists/${activeList.id}/stocks`);
      setItems(itemsData);
      setSelectedStock(null);
      setShowAddStock(false);
      // notify other pages that watchlists changed
      window.dispatchEvent(new CustomEvent('watchlists:changed', { detail: { listId: activeList.id } }));
      showToast('Stock added to watchlist', 'success');
    } catch (err: any) {
      showToast(err.message || 'Failed to add stock', 'error');
    } finally {
      setSubmitting(false);
    }
  };

  const handleRemoveStock = async (stockId: number) => {
    if (!activeList) return;
    try {
      await apiDelete(`/watchlists/${activeList.id}/stocks/${stockId}`);
      setItems(prev => prev.filter(i => i.stock_id !== stockId));
      showToast('Stock removed', 'success');
    } catch (err: any) {
      showToast(err.message || 'Failed to remove stock', 'error');
    }
  };

  const openTrade = (stock: StockResponse) => {
    setTradeStock(stock);
    setModalOpen(true);
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64"><Spinner size={32} /></div>;
  }

  return (
    <div className="space-y-6">
      <header className="mb-8 flex items-center justify-between">
        <div>
          <h2 className="font-headline-md text-headline-md font-bold text-on-surface">My Watchlist</h2>
          {watchlists.length > 1 && (
            <div className="flex gap-2 mt-3">
              {watchlists.map(list => (
                <button
                  key={list.id}
                  onClick={() => {
                    setActiveList(list);
                    apiGet<WatchlistItem[]>(`/watchlists/${list.id}/stocks`).then(setItems).catch(() => setItems([]));
                  }}
                  className={`px-3 py-1 rounded-lg text-xs font-label-caps transition-colors ${
                    activeList?.id === list.id ? 'bg-electric-crimson/20 text-electric-crimson border border-electric-crimson/40' : 'bg-surface-container text-on-surface-variant border border-border-muted'
                  }`}
                >
                  {list.name}
                </button>
              ))}
            </div>
          )}
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => { setShowAddList(true); setListName(''); setErrors({}); }}>
            <span className="material-symbols-outlined text-[18px]">playlist_add</span> New List
          </Button>
          <Button onClick={() => { setShowAddStock(true); setSelectedStock(null); setSearchQuery(''); setSearchResults([]); setErrors({}); }} disabled={!activeList}>
            <span className="material-symbols-outlined text-[18px]">add</span> Add Asset
          </Button>
        </div>
      </header>

      {!activeList ? (
        <div className="flex flex-col items-center justify-center py-24 text-center">
          <div className="w-16 h-16 rounded-full bg-surface-container flex items-center justify-center mb-4 border border-border-muted">
            <span className="material-symbols-outlined text-on-surface-variant text-3xl">visibility</span>
          </div>
          <h3 className="font-headline-md text-headline-md text-on-surface mb-2">No watchlist yet</h3>
          <p className="text-on-surface-variant mb-6 max-w-sm">Create a watchlist to track your favorite stocks.</p>
          <Button onClick={() => { setShowAddList(true); setListName(''); setErrors({}); }}>
            <span className="material-symbols-outlined text-[18px]">playlist_add</span> Create Watchlist
          </Button>
        </div>
      ) : items.length === 0 ? (
        <div className="bg-surface-dark border border-border-muted rounded-lg p-12 text-center">
          <p className="text-on-surface-variant mb-4">No stocks in this watchlist yet.</p>
          <Button onClick={() => { setShowAddStock(true); setSelectedStock(null); setSearchQuery(''); setSearchResults([]); setErrors({}); }}>
            <span className="material-symbols-outlined text-[18px]">add</span> Add Asset
          </Button>
        </div>
      ) : (
        <div className="bg-surface-dark border border-border-muted rounded-lg overflow-hidden shadow-2xl relative">
          <div className="absolute inset-0 bg-white/[0.02] pointer-events-none" />
          <table className="w-full text-left border-collapse relative z-10">
            <thead>
              <tr className="border-b border-border-muted bg-surface/50">
                <th className="py-3 px-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Asset</th>
                <th className="py-3 px-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Price</th>
                <th className="py-3 px-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">24h Change</th>
                <th className="py-3 px-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center hidden sm:table-cell">Trend</th>
                <th className="py-3 px-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Action</th>
              </tr>
            </thead>
            <tbody>
              {items.map(item => {
                const stock = item.stock;
                if (!stock) return null;
                const change = stock.ltp - stock.close;
                const changePct = stock.close > 0 ? (change / stock.close) * 100 : 0;
                const positive = change >= 0;
                const path = positive
                  ? 'M0,25 L20,20 L40,28 L60,15 L80,5 L100,2'
                  : 'M0,5 L20,15 L40,10 L60,25 L80,20 L100,28';
                return (
                  <tr key={item.id} className="border-b border-border-muted hover:bg-surface-container-low transition-colors group">
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-surface-container border border-border-muted flex items-center justify-center font-label-caps text-label-caps">{stock.symbol.slice(0, 4)}</div>
                        <div>
                          <div className="font-body-md font-bold text-on-surface">{stock.company_name}</div>
                          <div className="font-label-caps text-label-caps text-on-surface-variant">{stock.symbol}</div>
                        </div>
                      </div>
                    </td>
                    <td className="py-4 px-4 text-right">
                      <div className="font-data-tabular text-data-tabular text-on-surface">{formatCurrency(stock.ltp)}</div>
                    </td>
                    <td className="py-4 px-4 text-right">
                      <div className={`font-data-tabular text-data-tabular ${positive ? 'text-tertiary' : 'text-electric-crimson'}`}>{formatPercent(changePct)}</div>
                    </td>
                    <td className="py-4 px-4 text-center hidden sm:table-cell w-32">
                      <svg className={`w-full h-8 sparkline-svg ${positive ? 'sparkline-positive' : 'sparkline-negative'}`} preserveAspectRatio="none" viewBox="0 0 100 30">
                        <path d={path} />
                      </svg>
                    </td>
                    <td className="py-4 px-4 text-right">
                      <div className="flex justify-end gap-2">
                        <button
                          className="p-2 rounded bg-primary/10 text-electric-crimson hover:bg-primary/20 transition-colors border border-primary/20 font-label-caps text-label-caps"
                          onClick={() => openTrade(stock)}
                        >
                          Trade
                        </button>
                        <button
                          className="p-2 rounded text-on-surface-variant hover:text-error-pure hover:bg-error-pure/10 transition-colors opacity-0 group-hover:opacity-100"
                          onClick={() => handleRemoveStock(stock.id)}
                        >
                          <span className="material-symbols-outlined text-[18px]">close</span>
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* Create List Modal */}
      <Modal isOpen={showAddList} onClose={() => setShowAddList(false)} title="Create Watchlist">
        <div className="space-y-4">
          <Input
            id="listName"
            label="List Name"
            placeholder="My Watchlist"
            value={listName}
            onChange={e => setListName(e.target.value)}
            error={errors.name}
            icon="label"
          />
          <Button onClick={handleCreateList} disabled={submitting} className="w-full">
            {submitting ? <Spinner size={20} /> : 'Create'}
          </Button>
        </div>
      </Modal>

      {/* Add Stock Modal (searchable via backend API) */}
      <Modal isOpen={showAddStock} onClose={() => { if (searchDebounce.current) clearTimeout(searchDebounce.current); setShowAddStock(false); setSearchQuery(''); setSearchResults([]); setSelectedStock(null); setErrors({}); }} title="Add Asset">
        <div className="space-y-4">
          <div className="space-y-2">
            <label className="font-label-caps text-label-caps text-on-surface-variant uppercase block">Search Stock</label>
            <div className="relative">
              <input
                type="text"
                placeholder="Search by symbol or company name..."
                value={searchQuery}
                onChange={e => handleSearch(e.target.value)}
                className="w-full bg-surface-container-low border border-border-muted rounded-lg px-4 py-3 pr-10 text-on-surface focus:outline-none focus:border-electric-crimson focus:ring-1 focus:ring-electric-crimson transition-all font-data-tabular"
              />
              <span className="material-symbols-outlined text-[18px] absolute right-3 top-1/2 -translate-y-1/2 text-on-surface-variant">search</span>
            </div>

            {searchResults.length > 0 && (
              <div className="max-h-48 overflow-y-auto bg-surface-container border border-border-muted rounded-lg">
                {searchResults.map(s => (
                  <button
                    key={s.id}
                    type="button"
                    onClick={() => {
                      setSelectedStock(s);
                      setSearchResults([]);
                      setSearchQuery('');
                    }}
                    className={`w-full flex flex-col items-start px-4 py-2 text-left hover:bg-surface-container-low transition-colors ${
                      selectedStock?.id === s.id ? 'bg-primary/10 border-l-2 border-electric-crimson' : ''
                    }`}
                  >
                    <span className="font-label-caps text-label-caps text-on-surface">{s.symbol}</span>
                    <span className="text-xs text-on-surface-variant">{s.company_name}</span>
                  </button>
                ))}
              </div>
            )}

            {searchQuery.length >= 2 && searchResults.length === 0 && (
              <div className="px-4 py-3 text-sm text-on-surface-variant">No stock found</div>
            )}

            {selectedStock && (
              <div className="flex items-center justify-between px-3 py-2 bg-surface-container-low rounded-lg">
                <span className="text-sm text-on-surface-variant">Selected:</span>
                <span className="font-label-caps text-label-caps text-on-surface">{selectedStock.symbol}</span>
              </div>
            )}

            {errors.stock_id && <p className="text-error-pure text-xs font-label-caps">{errors.stock_id}</p>}
          </div>
          <Button onClick={handleAddStock} disabled={submitting || !selectedStock} className="w-full">
            {submitting ? <Spinner size={20} /> : 'Add to Watchlist'}
          </Button>
        </div>
      </Modal>

      <OrderModal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        stock={tradeStock}
        walletBalance={wallet?.balance ?? 0}
        onSuccess={() => { /* optional refresh */ }}
      />
    </div>
  );
}