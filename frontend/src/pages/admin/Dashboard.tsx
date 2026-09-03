import { Link } from 'react-router-dom';
import { listOrders } from '../../api/orders';
import { listProducts } from '../../api/products';
import { listCategories } from '../../api/categories';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import { formatPrice } from '../../utils/format';
import styles from './Dashboard.module.css';

function AdminDashboard() {
  const { data: orders, loading: ordersLoading } = useAsync(() => listOrders({ limit: 200 }), []);
  const { data: products, loading: productsLoading } = useAsync(() => listProducts({ limit: 200 }), []);
  const { data: categories, loading: categoriesLoading } = useAsync(() => listCategories({ limit: 200 }), []);

  if (ordersLoading || productsLoading || categoriesLoading) {
    return <StateMessage icon="⏳" title="Loading dashboard…" />;
  }

  const totalOrders = orders?.length ?? 0;
  const pendingOrders = orders?.filter((order) => order.status === 'pending').length ?? 0;
  const totalRevenue = (orders ?? [])
    .filter((order) => order.status !== 'cancelled')
    .reduce((sum, order) => sum + Number(order.total_amount), 0);

  const stats = [
    { label: 'Total Orders', value: String(totalOrders), icon: '📦', to: '/admin/orders' },
    { label: 'Pending Orders', value: String(pendingOrders), icon: '⏳', to: '/admin/orders' },
    { label: 'Total Revenue', value: formatPrice(totalRevenue), icon: '💰' },
    { label: 'Products', value: String(products?.length ?? 0), icon: '🥕', to: '/admin/products' },
    { label: 'Categories', value: String(categories?.length ?? 0), icon: '🧺', to: '/admin/categories' },
  ];

  const recentOrders = (orders ?? []).slice(0, 5);

  return (
    <section className={styles.page}>
      <h1>Dashboard</h1>

      <div className={styles.grid}>
        {stats.map((stat) => {
          const content = (
            <>
              <span className={styles.statIcon} aria-hidden="true">
                {stat.icon}
              </span>
              <span className={styles.statValue}>{stat.value}</span>
              <span className={styles.statLabel}>{stat.label}</span>
            </>
          );
          return stat.to ? (
            <Link key={stat.label} to={stat.to} className={styles.statCard}>
              {content}
            </Link>
          ) : (
            <div key={stat.label} className={styles.statCard}>
              {content}
            </div>
          );
        })}
      </div>

      <div className={styles.recent}>
        <h2>Recent Orders</h2>
        {recentOrders.length === 0 ? (
          <p className={styles.empty}>No orders yet.</p>
        ) : (
          <ul className={styles.recentList}>
            {recentOrders.map((order) => (
              <li key={order.id}>
                <Link to={`/orders/${order.id}`} className={styles.recentItem}>
                  <span>#{order.id.slice(0, 8)}</span>
                  <span className={styles.recentStatus}>{order.status.replace(/_/g, ' ')}</span>
                  <span>{formatPrice(order.total_amount)}</span>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}

export default AdminDashboard;
