import { useState } from 'react';
import {
  listDeliverySlots,
  createDeliverySlot,
  updateDeliverySlot,
  deleteDeliverySlot,
  type DeliverySlotInput,
} from '../../api/deliverySlots';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import Modal from '../../components/admin/Modal/Modal';
import ConfirmDialog from '../../components/admin/ConfirmDialog/ConfirmDialog';
import { ApiError } from '../../api/http';
import type { DeliverySlot } from '../../types';
import styles from './AdminTable.module.css';

const emptyForm: DeliverySlotInput = { label: '', start_time: '', end_time: '', is_active: true };

function AdminDeliverySlots() {
  const [refetchKey, setRefetchKey] = useState(0);
  const { data: slots, loading, error } = useAsync(() => listDeliverySlots({ limit: 200 }), [refetchKey]);

  const [creating, setCreating] = useState(false);
  const [editing, setEditing] = useState<DeliverySlot | null>(null);
  const [form, setForm] = useState<DeliverySlotInput>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<DeliverySlot | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  function openCreate() {
    setForm(emptyForm);
    setFormError(null);
    setCreating(true);
  }

  function openEdit(slot: DeliverySlot) {
    setForm({
      label: slot.label,
      start_time: slot.start_time.slice(0, 5),
      end_time: slot.end_time.slice(0, 5),
      is_active: slot.is_active,
    });
    setFormError(null);
    setEditing(slot);
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
        await updateDeliverySlot(editing.id, form);
      } else {
        await createDeliverySlot(form);
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
      await deleteDeliverySlot(deleting.id);
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
        <h1>Delivery Slots</h1>
        <button type="button" className={styles.addBtn} onClick={openCreate}>
          + Add Slot
        </button>
      </div>

      {loading && <StateMessage icon="⏳" title="Loading delivery slots…" />}
      {error && <StateMessage icon="⚠️" title="Couldn't load delivery slots" />}

      {!loading && !error && slots && (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Label</th>
                <th>Start</th>
                <th>End</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {slots.map((slot) => (
                <tr key={slot.id}>
                  <td>{slot.label}</td>
                  <td>{slot.start_time.slice(0, 5)}</td>
                  <td>{slot.end_time.slice(0, 5)}</td>
                  <td>
                    <span className={`${styles.badge} ${slot.is_active ? styles.badgeActive : styles.badgeInactive}`}>
                      {slot.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td>
                    <div className={styles.actions}>
                      <button type="button" className={styles.actionBtn} onClick={() => openEdit(slot)}>
                        Edit
                      </button>
                      <button
                        type="button"
                        className={`${styles.actionBtn} ${styles.deleteBtn}`}
                        onClick={() => {
                          setDeleteError(null);
                          setDeleting(slot);
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {slots.length === 0 && (
                <tr>
                  <td colSpan={5}>No delivery slots yet.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {(creating || editing) && (
        <Modal title={editing ? 'Edit Delivery Slot' : 'Add Delivery Slot'} onClose={closeForm}>
          <form
            className={styles.form}
            onSubmit={(event) => {
              event.preventDefault();
              handleSubmit();
            }}
          >
            {formError && <p className={styles.error}>{formError}</p>}

            <label className={styles.field}>
              <span>Label</span>
              <input
                type="text"
                required
                placeholder="Morning (6 AM - 9 AM)"
                value={form.label}
                onChange={(event) => setForm({ ...form, label: event.target.value })}
              />
            </label>

            <div className={styles.formRow}>
              <label className={styles.field}>
                <span>Start time</span>
                <input
                  type="time"
                  required
                  value={form.start_time}
                  onChange={(event) => setForm({ ...form, start_time: event.target.value })}
                />
              </label>

              <label className={styles.field}>
                <span>End time</span>
                <input
                  type="time"
                  required
                  value={form.end_time}
                  onChange={(event) => setForm({ ...form, end_time: event.target.value })}
                />
              </label>
            </div>

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
          title="Delete Delivery Slot"
          message={`Are you sure you want to delete "${deleting.label}"? This can't be undone.`}
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

export default AdminDeliverySlots;
