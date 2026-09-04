import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import AppLayout from './components/layout/AppLayout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Markets from './pages/Markets';
import StockDetails from './pages/StockDetails';
import Portfolio from './pages/Portfolio';
import Orders from './pages/Orders';
import Watchlist from './pages/Watchlist';

function RequireAuth({ children }: { children: JSX.Element }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return children;
}

export default function App() {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      <Route path="/login" element={isAuthenticated ? <Navigate to="/" replace /> : <Login />} />
      <Route path="/register" element={isAuthenticated ? <Navigate to="/" replace /> : <Register />} />
      <Route
        path="/"
        element={
          <RequireAuth>
            <AppLayout>
              <Dashboard />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route
        path="/markets"
        element={
          <RequireAuth>
            <AppLayout>
              <Markets />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route
        path="/stocks/:id"
        element={
          <RequireAuth>
            <AppLayout>
              <StockDetails />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route
        path="/portfolio"
        element={
          <RequireAuth>
            <AppLayout>
              <Portfolio />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route
        path="/orders"
        element={
          <RequireAuth>
            <AppLayout>
              <Orders />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route
        path="/watchlist"
        element={
          <RequireAuth>
            <AppLayout>
              <Watchlist />
            </AppLayout>
          </RequireAuth>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}