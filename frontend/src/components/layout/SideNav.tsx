import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

interface NavItem {
  path: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { path: '/', label: 'Dashboard', icon: 'dashboard' },
  { path: '/markets', label: 'Markets', icon: 'bar_chart' },
  { path: '/portfolio', label: 'Portfolio', icon: 'account_balance_wallet' },
  { path: '/orders', label: 'Orders', icon: 'receipt_long' },
  { path: '/watchlist', label: 'Watchlist', icon: 'visibility' },
];

export default function SideNav() {
  const { logout, user } = useAuth();
  const location = useLocation();

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  return (
    <nav className="hidden md:flex flex-col py-panel-padding bg-surface-dark fixed left-0 top-0 h-full w-60 border-r border-border-muted z-40">
      <div className="px-panel-padding mb-8">
        <h1 className="font-headline-md text-headline-md font-black text-on-surface tracking-tighter">My Demat</h1>
        <p className="font-label-caps text-label-caps text-on-surface-variant mt-1">Trading Terminal</p>
      </div>

      <div className="flex-1 px-4 space-y-2">
        {navItems.map(item => (
          <NavLink
            key={item.path}
            to={item.path}
            className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
              isActive(item.path)
                ? 'text-electric-crimson border-r-2 border-electric-crimson bg-primary/5 translate-x-1'
                : 'text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface'
            }`}
          >
            <span className="material-symbols-outlined">{item.icon}</span>
            <span className="font-label-caps text-label-caps">{item.label}</span>
          </NavLink>
        ))}
      </div>

      <div className="px-4 mt-auto space-y-2">
        <div className="px-4 py-2 text-on-surface-variant font-label-caps text-label-caps truncate">
          {user?.name ?? 'Trader'}
        </div>
        <button
          onClick={logout}
          className="flex items-center gap-3 px-4 py-3 rounded-lg text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface transition-all w-full text-left"
        >
          <span className="material-symbols-outlined">logout</span>
          <span className="font-label-caps text-label-caps">Logout</span>
        </button>
      </div>
    </nav>
  );
}