import { useState } from 'react';
import { Link } from 'react-router-dom';
import { listOrders, updateOrderStatus } from '../../api/orders';
import { getUser } from '../../api/users';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import { ORDER_STATUS_LABELS } from '../../components/OrderStatusBadge/OrderStatusBadge';
import { ApiError } from '../../api/http';
import { formatDate, formatPrice } from '../../utils/format';
import type { OrderStatus, User } from '../../types';
import styles from './AdminTable.module.css';

const STATUS_OPTIONS: OrderStatus[] = ['pending', 'confirmed', 'out_for_delivery', 'delivered', 'cancelled'];

function AdminOrders() {
  const [statusFilter, setStatusFilter] = useState<OrderStatus | ''>('');
  const [refetchKey, setRefetchKey] = useState(0);
  const {
    data: orders,
    loading,
    error,
  } = useAsync(() => listOrders({ status: statusFilter || undefined, limit: 200 }), [statusFilter, refetchKey]);

  const { data: userMap } = useAsync(async () => {
    const uniqueUserIds = [...new Set((orders ?? []).map((order) => order.user_id))];
    const users = await Promise.all(uniqueUserIds.map((id) => getUser(id).catch(() => null)));
    const map = new Map<string, User>();
    users.forEach((user) => {
      if (user) map.set(user.id, user);
    });
    return map;
  }, [orders]);

  const [updatingId, setUpdatingId] = useState<string | null>(null);
  const [rowError, setRowError] = useState<{ id: string; message: string } | null>(null);

  async function handleStatusChange(orderId: string, status: OrderStatus) {
    setRowError(null);
    setUpdatingId(orderId);
    try {
      await updateOrderStatus(orderId, status);
      setRefetchKey((k) => k + 1);
    } catch (err) {
      setRowError({ id: orderId, message: err instanceof ApiError ? err.message : 'Failed to update status.' });
    } finally {
      setUpdatingId(null);
    }
  }

  return (
    <section className={styles.page}>
      <div className={styles.header}>
        <h1>Orders</h1>
        <select
          className={styles.filterSelect}
          value={statusFilter}
          onChange={(event) => setStatusFilter(event.target.value as OrderStatus | '')}
        >
          <option value="">All statuses</option>
          {STATUS_OPTIONS.map((status) => (
            <option key={status} value={status}>
              {ORDER_STATUS_LABELS[status]}
            </option>
          ))}
        </select>
      </div>

      {loading && <StateMessage icon="⏳" title="Loading orders…" />}
      {error && <StateMessage icon="⚠️" title="Couldn't load orders" />}

      {!loading && !error && orders && (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Order</th>
                <th>Customer</th>
                <th>Date</th>
                <th>Total</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr key={order.id}>
                  <td>#{order.id.slice(0, 8)}</td>
                  <td>{userMap?.get(order.user_id)?.name ?? '—'}</td>
                  <td>{formatDate(order.created_at)}</td>
                  <td>{formatPrice(order.total_amount)}</td>
                  <td>
                    <select
                      className={styles.filterSelect}
                      value={order.status}
                      disabled={updatingId === order.id}
                      onChange={(event) => handleStatusChange(order.id, event.target.value as OrderStatus)}
                    >
                      {STATUS_OPTIONS.map((status) => (
                        <option key={status} value={status}>
                          {ORDER_STATUS_LABELS[status]}
                        </option>
                      ))}
                    </select>
                    {rowError?.id === order.id && <p className={styles.error}>{rowError.message}</p>}
                  </td>
                  <td>
                    <Link to={`/orders/${order.id}`} className={styles.actionBtn}>
                      View
                    </Link>
                  </td>
                </tr>
              ))}
              {orders.length === 0 && (
                <tr>
                  <td colSpan={6}>No orders found.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

export default AdminOrders;
