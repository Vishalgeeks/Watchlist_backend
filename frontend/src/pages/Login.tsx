import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../context/ToastContext';
import { loginSchema, formatValidationErrors } from '../lib/validation';
import Input from '../components/ui/Input';
import Button from '../components/ui/Button';
import Spinner from '../components/ui/Spinner';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const { login } = useAuth();
  const { showToast } = useToast();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const result = loginSchema.safeParse({ email, password });
    if (!result.success) {
      setErrors(formatValidationErrors(result.error));
      return;
    }
    setErrors({});
    setLoading(true);
    try {
      await login(email, password);
      showToast('Login successful', 'success');
      navigate('/');
    } catch (err: any) {
      showToast(err.message || 'Login failed', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 antialiased bg-background">
      <div className="fixed inset-0 pointer-events-none opacity-20" style={{ background: 'radial-gradient(circle at 50% 0%, rgba(255,0,94,0.15) 0%, transparent 50%)' }} />
      <main className="w-full max-w-md relative z-10">
        <div className="glass-panel rounded-xl p-8 md:p-10 w-full flex flex-col items-center">
          <div className="mb-8 flex flex-col items-center">
            <div className="w-16 h-16 rounded-full bg-electric-crimson flex items-center justify-center mb-4">
              <span className="font-headline-md text-headline-md font-black text-on-primary text-2xl">B</span>
            </div>
            <h1 className="font-headline-md text-headline-md text-on-surface mb-2 tracking-tight">Welcome back</h1>
            <p className="font-body-md text-body-md text-on-surface-variant text-center">Enter your credentials to access the terminal.</p>
          </div>

          <form className="w-full space-y-6" onSubmit={handleSubmit}>
            <Input
              id="email"
              label="Email Address"
              type="email"
              placeholder="trader@boro.pro"
              value={email}
              onChange={e => setEmail(e.target.value)}
              error={errors.email}
              icon="mail"
            />
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <label className="font-label-caps text-label-caps text-on-surface-variant uppercase block" htmlFor="password">
                  Password
                </label>
                <a className="font-label-caps text-label-caps text-electric-crimson hover:text-primary transition-colors" href="#">
                  Forgot password?
                </a>
              </div>
              <div className="relative group">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  className={`w-full bg-surface-container-low border rounded-lg px-4 py-3 font-data-tabular text-data-tabular tracking-widest placeholder:text-surface-variant placeholder:tracking-normal pr-12
                    focus:outline-none focus:border-electric-crimson focus:ring-1 focus:ring-electric-crimson transition-all
                    ${errors.password ? 'border-error-pure' : 'border-border-muted'}`}
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  aria-label="Toggle password visibility"
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-on-surface-variant hover:text-on-surface transition-colors p-1"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  <span className="material-symbols-outlined text-[20px]">
                    {showPassword ? 'visibility' : 'visibility_off'}
                  </span>
                </button>
              </div>
              {errors.password && <p className="text-error-pure text-xs font-label-caps">{errors.password}</p>}
            </div>

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? <Spinner size={20} /> : (
                <>
                  Login
                  <span className="material-symbols-outlined text-[20px]">arrow_forward</span>
                </>
              )}
            </Button>
          </form>

          <div className="mt-8 pt-6 border-t border-border-muted w-full text-center">
            <p className="font-body-md text-body-md text-on-surface-variant">
              Don't have an account?{' '}
              <Link to="/register" className="text-electric-crimson hover:text-primary transition-colors font-medium">
                Create account
              </Link>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}