import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import StateMessage from '../components/StateMessage/StateMessage';

function AdminRoute() {
  const { user, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <StateMessage icon="⏳" title="Loading…" />;
  }

  if (!user) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  }

  // Logged in but not an admin — send them back to the store, not to /login
  // (they're already authenticated, just not authorized for this).
  if (user.role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
}

export default AdminRoute;
