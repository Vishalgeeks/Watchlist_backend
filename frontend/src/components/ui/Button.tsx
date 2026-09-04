import { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost';
  children: ReactNode;
}

export default function Button({ variant = 'primary', children, className = '', disabled, ...props }: ButtonProps) {
  const base = 'font-label-caps text-label-caps px-6 py-2.5 rounded-lg transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed';
  const variants = {
    primary: 'bg-electric-crimson hover:bg-primary-container text-on-primary shadow-[0_0_15px_rgba(255,0,94,0.3)]',
    secondary: 'bg-surface-container-high hover:bg-surface-bright border border-border-muted text-on-surface',
    ghost: 'text-on-surface-variant hover:text-on-surface',
  };
  return (
    <button className={`${base} ${variants[variant]} ${className}`} disabled={disabled} {...props}>
      {children}
    </button>
  );
}