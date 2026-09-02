import { useNavigate } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import styles from './Account.module.css';

function Account() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  // ProtectedRoute guarantees user is set whenever this page actually renders.
  if (!user) return null;

  async function handleLogout() {
    await logout();
    navigate('/');
  }

  return (
    <section className={styles.page}>
      <h1>Welcome, {user.name}</h1>

      <div className={styles.card}>
        <div className={styles.row}>
          <span className={styles.label}>Phone</span>
          <span>{user.phone}</span>
        </div>
        {user.email && (
          <div className={styles.row}>
            <span className={styles.label}>Email</span>
            <span>{user.email}</span>
          </div>
        )}
        <div className={styles.row}>
          <span className={styles.label}>Account type</span>
          <span className={styles.roleBadge}>{user.role}</span>
        </div>
      </div>

      <button type="button" className={styles.logoutBtn} onClick={handleLogout}>
        Log Out
      </button>
    </section>
  );
}

export default Account;
