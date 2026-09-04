export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}

export interface Stock {
  id: number;
  exchange_instrument_id: string;
  segment: string;
  instrument_type: string;
  symbol: string;
  display_name: string;
  company_name: string;
  isin: string;
  series: string;
  exchange: string;
  contract_expiration: string;
  strike: number;
  option_type: string;
  underlying_symbol_id: string;
  underlying_symbol: string;
  lot_size: number;
  tick_size: number;
  upper_circuit: number;
  lower_circuit: number;
  freeze_qty: number;
  description: string;
  ltp: number;
  open: number;
  high: number;
  low: number;
  close: number;
  vol: number;
  oi: number;
  bid: number;
  ask: number;
  bid_qty: number;
  ask_qty: number;
  cautionary_message_info: string;
  last_updated: string;
}

export interface StockResponse {
  id: number;
  exchange_instrument_id: string;
  segment: string;
  instrument_type: string;
  symbol: string;
  display_name: string;
  company_name: string;
  isin: string;
  series: string;
  exchange: string;
  contract_expiration: string;
  strike: number;
  option_type: string;
  underlying_symbol_id: string;
  underlying_symbol: string;
  lot_size: number;
  tick_size: number;
  upper_circuit: number;
  lower_circuit: number;
  freeze_qty: number;
  description: string;
  ltp: number;
  open: number;
  high: number;
  low: number;
  close: number;
  vol: number;
  oi: number;
  bid: number;
  ask: number;
  bid_qty: number;
  ask_qty: number;
  cautionary_message_info: string;
  last_updated: string;
}

export interface Watchlist {
  id: number;
  user_id: number;
  name: string;
  created_at: string;
  stock_count: number;
}

export interface WatchlistItem {
  id: number;
  watchlist_id: number;
  stock_id: number;
  added_at: string;
  stock?: StockResponse;
}

export interface Order {
  id: number;
  user_id: number;
  stock_id: number;
  side: 'BUY' | 'SELL';
  order_type: 'MARKET' | 'LIMIT';
  quantity: number;
  price: number | null;
  status: 'PENDING' | 'FILLED' | 'CANCELLED' | 'REJECTED';
  created_at: string;
  updated_at: string;
  stock?: StockResponse;
}

export interface CreateOrderRequest {
  stock_id: number;
  side: 'BUY' | 'SELL';
  order_type: 'MARKET' | 'LIMIT';
  quantity: number;
  price?: number;
}

export interface Wallet {
  id: number;
  user_id: number;
  balance: number;
  created_at: string;
  updated_at: string;
}

export interface WalletTransaction {
  id: number;
  wallet_id: number;
  user_id: number;
  type: string;
  amount: number;
  balance_after: number;
  reference_type: string | null;
  reference_id: number | null;
  description: string | null;
  created_at: string;
}

export interface Portfolio {
  id: number;
  user_id: number;
  stock_id: number;
  quantity: number;
  avg_price: number;
  created_at: string;
  updated_at: string;
  stock?: StockResponse;
  current_value: number;
  unrealized_pnl: number;
  pnl_percentage: number;
}

export interface PortfolioSummary {
  total_value: number;
  total_cost_basis: number;
  total_unrealized_pnl: number;
  total_pnl_percentage: number;
  holdings_count: number;
}

export interface Trade {
  id: number;
  order_id: number;
  user_id: number;
  stock_id: number;
  side: string;
  quantity: number;
  execution_price: number;
  total_amount: number;
  executed_at: string;
  stock?: StockResponse;
}

export interface DepositRequest {
  amount: number;
  description?: string;
}

export interface WithdrawRequest {
  amount: number;
  description?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}

export interface ValidationError {
  field: string;
  message: string;
}