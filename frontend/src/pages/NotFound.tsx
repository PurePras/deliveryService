import { Link } from 'react-router-dom';
import styles from './PlaceholderPage.module.css';

function NotFound() {
  return (
    <section className={styles.placeholder}>
      <span className={styles.icon} aria-hidden="true">
        🍂
      </span>
      <h1>404</h1>
      <p>This page doesn't exist.</p>
      <Link to="/">Go back home</Link>
    </section>
  );
}

export default NotFound;
