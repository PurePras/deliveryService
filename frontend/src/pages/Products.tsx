import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { listProducts } from '../api/products';
import { listCategories } from '../api/categories';
import { useAsync } from '../hooks/useAsync';
import ProductCard from '../components/ProductCard/ProductCard';
import StateMessage from '../components/StateMessage/StateMessage';
import type { Product } from '../types';
import styles from './Products.module.css';

const PAGE_SIZE = 12;

function Products() {
  const [searchParams, setSearchParams] = useSearchParams();
  const categoryId = searchParams.get('category') ?? undefined;
  const urlQuery = searchParams.get('q') ?? '';

  const [searchInput, setSearchInput] = useState(urlQuery);
  const [products, setProducts] = useState<Product[]>([]);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);

  const { data: categories } = useAsync(() => listCategories({ is_active: true, limit: 100 }), []);

  useEffect(() => {
    const handle = setTimeout(() => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          if (searchInput) next.set('q', searchInput);
          else next.delete('q');
          return next;
        },
        { replace: true },
      );
    }, 300);
    return () => clearTimeout(handle);
  }, [searchInput, setSearchParams]);

  const { data, loading, error } = useAsync(
    () => listProducts({ category_id: categoryId, q: urlQuery || undefined, limit: PAGE_SIZE, offset: 0 }),
    [categoryId, urlQuery],
  );

  useEffect(() => {
    if (data) {
      setProducts(data);
      setOffset(data.length);
      setHasMore(data.length === PAGE_SIZE);
    }
  }, [data]);

  async function loadMore() {
    setLoadingMore(true);
    try {
      const more = await listProducts({
        category_id: categoryId,
        q: urlQuery || undefined,
        limit: PAGE_SIZE,
        offset,
      });
      setProducts((prev) => [...prev, ...more]);
      setOffset((prev) => prev + more.length);
      setHasMore(more.length === PAGE_SIZE);
    } finally {
      setLoadingMore(false);
    }
  }

  function selectCategory(id: string | undefined) {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (id) next.set('category', id);
      else next.delete('category');
      return next;
    });
  }

  return (
    <section className={styles.page}>
      <div className={styles.header}>
        <h1>Products</h1>
        <input
          type="search"
          placeholder="Search vegetables, fruits…"
          value={searchInput}
          onChange={(event) => setSearchInput(event.target.value)}
          className={styles.search}
        />
      </div>

      {categories && categories.length > 0 && (
        <div className={styles.filters}>
          <button
            type="button"
            className={`${styles.chip} ${!categoryId ? styles.chipActive : ''}`}
            onClick={() => selectCategory(undefined)}
          >
            All
          </button>
          {categories.map((category) => (
            <button
              key={category.id}
              type="button"
              className={`${styles.chip} ${categoryId === category.id ? styles.chipActive : ''}`}
              onClick={() => selectCategory(category.id)}
            >
              {category.name}
            </button>
          ))}
        </div>
      )}

      {loading && <StateMessage icon="⏳" title="Loading products…" />}
      {error && (
        <StateMessage icon="⚠️" title="Couldn't load products" description="Please try again in a moment." />
      )}
      {!loading && !error && products.length === 0 && (
        <StateMessage icon="🔍" title="No products found" description="Try a different search or category." />
      )}

      {!loading && !error && products.length > 0 && (
        <>
          <div className={styles.grid}>
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
          {hasMore && (
            <button type="button" className={styles.loadMore} onClick={loadMore} disabled={loadingMore}>
              {loadingMore ? 'Loading…' : 'Load more'}
            </button>
          )}
        </>
      )}
    </section>
  );
}

export default Products;
