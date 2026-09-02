import type { ReactNode } from 'react';
import styles from './StateMessage.module.css';

interface StateMessageProps {
  icon?: string;
  title: string;
  description?: string;
  children?: ReactNode;
}

function StateMessage({ icon, title, description, children }: StateMessageProps) {
  return (
    <div className={styles.state}>
      {icon && (
        <span className={styles.icon} aria-hidden="true">
          {icon}
        </span>
      )}
      <p className={styles.title}>{title}</p>
      {description && <p className={styles.description}>{description}</p>}
      {children}
    </div>
  );
}

export default StateMessage;
