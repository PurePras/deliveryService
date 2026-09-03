import { useState } from 'react';
import { listProducts, createProduct, updateProduct, deleteProduct, type ProductInput } from '../../api/products';
import { listCategories } from '../../api/categories';
import { useAsync } from '../../hooks/useAsync';
import StateMessage from '../../components/StateMessage/StateMessage';
import Modal from '../../components/admin/Modal/Modal';
import ConfirmDialog from '../../components/admin/ConfirmDialog/ConfirmDialog';
import { ApiError } from '../../api/http';
import { formatPricePerUnit } from '../../utils/format';
import type { Product } from '../../types';
import styles from './AdminTable.module.css';

const UNITS = ['kg', 'g', 'piece', 'bundle', 'dozen'];

function emptyForm(defaultCategoryId: string): ProductInput {
  return {
    category_id: defaultCategoryId,
    name: '',
    slug: '',
    description: '',
    unit: 'kg',
    price: '',
    stock_quantity: '',
    image_url: '',
    is_available: true,
  };
}

function AdminProducts() {
  const [refetchKey, setRefetchKey] = useState(0);
  const { data: products, loading, error } = useAsync(() => listProducts({ limit: 200 }), [refetchKey]);
  const { data: categories } = useAsync(() => listCategories({ limit: 200 }), []);

  const [creating, setCreating] = useState(false);
  const [editing, setEditing] = useState<Product | null>(null);
  const [form, setForm] = useState<ProductInput>(emptyForm(''));
  const [formError, setFormError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<Product | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  function categoryName(categoryId: string): string {
    return categories?.find((c) => c.id === categoryId)?.name ?? '—';
  }

  function openCreate() {
    setForm(emptyForm(categories?.[0]?.id ?? ''));
    setFormError(null);
    setCreating(true);
  }

  function openEdit(product: Product) {
    setForm({
      category_id: product.category_id,
      name: product.name,
      slug: product.slug,
      description: product.description ?? '',
      unit: product.unit,
      price: product.price,
      stock_quantity: product.stock_quantity,
      image_url: product.image_url ?? '',
      is_available: product.is_available,
    });
    setFormError(null);
    setEditing(product);
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
        await updateProduct(editing.id, form);
      } else {
        await createProduct(form);
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
      await deleteProduct(deleting.id);
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
        <h1>Products</h1>
        <button type="button" className={styles.addBtn} onClick={openCreate}>
          + Add Product
        </button>
      </div>

      {loading && <StateMessage icon="⏳" title="Loading products…" />}
      {error && <StateMessage icon="⚠️" title="Couldn't load products" />}

      {!loading && !error && products && (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Name</th>
                <th>Category</th>
                <th>Price</th>
                <th>Stock</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id}>
                  <td>{product.name}</td>
                  <td>{categoryName(product.category_id)}</td>
                  <td>{formatPricePerUnit(product.price, product.unit)}</td>
                  <td>{product.stock_quantity}</td>
                  <td>
                    <span className={`${styles.badge} ${product.is_available ? styles.badgeActive : styles.badgeInactive}`}>
                      {product.is_available ? 'Available' : 'Unavailable'}
                    </span>
                  </td>
                  <td>
                    <div className={styles.actions}>
                      <button type="button" className={styles.actionBtn} onClick={() => openEdit(product)}>
                        Edit
                      </button>
                      <button
                        type="button"
                        className={`${styles.actionBtn} ${styles.deleteBtn}`}
                        onClick={() => {
                          setDeleteError(null);
                          setDeleting(product);
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {products.length === 0 && (
                <tr>
                  <td colSpan={6}>No products yet.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {(creating || editing) && (
        <Modal title={editing ? 'Edit Product' : 'Add Product'} onClose={closeForm}>
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
              <span>Category</span>
              <select
                required
                value={form.category_id}
                onChange={(event) => setForm({ ...form, category_id: event.target.value })}
              >
                <option value="" disabled>
                  Select a category
                </option>
                {categories?.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </label>

            <label className={styles.field}>
              <span>Description</span>
              <textarea
                rows={2}
                value={form.description}
                onChange={(event) => setForm({ ...form, description: event.target.value })}
              />
            </label>

            <div className={styles.formRow}>
              <label className={styles.field}>
                <span>Unit</span>
                <select value={form.unit} onChange={(event) => setForm({ ...form, unit: event.target.value })}>
                  {UNITS.map((unit) => (
                    <option key={unit} value={unit}>
                      {unit}
                    </option>
                  ))}
                </select>
              </label>

              <label className={styles.field}>
                <span>Price (₹)</span>
                <input
                  type="text"
                  required
                  inputMode="decimal"
                  value={form.price}
                  onChange={(event) => setForm({ ...form, price: event.target.value })}
                />
              </label>
            </div>

            <div className={styles.formRow}>
              <label className={styles.field}>
                <span>Stock quantity</span>
                <input
                  type="text"
                  required
                  inputMode="decimal"
                  value={form.stock_quantity}
                  onChange={(event) => setForm({ ...form, stock_quantity: event.target.value })}
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
            </div>

            <label className={`${styles.field} ${styles.checkboxField}`}>
              <input
                type="checkbox"
                checked={!!form.is_available}
                onChange={(event) => setForm({ ...form, is_available: event.target.checked })}
              />
              <span>Available</span>
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
          title="Delete Product"
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

export default AdminProducts;
