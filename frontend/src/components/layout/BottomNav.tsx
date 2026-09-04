import { NavLink, useLocation } from 'react-router-dom';

interface NavItem {
  path: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { path: '/', label: 'Home', icon: 'home' },
  { path: '/markets', label: 'Markets', icon: 'trending_up' },
  { path: '/portfolio', label: 'Portfolio', icon: 'account_balance_wallet' },
  { path: '/orders', label: 'Orders', icon: 'swap_horiz' },
  { path: '/watchlist', label: 'Watchlist', icon: 'person' },
];

export default function BottomNav() {
  const location = useLocation();
  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  return (
    <nav className="md:hidden fixed bottom-6 left-1/2 -translate-x-1/2 w-[92%] z-50 flex justify-around items-center p-2 bg-surface-container/90 backdrop-blur-md rounded-full border border-border-muted shadow-xl">
      {navItems.map(item => (
        <NavLink
          key={item.path}
          to={item.path}
          className={`flex flex-col items-center justify-center px-4 py-2 rounded-full transition-colors ${
            isActive(item.path)
              ? 'bg-primary-container text-on-primary-container scale-95'
              : 'text-on-surface-variant active:bg-surface-variant'
          }`}
        >
          <span className="material-symbols-outlined mb-1" style={isActive(item.path) ? { fontVariationSettings: "'FILL' 1" } : undefined}>
            {item.icon}
          </span>
          <span className="font-label-caps text-label-caps text-[10px]">{item.label}</span>
        </NavLink>
      ))}
    </nav>
  );
}