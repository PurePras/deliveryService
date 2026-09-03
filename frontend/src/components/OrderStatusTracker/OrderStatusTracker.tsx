import type { OrderStatus } from '../../types';
import styles from './OrderStatusTracker.module.css';

const STEPS: { key: OrderStatus; label: string }[] = [
  { key: 'pending', label: 'Order Placed' },
  { key: 'confirmed', label: 'Confirmed' },
  { key: 'out_for_delivery', label: 'Out for Delivery' },
  { key: 'delivered', label: 'Delivered' },
];

interface OrderStatusTrackerProps {
  status: OrderStatus;
}

function OrderStatusTracker({ status }: OrderStatusTrackerProps) {
  if (status === 'cancelled') {
    return (
      <div className={styles.cancelled}>
        <span aria-hidden="true">✕</span> Order Cancelled
      </div>
    );
  }

  const currentIndex = Math.max(
    STEPS.findIndex((step) => step.key === status),
    0,
  );

  return (
    <div className={styles.tracker}>
      {STEPS.map((step, index) => {
        const isDone = index < currentIndex;
        const isCurrent = index === currentIndex;
        const lineActive = index <= currentIndex;

        return (
          <div key={step.key} className={styles.step}>
            {index > 0 && <span className={`${styles.line} ${lineActive ? styles.lineActive : ''}`} />}
            <span className={`${styles.dot} ${isDone ? styles.dotDone : ''} ${isCurrent ? styles.dotCurrent : ''}`}>
              {isDone ? '✓' : index + 1}
            </span>
            <span className={`${styles.label} ${isCurrent || isDone ? styles.labelActive : ''}`}>{step.label}</span>
          </div>
        );
      })}
    </div>
  );
}

export default OrderStatusTracker;
