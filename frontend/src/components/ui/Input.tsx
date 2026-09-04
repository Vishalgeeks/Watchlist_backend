import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  icon?: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(({ label, error, icon, className = '', ...props }, ref) => {
  return (
    <div className="space-y-2">
      {label && (
        <label className="font-label-caps text-label-caps text-on-surface-variant uppercase block" htmlFor={props.id}>
          {label}
        </label>
      )}
      <div className="relative">
        {icon && (
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm">
            {icon}
          </span>
        )}
        <input
          ref={ref}
          className={`w-full bg-surface-container-low border rounded-lg px-4 py-3 text-on-surface placeholder:text-surface-variant
            focus:outline-none focus:border-electric-crimson focus:ring-1 focus:ring-electric-crimson transition-all
            ${icon ? 'pl-10' : ''} ${error ? 'border-error-pure' : 'border-border-muted'} ${className}`}
          {...props}
        />
      </div>
      {error && <p className="text-error-pure text-xs font-label-caps">{error}</p>}
    </div>
  );
});

Input.displayName = 'Input';
export default Input;