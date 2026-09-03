import { useState } from 'react';
import { listCategories, createCategory, updateCategory, deleteCategory, type CategoryInput } from '../../api/categories';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import Modal from '../../components/admin/Modal/Modal';
import ConfirmDialog from '../../components/admin/ConfirmDialog/ConfirmDialog';
import { ApiError } from '../../api/http';
import type { Category } from '../../types';
import styles from './AdminTable.module.css';

const emptyForm: CategoryInput = { name: '', slug: '', description: '', image_url: '', is_active: true };

function AdminCategories() {
  const [refetchKey, setRefetchKey] = useState(0);
  const { data: categories, loading, error } = useAsync(() => listCategories({ limit: 200 }), [refetchKey]);

  const [creating, setCreating] = useState(false);
  const [editing, setEditing] = useState<Category | null>(null);
  const [form, setForm] = useState<CategoryInput>(emptyForm);
  const [formError, setFormError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<Category | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  function openCreate() {
    setForm(emptyForm);
    setFormError(null);
    setCreating(true);
  }

  function openEdit(category: Category) {
    setForm({
      name: category.name,
      slug: category.slug,
      description: category.description ?? '',
      image_url: category.image_url ?? '',
      is_active: category.is_active,
    });
    setFormError(null);
    setEditing(category);
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
        await updateCategory(editing.id, form);
      } else {
        await createCategory(form);
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
      await deleteCategory(deleting.id);
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
        <h1>Categories</h1>
        <button type="button" className={styles.addBtn} onClick={openCreate}>
          + Add Category
        </button>
      </div>

      {loading && <StateMessage icon="⏳" title="Loading categories…" />}
      {error && <StateMessage icon="⚠️" title="Couldn't load categories" />}

      {!loading && !error && categories && (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Name</th>
                <th>Slug</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {categories.map((category) => (
                <tr key={category.id}>
                  <td>{category.name}</td>
                  <td>{category.slug}</td>
                  <td>
                    <span className={`${styles.badge} ${category.is_active ? styles.badgeActive : styles.badgeInactive}`}>
                      {category.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td>
                    <div className={styles.actions}>
                      <button type="button" className={styles.actionBtn} onClick={() => openEdit(category)}>
                        Edit
                      </button>
                      <button
                        type="button"
                        className={`${styles.actionBtn} ${styles.deleteBtn}`}
                        onClick={() => {
                          setDeleteError(null);
                          setDeleting(category);
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {categories.length === 0 && (
                <tr>
                  <td colSpan={4}>No categories yet.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {(creating || editing) && (
        <Modal title={editing ? 'Edit Category' : 'Add Category'} onClose={closeForm}>
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
              <span>Slug</span>
              <input type="text" required value={form.slug} onChange={(event) => setForm({ ...form, slug: event.target.value })} />
            </label>

            <label className={styles.field}>
              <span>Description</span>
              <textarea
                rows={2}
                value={form.description}
                onChange={(event) => setForm({ ...form, description: event.target.value })}
              />
            </label>

            <label className={styles.field}>
              <span>Image URL</span>
              <input
                type="text"
                value={form.image_url}
                onChange={(event) => setForm({ ...form, image_url: event.target.value })}
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
          title="Delete Category"
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

export default AdminCategories;
