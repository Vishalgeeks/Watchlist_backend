import { useState, useEffect } from 'react';
import { apiGet, apiPut } from '../lib/api';
import { useToast } from '../context/ToastContext';
import type { Order } from '../types';
import { formatCurrency } from '../lib/format';
import Spinner from '../components/ui/Spinner';

type Filter = 'All' | 'Open' | 'Completed' | 'Cancelled';

function statusBadge(status: Order['status']) {
  switch (status) {
    case 'FILLED':
      return <span className="badge-completed">Completed</span>;
    case 'PENDING':
      return <span className="badge-pending">Pending</span>;
    case 'CANCELLED':
    case 'REJECTED':
      return <span className="badge-cancelled">Cancelled</span>;
    default:
      return <span className="badge-pending">{status}</span>;
  }
}

export default function Orders() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<Filter>('All');
  const { showToast } = useToast();

  useEffect(() => {
    const load = async () => {
      try {
        const data = await apiGet<Order[]>('/orders');
        setOrders(data);
      } catch (err: any) {
        showToast(err.message || 'Failed to load orders', 'error');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [showToast]);

  const handleCancel = async (id: number) => {
    try {
      await apiPut(`/orders/${id}/cancel`);
      setOrders(prev => prev.map(o => o.id === id ? { ...o, status: 'CANCELLED' } : o));
      showToast('Order cancelled', 'success');
    } catch (err: any) {
      showToast(err.message || 'Failed to cancel order', 'error');
    }
  };

  const filtered = orders.filter(o => {
    if (filter === 'All') return true;
    if (filter === 'Open') return o.status === 'PENDING';
    if (filter === 'Completed') return o.status === 'FILLED';
    if (filter === 'Cancelled') return o.status === 'CANCELLED' || o.status === 'REJECTED';
    return true;
  });

  const filters: Filter[] = ['All', 'Open', 'Completed', 'Cancelled'];

  return (
    <div className="flex flex-col min-h-screen">
      {/* Page Header */}
      <div className="pt-8 pb-6 border-b border-border-muted bg-surface-dark/50 sticky top-0 z-40">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div>
            <h1 className="font-headline-md text-headline-md md:font-display-lg md:text-display-lg font-bold text-on-surface mb-2">Orders</h1>
            <p className="text-on-surface-variant text-sm">Manage your active and historical trades.</p>
          </div>
          <div className="flex bg-surface-container rounded-lg p-1 border border-border-muted overflow-x-auto">
            {filters.map(f => (
              <button
                key={f}
                onClick={() => setFilter(f)}
                className={`px-4 py-1.5 rounded-md font-label-caps text-label-caps whitespace-nowrap transition-colors ${
                  filter === f ? 'bg-surface-variant text-on-surface shadow-sm border border-border-muted' : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-container-high'
                }`}
              >
                {f}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Orders Table */}
      <div className="flex-1 p-margin-mobile md:p-margin-desktop overflow-x-auto">
        {loading ? (
          <div className="flex items-center justify-center h-64"><Spinner size={32} /></div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <div className="w-16 h-16 rounded-full bg-surface-container flex items-center justify-center mb-4 border border-border-muted">
              <span className="material-symbols-outlined text-on-surface-variant text-3xl">receipt_long</span>
            </div>
            <h3 className="font-headline-md text-headline-md text-on-surface mb-2">No orders yet</h3>
            <p className="text-on-surface-variant mb-6 max-w-sm">You haven't placed any trades. Head over to the markets to start building your portfolio.</p>
          </div>
        ) : (
          <div className="min-w-[800px] w-full bg-surface-dark rounded-xl border border-border-muted overflow-hidden shadow-2xl relative">
            <div className="grid grid-cols-8 gap-4 px-6 py-4 border-b border-border-muted bg-surface-container-low font-label-caps text-label-caps text-on-surface-variant">
              <div className="col-span-2">Stock</div>
              <div>Type</div>
              <div>Order</div>
              <div className="text-right">Quantity</div>
              <div className="text-right">Price</div>
              <div className="text-right">Total</div>
              <div className="text-center">Status</div>
            </div>
            <div className="flex flex-col">
              {filtered.map(order => {
                const price = order.price ?? order.stock?.ltp ?? 0;
                const total = price * order.quantity;
                const isBuy = order.side === 'BUY';
                return (
                  <div key={order.id} className="grid grid-cols-8 gap-4 px-6 py-4 border-b border-border-muted/50 hover:bg-surface-light transition-colors group items-center">
                    <div className="col-span-2 flex flex-col">
                      <span className="font-label-caps text-label-caps text-on-surface">{order.stock?.symbol}</span>
                      <span className="text-xs text-on-surface-variant truncate">{order.stock?.company_name}</span>
                    </div>
                    <div className={`font-data-tabular text-data-tabular ${isBuy ? 'text-tertiary' : 'text-error-pure'}`}>{order.side}</div>
                    <div className="font-data-tabular text-data-tabular text-on-surface">{order.order_type}</div>
                    <div className="font-data-tabular text-data-tabular text-on-surface text-right">{order.quantity}</div>
                    <div className="font-data-tabular text-data-tabular text-on-surface text-right">{formatCurrency(price)}</div>
                    <div className="font-data-tabular text-data-tabular text-on-surface text-right">{formatCurrency(total)}</div>
                    <div className="flex justify-center items-center gap-2">
                      {statusBadge(order.status)}
                      {order.status === 'PENDING' && (
                        <button
                          onClick={() => handleCancel(order.id)}
                          className="text-on-surface-variant hover:text-error-pure opacity-0 group-hover:opacity-100 transition-opacity"
                          title="Cancel order"
                        >
                          <span className="material-symbols-outlined text-[18px]">close</span>
                        </button>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}