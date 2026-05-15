import { Navigate } from 'react-router-dom';

import { useAuth } from '../context/AuthContext';
import { ProtectedRoute } from './ProtectedRoute';

export function OrganizerRoute({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute>
      <OrganizerOnly>{children}</OrganizerOnly>
    </ProtectedRoute>
  );
}

function OrganizerOnly({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();

  if (loading || !user) {
    return null;
  }

  if (user.platform_role !== 'organizer') {
    return <Navigate to="/" replace />;
  }

  return children;
}
