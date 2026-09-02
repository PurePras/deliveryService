import styles from './QuantityStepper.module.css';

interface QuantityStepperProps {
  quantity: number;
  onIncrease: () => void;
  onDecrease: () => void;
  size?: 'sm' | 'md';
}

function QuantityStepper({ quantity, onIncrease, onDecrease, size = 'md' }: QuantityStepperProps) {
  return (
    <div className={`${styles.stepper} ${size === 'sm' ? styles.sm : ''}`}>
      <button type="button" className={styles.btn} onClick={onDecrease} aria-label="Decrease quantity">
        −
      </button>
      <span className={styles.value}>{quantity}</span>
      <button type="button" className={styles.btn} onClick={onIncrease} aria-label="Increase quantity">
        +
      </button>
    </div>
  );
}

export default QuantityStepper;
