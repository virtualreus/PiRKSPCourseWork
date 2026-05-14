import { Navigate } from 'react-router-dom';

import { useAuth } from '../context/AuthContext';

export function GuestRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div className="page-loading">Загрузка…</div>;
  }

  if (user) {
    return <Navigate to="/" replace />;
  }

  return children;
}
