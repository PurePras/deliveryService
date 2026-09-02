import { listCategories } from '../api/categories';
import { useAsync } from '../hooks/useAsync';
import CategoryCard from '../components/CategoryCard/CategoryCard';
import StateMessage from '../components/StateMessage/StateMessage';
import styles from './Categories.module.css';

function Categories() {
  const { data: categories, loading, error } = useAsync(() => listCategories({ is_active: true, limit: 100 }), []);

  return (
    <section className={styles.page}>
      <h1>Categories</h1>

      {loading && <StateMessage icon="⏳" title="Loading categories…" />}
      {error && (
        <StateMessage icon="⚠️" title="Couldn't load categories" description="Please try again in a moment." />
      )}
      {!loading && !error && categories && categories.length === 0 && (
        <StateMessage icon="🧺" title="No categories yet" />
      )}

      {!loading && !error && categories && categories.length > 0 && (
        <div className={styles.grid}>
          {categories.map((category) => (
            <CategoryCard key={category.id} category={category} />
          ))}
        </div>
      )}
    </section>
  );
}

export default Categories;
