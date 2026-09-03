import { Link } from 'react-router-dom';
import { listOrders } from '../api/orders';
import { useAsync } from '../hooks/useAsync';
import StateMessage from '../components/StateMessage/StateMessage';
import OrderStatusBadge from '../components/OrderStatusBadge/OrderStatusBadge';
import { formatDate, formatPrice } from '../utils/format';
import styles from './Orders.module.css';

function Orders() {
  const { data: orders, loading, error } = useAsync(() => listOrders({ limit: 50 }), []);

  if (loading) {
    return <StateMessage icon="⏳" title="Loading your orders…" />;
  }

  if (error) {
    return <StateMessage icon="⚠️" title="Couldn't load your orders" description="Please try again in a moment." />;
  }

  if (!orders || orders.length === 0) {
    return (
      <StateMessage icon="📦" title="No orders yet" description="Once you place an order, it'll show up here.">
        <Link to="/products" className={styles.browseLink}>
          Browse products →
        </Link>
      </StateMessage>
    );
  }

  return (
    <section className={styles.page}>
      <h1>Your Orders</h1>

      <ul className={styles.list}>
        {orders.map((order) => (
          <li key={order.id}>
            <Link to={`/orders/${order.id}`} className={styles.card}>
              <div className={styles.cardHeader}>
                <span className={styles.orderId}>#{order.id.slice(0, 8)}</span>
                <OrderStatusBadge status={order.status} />
              </div>
              <div className={styles.cardMeta}>
                <span>{formatDate(order.created_at)}</span>
                <span className={styles.total}>{formatPrice(order.total_amount)}</span>
              </div>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}

export default Orders;
