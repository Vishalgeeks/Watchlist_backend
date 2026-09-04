import { useNavigate } from 'react-router-dom';

export default function TopNav() {
  const navigate = useNavigate();
  return (
    <header className="md:hidden flex justify-between items-center h-16 px-margin-mobile w-full bg-surface-dark/80 backdrop-blur-xl border-b border-border-muted fixed top-0 z-50">
      <div className="flex items-center gap-2" onClick={() => navigate('/')} role="button">
        <div className="w-8 h-8 rounded-lg bg-electric-crimson flex items-center justify-center">
          <span className="font-headline-md-mobile text-headline-md-mobile font-black text-on-primary">B</span>
        </div>
        <span className="font-headline-md-mobile text-headline-md-mobile font-bold text-on-surface tracking-tighter">My Demat</span>
      </div>
      <div className="flex items-center gap-4 text-electric-crimson">
        <button className="hover:text-primary transition-colors duration-200">
          <span className="material-symbols-outlined">notifications</span>
        </button>
        <div className="w-8 h-8 rounded-full bg-surface-container-high border border-border-muted overflow-hidden flex items-center justify-center">
          <span className="material-symbols-outlined text-on-surface-variant text-sm">person</span>
        </div>
      </div>
    </header>
  );
}