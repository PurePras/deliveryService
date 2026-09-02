import { Link } from 'react-router-dom';
import { useCart } from '../store/cart/CartContext';
import ProductImage from '../components/ProductImage/ProductImage';
import QuantityStepper from '../components/QuantityStepper/QuantityStepper';
import StateMessage from '../components/StateMessage/StateMessage';
import { formatPrice } from '../utils/format';
import styles from './Cart.module.css';

function Cart() {
  const { items, itemCount, subtotal, setQuantity, removeItem } = useCart();

  if (items.length === 0) {
    return (
      <StateMessage icon="🛒" title="Your cart is empty" description="Add some fresh produce to get started.">
        <Link to="/products" className={styles.browseLink}>
          Browse products →
        </Link>
      </StateMessage>
    );
  }

  return (
    <section className={styles.page}>
      <h1>Your Cart</h1>

      <div className={styles.layout}>
        <ul className={styles.items}>
          {items.map((item) => (
            <li key={item.productId} className={styles.item}>
              <Link to={`/products/${item.productId}`} className={styles.itemImage}>
                <ProductImage src={item.imageUrl} alt={item.name} />
              </Link>

              <div className={styles.itemInfo}>
                <Link to={`/products/${item.productId}`} className={styles.itemName}>
                  {item.name}
                </Link>
                <p className={styles.itemUnitPrice}>
                  {formatPrice(item.price)} / {item.unit}
                </p>
                {!item.isAvailable && <p className={styles.itemUnavailable}>No longer available</p>}
              </div>

              <div className={styles.itemActions}>
                <QuantityStepper
                  quantity={item.quantity}
                  onIncrease={() => setQuantity(item.productId, item.quantity + 1)}
                  onDecrease={() => setQuantity(item.productId, item.quantity - 1)}
                />

                <p className={styles.itemSubtotal}>{formatPrice(Number(item.price) * item.quantity)}</p>

                <button
                  type="button"
                  className={styles.remove}
                  onClick={() => removeItem(item.productId)}
                  aria-label={`Remove ${item.name} from cart`}
                >
                  ✕
                </button>
              </div>
            </li>
          ))}
        </ul>

        <aside className={styles.summary}>
          <h2>Order Summary</h2>
          <div className={styles.summaryRow}>
            <span>Items ({itemCount})</span>
            <span>{formatPrice(subtotal)}</span>
          </div>
          <div className={styles.summaryTotal}>
            <span>Subtotal</span>
            <span>{formatPrice(subtotal)}</span>
          </div>
          <Link to="/checkout" className={styles.checkoutBtn}>
            Proceed to Checkout
          </Link>
        </aside>
      </div>
    </section>
  );
}

export default Cart;
