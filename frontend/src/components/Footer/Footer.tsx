import { useEffect, useState } from 'react';
import { getHealth } from '../../api/health';
import styles from './Footer.module.css';

type ApiStatus = 'checking' | 'ok' | 'down';

function Footer() {
  const [status, setStatus] = useState<ApiStatus>('checking');

  useEffect(() => {
    let cancelled = false;

    getHealth()
      .then((health) => {
        if (!cancelled) {
          setStatus(health.status === 'ok' ? 'ok' : 'down');
        }
      })
      .catch(() => {
        if (!cancelled) {
          setStatus('down');
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const statusLabel = status === 'checking' ? 'checking…' : status === 'ok' ? 'online' : 'offline';

  return (
    <footer className={styles.footer}>
      <p className={styles.tagline}>ताज़ी सब्ज़ियाँ, अब आपके घर तक</p>
      <div className={styles.meta}>
        <span>&copy; {new Date().getFullYear()} Made by Prashant</span>
        <span className={styles.status}>
          <span className={`${styles.dot} ${styles[status]}`} />
          API {statusLabel}
        </span>
      </div>
    </footer>
  );
}

export default Footer;
