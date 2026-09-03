import type { OrderStatus } from '../../types';
import styles from './OrderStatusBadge.module.css';

export const ORDER_STATUS_LABELS: Record<OrderStatus, string> = {
  pending: 'Order Placed',
  confirmed: 'Confirmed',
  out_for_delivery: 'Out for Delivery',
  delivered: 'Delivered',
  cancelled: 'Cancelled',
};

interface OrderStatusBadgeProps {
  status: OrderStatus;
}

function OrderStatusBadge({ status }: OrderStatusBadgeProps) {
  return <span className={`${styles.badge} ${styles[status]}`}>{ORDER_STATUS_LABELS[status]}</span>;
}

export default OrderStatusBadge;
