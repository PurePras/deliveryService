import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../store/cart/CartContext';
import { listDeliveryAreas } from '../api/deliveryAreas';
import { listDeliverySlots } from '../api/deliverySlots';
import { createOrder } from '../api/orders';
import { useAsync } from '../hooks/useAsync';
import ProductImage from '../components/ProductImage/ProductImage';
import StateMessage from '../components/StateMessage/StateMessage';
import { formatPrice } from '../utils/format';
import { ApiError } from '../api/http';
import styles from './Checkout.module.css';

function todayISO(): string {
  const now = new Date();
  const local = new Date(now.getTime() - now.getTimezoneOffset() * 60000);
  return local.toISOString().slice(0, 10);
}

function Checkout() {
  const { items, subtotal, clearCart } = useCart();
  const navigate = useNavigate();
  const { data: areas } = useAsync(() => listDeliveryAreas({ is_active: true, limit: 100 }), []);
  const { data: slots } = useAsync(() => listDeliverySlots({ is_active: true, limit: 100 }), []);

  const [address, setAddress] = useState('');
  const [areaId, setAreaId] = useState('');
  const [slotId, setSlotId] = useState('');
  const [date, setDate] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const created = await createOrder({
        delivery_area_id: areaId,
        delivery_slot_id: slotId,
        delivery_date: date,
        delivery_address: address,
        items: items.map((item) => ({ product_id: item.productId, quantity: String(item.quantity) })),
      });
      clearCart();
      navigate(`/orders/${created.id}`, { state: { justPlaced: true }, replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
      setSubmitting(false);
    }
  }

  if (items.length === 0) {
    return (
      <StateMessage icon="🛒" title="Your cart is empty" description="Add some fresh produce before checking out.">
        <Link to="/products" className={styles.browseLink}>
          Browse products →
        </Link>
      </StateMessage>
    );
  }

  return (
    <section className={styles.page}>
      <h1>Checkout</h1>

      <div className={styles.layout}>
        <ul className={styles.review}>
          {items.map((item) => (
            <li key={item.productId} className={styles.reviewItem}>
              <div className={styles.reviewImage}>
                <ProductImage src={item.imageUrl} alt={item.name} />
              </div>
              <div className={styles.reviewInfo}>
                <span className={styles.reviewName}>{item.name}</span>
                <span className={styles.reviewQty}>
                  {item.quantity} × {formatPrice(item.price)}
                </span>
              </div>
              <span className={styles.reviewSubtotal}>{formatPrice(Number(item.price) * item.quantity)}</span>
            </li>
          ))}
          <li className={styles.reviewTotal}>
            <span>Subtotal</span>
            <span>{formatPrice(subtotal)}</span>
          </li>
        </ul>

        <form className={styles.form} onSubmit={handleSubmit}>
          {error && <p className={styles.error}>{error}</p>}

          <label className={styles.field}>
            <span>Delivery address</span>
            <textarea
              required
              rows={3}
              value={address}
              onChange={(event) => setAddress(event.target.value)}
              placeholder="House no., street, landmark…"
            />
          </label>

          <label className={styles.field}>
            <span>Delivery area</span>
            <select required value={areaId} onChange={(event) => setAreaId(event.target.value)}>
              <option value="" disabled>
                Select an area
              </option>
              {areas?.map((area) => (
                <option key={area.id} value={area.id}>
                  {area.name}, {area.city} ({area.pincode})
                </option>
              ))}
            </select>
          </label>

          <label className={styles.field}>
            <span>Delivery slot</span>
            <select required value={slotId} onChange={(event) => setSlotId(event.target.value)}>
              <option value="" disabled>
                Select a slot
              </option>
              {slots?.map((slot) => (
                <option key={slot.id} value={slot.id}>
                  {slot.label}
                </option>
              ))}
            </select>
          </label>

          <label className={styles.field}>
            <span>Delivery date</span>
            <input type="date" required min={todayISO()} value={date} onChange={(event) => setDate(event.target.value)} />
          </label>

          <button type="submit" className={styles.submit} disabled={submitting}>
            {submitting ? 'Placing order…' : `Place Order · ${formatPrice(subtotal)}`}
          </button>
        </form>
      </div>
    </section>
  );
}

export default Checkout;
