import { useState } from 'react';
import {
  listDeliveryAreas,
  createDeliveryArea,
  updateDeliveryArea,
  deleteDeliveryArea,
  type DeliveryAreaInput,
} from '../../api/deliveryAreas';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import Modal from '../../components/admin/Modal/Modal';
import ConfirmDialog from '../../components/admin/ConfirmDialog/ConfirmDialog';
import { ApiError } from '../../api/http';
import type { DeliveryArea } from '../../types';
import styles from './AdminTable.module.css';

const emptyForm: DeliveryAreaInput = { name: '', city: '', pincode: '', is_active: true };

function AdminDeliveryAreas() {
  const [refetchKey, setRefetchKey] = useState(0);
  const { data: areas, loading, error } = useAsync(() => listDeliveryAreas({ limit: 200 }), [refetchKey]);

  const [creating, setCreating] = useState(false);
  const [editing, setEditing] = useState<DeliveryArea | null>(null);
  const [form, setForm] = useState<DeliveryAreaInput>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<DeliveryArea | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  function openCreate() {
    setForm(emptyForm);
    setFormError(null);
    setCreating(true);
  }

  function openEdit(area: DeliveryArea) {
    setForm({ name: area.name, city: area.city, pincode: area.pincode, is_active: area.is_active });
    setFormError(null);
    setEditing(area);
  }

  function closeForm() {
    setCreating(false);
    setEditing(null);
  }

  async function handleSubmit() {
    setFormError(null);
    setSubmitting(true);
    try {
      if (editing) {
        await updateDeliveryArea(editing.id, form);
      } else {
        await createDeliveryArea(form);
      }
      closeForm();
      setRefetchKey((k) => k + 1);
    } catch (err) {
      setFormError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete() {
    if (!deleting) return;
    setDeleteError(null);
    setSubmitting(true);
    try {
      await deleteDeliveryArea(deleting.id);
      setDeleting(null);
      setRefetchKey((k) => k + 1);
    } catch (err) {
      setDeleteError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className={styles.page}>
      <div className={styles.header}>
        <h1>Delivery Areas</h1>
        <button type="button" className={styles.addBtn} onClick={openCreate}>
          + Add Area
        </button>
      </div>

      {loading && <StateMessage icon="⏳" title="Loading delivery areas…" />}
      {error && <StateMessage icon="⚠️" title="Couldn't load delivery areas" />}

      {!loading && !error && areas && (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Name</th>
                <th>City</th>
                <th>Pincode</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {areas.map((area) => (
                <tr key={area.id}>
                  <td>{area.name}</td>
                  <td>{area.city}</td>
                  <td>{area.pincode}</td>
                  <td>
                    <span className={`${styles.badge} ${area.is_active ? styles.badgeActive : styles.badgeInactive}`}>
                      {area.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td>
                    <div className={styles.actions}>
                      <button type="button" className={styles.actionBtn} onClick={() => openEdit(area)}>
                        Edit
                      </button>
                      <button
                        type="button"
                        className={`${styles.actionBtn} ${styles.deleteBtn}`}
                        onClick={() => {
                          setDeleteError(null);
                          setDeleting(area);
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {areas.length === 0 && (
                <tr>
                  <td colSpan={5}>No delivery areas yet.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {(creating || editing) && (
        <Modal title={editing ? 'Edit Delivery Area' : 'Add Delivery Area'} onClose={closeForm}>
          <form
            className={styles.form}
            onSubmit={(event) => {
              event.preventDefault();
              handleSubmit();
            }}
          >
            {formError && <p className={styles.error}>{formError}</p>}

            <label className={styles.field}>
              <span>Name</span>
              <input type="text" required value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
            </label>

            <label className={styles.field}>
              <span>City</span>
              <input type="text" required value={form.city} onChange={(event) => setForm({ ...form, city: event.target.value })} />
            </label>

            <label className={styles.field}>
              <span>Pincode</span>
              <input
                type="text"
                required
                value={form.pincode}
                onChange={(event) => setForm({ ...form, pincode: event.target.value })}
              />
            </label>

            <label className={`${styles.field} ${styles.checkboxField}`}>
              <input
                type="checkbox"
                checked={!!form.is_active}
                onChange={(event) => setForm({ ...form, is_active: event.target.checked })}
              />
              <span>Active</span>
            </label>

            <div className={styles.formActions}>
              <button type="button" className={styles.cancelBtn} onClick={closeForm} disabled={submitting}>
                Cancel
              </button>
              <button type="submit" className={styles.submitBtn} disabled={submitting}>
                {submitting ? 'Saving…' : 'Save'}
              </button>
            </div>
          </form>
        </Modal>
      )}

      {deleting && (
        <ConfirmDialog
          title="Delete Delivery Area"
          message={`Are you sure you want to delete "${deleting.name}"? This can't be undone.`}
          confirmLabel="Delete"
          danger
          busy={submitting}
          error={deleteError}
          onConfirm={handleDelete}
          onCancel={() => setDeleting(null)}
        />
      )}
    </section>
  );
}

export default AdminDeliveryAreas;
