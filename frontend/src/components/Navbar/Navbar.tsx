import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useCart } from '../../store/cart/CartContext';
import { useAuth } from '../../store/auth/AuthContext';
import styles from './Navbar.module.css';

const links = [
  { to: '/products', label: 'Products' },
  { to: '/categories', label: 'Categories' },
];

function Navbar() {
  const [isOpen, setIsOpen] = useState(false);
  const { itemCount } = useCart();
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    setIsOpen(false);
    await logout();
    navigate('/');
  }

  return (
    <header className={styles.navbar}>
      <div className={styles.inner}>
        <NavLink to="/" className={styles.brand} onClick={() => setIsOpen(false)}>
          <span className={styles.brandIcon} aria-hidden="true">
            🥬
          </span>
          Shri Ram Sabji Delivery
        </NavLink>

        <div className={styles.right}>
          <nav className={`${styles.links} ${isOpen ? styles.linksOpen : ''}`}>
            {links.map((link) => (
              <NavLink
                key={link.to}
                to={link.to}
                className={({ isActive }) => `${styles.link} ${isActive ? styles.linkActive : ''}`}
                onClick={() => setIsOpen(false)}
              >
                {link.label}
              </NavLink>
            ))}

            {user ? (
              <>
                <NavLink
                  to="/orders"
                  className={({ isActive }) => `${styles.link} ${isActive ? styles.linkActive : ''}`}
                  onClick={() => setIsOpen(false)}
                >
                  Orders
                </NavLink>
                <NavLink
                  to="/account"
                  className={({ isActive }) => `${styles.link} ${isActive ? styles.linkActive : ''}`}
                  onClick={() => setIsOpen(false)}
                >
                  {user.name.split(' ')[0]}
                </NavLink>
                {user.role === 'admin' && (
                  <NavLink
                    to="/admin"
                    className={({ isActive }) => `${styles.link} ${isActive ? styles.linkActive : ''}`}
                    onClick={() => setIsOpen(false)}
                  >
                    Admin
                  </NavLink>
                )}
                <button type="button" className={styles.link} onClick={handleLogout}>
                  Logout
                </button>
              </>
            ) : (
              <NavLink
                to="/login"
                className={({ isActive }) => `${styles.link} ${isActive ? styles.linkActive : ''}`}
                onClick={() => setIsOpen(false)}
              >
                Login
              </NavLink>
            )}
          </nav>

          <NavLink to="/cart" className={styles.cart} aria-label="Cart" onClick={() => setIsOpen(false)}>
            <span aria-hidden="true">🛒</span>
            {itemCount > 0 && <span className={styles.cartBadge}>{itemCount}</span>}
          </NavLink>

          <button
            type="button"
            className={styles.toggle}
            aria-label="Toggle navigation menu"
            aria-expanded={isOpen}
            onClick={() => setIsOpen((open) => !open)}
          >
            <span className={styles.toggleBar} />
            <span className={styles.toggleBar} />
            <span className={styles.toggleBar} />
          </button>
        </div>
      </div>
    </header>
  );
}

export default Navbar;
