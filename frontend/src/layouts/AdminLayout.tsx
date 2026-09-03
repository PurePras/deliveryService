import { Link, NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import styles from './AdminLayout.module.css';

const links = [
  { to: '/admin', label: 'Dashboard', end: true },
  { to: '/admin/products', label: 'Products' },
  { to: '/admin/categories', label: 'Categories' },
  { to: '/admin/orders', label: 'Orders' },
  { to: '/admin/delivery-areas', label: 'Delivery Areas' },
  { to: '/admin/delivery-slots', label: 'Delivery Slots' },
  { to: '/admin/settings', label: 'Settings' },
];

function AdminLayout() {
  const { user } = useAuth();

  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.brand}>
          <span aria-hidden="true">🥬</span> Admin
        </div>

        <nav className={styles.nav}>
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.end}
              className={({ isActive }) => `${styles.navLink} ${isActive ? styles.navLinkActive : ''}`}
            >
              {link.label}
            </NavLink>
          ))}
        </nav>

        <div className={styles.footer}>
          <span className={styles.adminName}>{user?.name}</span>
          <Link to="/" className={styles.backLink}>
            ← Back to Store
          </Link>
        </div>
      </aside>

      <main className={styles.content}>
        <Outlet />
      </main>
    </div>
  );
}

export default AdminLayout;
