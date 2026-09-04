import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useToast } from '../context/ToastContext';
import { registerSchema, formatValidationErrors } from '../lib/validation';
import Input from '../components/ui/Input';
import Button from '../components/ui/Button';
import Spinner from '../components/ui/Spinner';

export default function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const { showToast } = useToast();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const result = registerSchema.safeParse({ name, email, password, confirmPassword });
    if (!result.success) {
      setErrors(formatValidationErrors(result.error));
      return;
    }
    setErrors({});
    setLoading(true);
    try {
      await register(name, email, password);
      showToast('Account created successfully', 'success');
      navigate('/');
    } catch (err: any) {
      showToast(err.message || 'Registration failed', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-surface-dark text-on-surface font-body-md min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-surface-container rounded-xl border border-border-muted overflow-hidden shadow-2xl relative">
        <div className="absolute inset-0 bg-gradient-to-br from-electric-crimson/5 to-transparent pointer-events-none" />
        <div className="p-8 relative z-10">
          <div className="flex justify-center mb-8">
            <div className="w-16 h-16 rounded-full bg-electric-crimson flex items-center justify-center">
              <span className="font-headline-md text-headline-md font-black text-on-primary text-2xl">B</span>
            </div>
          </div>
          <h1 className="font-headline-md text-headline-md text-center mb-2">Create your account</h1>
          <p className="font-body-md text-body-md text-on-surface-variant text-center mb-8">Join Boro for high-frequency trading.</p>

          <form className="space-y-4" onSubmit={handleSubmit}>
            <Input
              id="fullName"
              label="Full Name"
              placeholder="John Doe"
              value={name}
              onChange={e => setName(e.target.value)}
              error={errors.name}
              icon="person"
            />
            <Input
              id="email"
              label="Email Address"
              type="email"
              placeholder="name@example.com"
              value={email}
              onChange={e => setEmail(e.target.value)}
              error={errors.email}
              icon="mail"
            />
            <Input
              id="password"
              label="Password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={e => setPassword(e.target.value)}
              error={errors.password}
              icon="lock"
            />
            <Input
              id="confirmPassword"
              label="Confirm Password"
              type="password"
              placeholder="••••••••"
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
              error={errors.confirmPassword}
              icon="lock"
            />

            <div className="pt-4">
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? <Spinner size={20} /> : (
                  <>
                    Create Account
                    <span className="material-symbols-outlined">arrow_forward</span>
                  </>
                )}
              </Button>
            </div>
          </form>

          <div className="mt-6 text-center">
            <Link to="/login" className="font-body-md text-body-md text-on-surface-variant hover:text-electric-crimson transition-colors">
              Already have an account?{' '}
              <span className="text-on-surface font-semibold underline decoration-border-muted hover:decoration-electric-crimson underline-offset-4">
                Login
              </span>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}