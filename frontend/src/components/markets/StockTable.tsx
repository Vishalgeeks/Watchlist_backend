import { useNavigate } from 'react-router-dom';
import type { StockResponse } from '../../types';
import { formatCurrency, formatPercent } from '../../lib/format';

interface StockTableProps {
  stocks: StockResponse[];
  onTrade: (stock: StockResponse) => void;
  filter: 'All' | 'Gainers' | 'Losers' | 'Most Active';
}

export default function StockTable({ stocks, onTrade, filter }: StockTableProps) {
  const navigate = useNavigate();
  const filtered = stocks.filter(s => {
    if (filter === 'Gainers') return s.ltp >= s.close;
    if (filter === 'Losers') return s.ltp < s.close;
    return true;
  });

  if (filtered.length === 0) {
    return <p className="text-on-surface-variant text-center py-8">No stocks match the current filter.</p>;
  }

  return (
    <div className="overflow-x-auto w-full">
      <table className="w-full text-left border-collapse min-w-[800px]">
        <thead>
          <tr className="border-b border-border-muted bg-surface-dark/30">
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium w-1/4">Stock</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-right">Price</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-right">Change</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-right">Change %</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-right">Volume</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-right">Market Cap</th>
            <th className="p-4 font-label-caps text-label-caps text-on-surface-variant font-medium text-center">Action</th>
          </tr>
        </thead>
        <tbody className="font-data-tabular text-data-tabular divide-y divide-border-muted/50">
          {filtered.map(stock => {
            const change = stock.ltp - stock.close;
            const changePct = stock.close > 0 ? (change / stock.close) * 100 : 0;
            const positive = change >= 0;
            return (
              <tr key={stock.id} className="hover:bg-surface-light transition-colors group cursor-pointer" onClick={() => navigate(`/stocks/${stock.id}`)}>
                <td className="p-4">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded bg-surface-container flex items-center justify-center font-bold text-on-surface text-xs border border-border-muted">
                      {stock.symbol.slice(0, 4)}
                    </div>
                    <div>
                      <div className="text-on-surface font-body-md font-semibold text-sm">{stock.symbol}</div>
                      <div className="text-on-surface-variant text-xs font-body-md">{stock.company_name}</div>
                    </div>
                  </div>
                </td>
                <td className="p-4 text-right text-on-surface">{formatCurrency(stock.ltp)}</td>
                <td className={`p-4 text-right ${positive ? 'text-tertiary' : 'text-electric-crimson'}`}>
                  {positive ? '+' : ''}{formatCurrency(change)}
                </td>
                <td className="p-4 text-right">
                  <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs ${positive ? 'bg-tertiary/10 text-tertiary' : 'bg-error-container/20 text-electric-crimson'}`}>
                    <span className="material-symbols-outlined text-[10px]">{positive ? 'arrow_upward' : 'arrow_downward'}</span>
                    {formatPercent(changePct)}
                  </span>
                </td>
                <td className="p-4 text-right text-on-surface-variant">{stock.vol.toLocaleString()}</td>
                <td className="p-4 text-right text-on-surface-variant">-</td>
                <td className="p-4 text-center" onClick={e => e.stopPropagation()}>
                  <button
                    onClick={() => onTrade(stock)}
                    className="bg-primary-container/20 text-electric-crimson hover:bg-primary-container hover:text-on-primary-container px-4 py-1.5 rounded-lg text-xs font-label-caps text-label-caps transition-all opacity-0 group-hover:opacity-100"
                  >
                    Trade
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}