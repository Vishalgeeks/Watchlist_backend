import { useState } from 'react';
import Modal from '../ui/Modal';
import Button from '../ui/Button';
import Input from '../ui/Input';
import Spinner from '../ui/Spinner';
import { useToast } from '../../context/ToastContext';
import { apiPost } from '../../lib/api';
import { createOrderSchema, formatValidationErrors } from '../../lib/validation';
import type { StockResponse, Order } from '../../types';
import { formatCurrency } from '../../lib/format';

interface OrderModalProps {
  isOpen: boolean;
  onClose: () => void;
  stock: StockResponse | null;
  walletBalance: number;
  defaultSide?: 'BUY' | 'SELL';
  onSuccess?: () => void;
}

export default function OrderModal({ isOpen, onClose, stock, walletBalance, defaultSide = 'BUY', onSuccess }: OrderModalProps) {
  const [side, setSide] = useState<'BUY' | 'SELL'>(defaultSide);
  const [orderType, setOrderType] = useState<'MARKET' | 'LIMIT'>('MARKET');
  const [quantity, setQuantity] = useState(1);
  const [price, setPrice] = useState<number | ''>('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  const resetForm = () => {
    setSide(defaultSide);
    setOrderType('MARKET');
    setQuantity(1);
    setPrice('');
    setErrors({});
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const basePrice = stock ? (stock.ltp > 0 ? stock.ltp : stock.close > 0 ? stock.close : stock.open > 0 ? stock.open : 100) : 100;
  const estCost = stock ? (orderType === 'LIMIT' && price !== '' ? Number(price) : basePrice) * quantity : 0;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stock) return;

    const result = createOrderSchema.safeParse({
      stock_id: stock.id,
      side,
      order_type: orderType,
      quantity,
      price: orderType === 'LIMIT' && price !== '' ? Number(price) : undefined,
    });

    if (!result.success) {
      setErrors(formatValidationErrors(result.error));
      return;
    }
    setErrors({});

    if (side === 'BUY' && estCost > walletBalance) {
      showToast('Insufficient wallet balance', 'error');
      return;
    }

    setLoading(true);
    try {
      const body: Record<string, unknown> = {
        stock_id: stock.id,
        side,
        order_type: orderType,
        quantity,
      };
      if (orderType === 'LIMIT' && price !== '') {
        body.price = Number(price);
      }
      const order = await apiPost<Order>('/orders', body);
      const executionPrice = orderType === 'LIMIT' && price !== '' ? Number(price) : basePrice;
      await apiPost(`/orders/${order.id}/execute`, { execution_price: executionPrice });
      showToast(`${side} order filled successfully`, 'success');
      handleClose();
      onSuccess?.();
    } catch (err: any) {
      showToast(err.message || 'Failed to place order', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title={stock ? `Trade ${stock.symbol}` : 'Trade'} maxWidth="max-w-lg">
      {stock && (
        <form className="space-y-5" onSubmit={handleSubmit}>
          <div className="flex items-center justify-between text-sm">
            <span className="font-label-caps text-label-caps text-on-surface-variant">Available Cash</span>
            <span className="font-data-tabular text-data-tabular text-on-surface">{formatCurrency(walletBalance)}</span>
          </div>

          <div>
            <label className="font-label-caps text-label-caps text-on-surface-variant uppercase block mb-2">Side</label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setSide('BUY')}
                className={`py-2.5 rounded-lg font-label-caps text-label-caps transition-all ${
                  side === 'BUY' ? 'bg-tertiary/20 text-tertiary border border-tertiary/40' : 'bg-surface-dark text-on-surface-variant border border-border-muted'
                }`}
              >
                Buy
              </button>
              <button
                type="button"
                onClick={() => setSide('SELL')}
                className={`py-2.5 rounded-lg font-label-caps text-label-caps transition-all ${
                  side === 'SELL' ? 'bg-error-pure/20 text-error-pure border border-error-pure/40' : 'bg-surface-dark text-on-surface-variant border border-border-muted'
                }`}
              >
                Sell
              </button>
            </div>
          </div>

          <div>
            <label className="font-label-caps text-label-caps text-on-surface-variant uppercase block mb-2">Order Type</label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setOrderType('MARKET')}
                className={`py-2.5 rounded-lg font-label-caps text-label-caps transition-all ${
                  orderType === 'MARKET' ? 'bg-electric-crimson/20 text-electric-crimson border border-electric-crimson/40' : 'bg-surface-dark text-on-surface-variant border border-border-muted'
                }`}
              >
                Market
              </button>
              <button
                type="button"
                onClick={() => setOrderType('LIMIT')}
                className={`py-2.5 rounded-lg font-label-caps text-label-caps transition-all ${
                  orderType === 'LIMIT' ? 'bg-electric-crimson/20 text-electric-crimson border border-electric-crimson/40' : 'bg-surface-dark text-on-surface-variant border border-border-muted'
                }`}
              >
                Limit
              </button>
            </div>
          </div>

          <Input
            id="quantity"
            label="Quantity"
            type="number"
            min={1}
            value={quantity}
            onChange={e => setQuantity(Math.max(1, Number(e.target.value)))}
            error={errors.quantity}
          />

          {orderType === 'LIMIT' && (
            <Input
              id="limitPrice"
              label="Limit Price"
              type="number"
              min={0}
              step="0.01"
              placeholder={stock.ltp.toFixed(2)}
              value={price}
              onChange={e => setPrice(e.target.value === '' ? '' : Number(e.target.value))}
              error={errors.price}
              icon="sell"
            />
          )}

          <div className="flex items-center justify-between p-3 bg-surface-dark rounded-lg border border-border-muted">
            <span className="font-label-caps text-label-caps text-on-surface-variant">Est. Cost</span>
            <span className="font-data-tabular text-data-tabular text-on-surface text-lg font-bold">{formatCurrency(estCost)}</span>
          </div>

          {errors.stock_id && <p className="text-error-pure text-xs font-label-caps">{errors.stock_id}</p>}

          <div className="flex gap-3 pt-2">
            <Button type="submit" className="flex-1" disabled={loading}>
              {loading ? <Spinner size={20} /> : side === 'BUY' ? 'Buy' : 'Sell'}
            </Button>
            <Button type="button" variant="secondary" onClick={handleClose} disabled={loading}>
              Cancel
            </Button>
          </div>
        </form>
      )}
    </Modal>
  );
}