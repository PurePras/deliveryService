import { Link, useLocation, useParams } from 'react-router-dom';
import { getOrder } from '../api/orders';
import { getProduct } from '../api/products';
import { useAsync } from '../hooks/useAsync';
import StateMessage from '../components/StateMessage/StateMessage';
import OrderStatusTracker from '../components/OrderStatusTracker/OrderStatusTracker';
import ProductImage from '../components/ProductImage/ProductImage';
import { formatDate, formatPrice } from '../utils/format';
import type { OrderWithItems, Product } from '../types';
import styles from './OrderDetail.module.css';

interface OrderDetailData {
  order: OrderWithItems;
  products: Map<string, Product>;
}

function OrderDetail() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const justPlaced = Boolean((location.state as { justPlaced?: boolean } | null)?.justPlaced);

  const { data, loading, error } = useAsync<OrderDetailData>(async () => {
    if (!id) throw new Error('Missing order id');
    const order = await getOrder(id);

    // The order API only returns product_id per line item — look products up
    // client-side to show real names/images instead of raw ids.
    const uniqueIds = [...new Set(order.items.map((item) => item.product_id))];
    const fetched = await Promise.all(uniqueIds.map((pid) => getProduct(pid).catch(() => null)));
    const products = new Map<string, Product>();
    fetched.forEach((product) => {
      if (product) products.set(product.id, product);
    });

    return { order, products };
  }, [id]);

  if (loading) {
    return <StateMessage icon="⏳" title="Loading order…" />;
  }

  if (error || !data) {
    return (
      <StateMessage icon="⚠️" title="Order not found" description="It may have been removed, or the link is incorrect." />
    );
  }

  const { order, products } = data;

  return (
    <section className={styles.page}>
      <nav className={styles.breadcrumb} aria-label="Breadcrumb">
        <Link to="/orders">Your Orders</Link>
        <span>/</span>
        <span className={styles.current}>#{order.id.slice(0, 8)}</span>
      </nav>

      {justPlaced && (
        <div className={styles.banner}>
          <span aria-hidden="true">🎉</span> Order placed successfully!
        </div>
      )}

      <div className={styles.header}>
        <div>
          <h1>Order #{order.id.slice(0, 8)}</h1>
          <p className={styles.placedOn}>Placed on {formatDate(order.created_at)}</p>
        </div>
        <p className={styles.total}>{formatPrice(order.total_amount)}</p>
      </div>

      <div className={styles.trackerCard}>
        <OrderStatusTracker status={order.status} />
      </div>

      <div className={styles.layout}>
        <ul className={styles.items}>
          {order.items.map((item) => {
            const product = products.get(item.product_id);
            return (
              <li key={item.id} className={styles.item}>
                <div className={styles.itemImage}>
                  <ProductImage src={product?.image_url} alt={product?.name ?? 'Product'} />
                </div>
                <div className={styles.itemInfo}>
                  <span className={styles.itemName}>{product?.name ?? 'Product'}</span>
                  <span className={styles.itemQty}>
                    {item.quantity} × {formatPrice(item.unit_price)}
                  </span>
                </div>
                <span className={styles.itemSubtotal}>{formatPrice(item.subtotal)}</span>
              </li>
            );
          })}
        </ul>

        <aside className={styles.delivery}>
          <h2>Delivery Details</h2>
          <div className={styles.deliveryRow}>
            <span className={styles.label}>Address</span>
            <span>{order.delivery_address}</span>
          </div>
          <div className={styles.deliveryRow}>
            <span className={styles.label}>Date</span>
            <span>{formatDate(order.delivery_date)}</span>
          </div>
        </aside>
      </div>
    </section>
  );
}

export default OrderDetail;
