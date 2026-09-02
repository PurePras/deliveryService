import styles from './PlaceholderPage.module.css';

function Checkout() {
  return (
    <section className={styles.placeholder}>
      <span className={styles.icon} aria-hidden="true">
        🧾
      </span>
      <h1>Checkout</h1>
      <p>Checkout is coming soon.</p>
    </section>
  );
}

export default Checkout;
