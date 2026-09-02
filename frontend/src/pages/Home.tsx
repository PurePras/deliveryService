import { Link } from 'react-router-dom';
import { listCategories } from '../api/categories';
import { listProducts } from '../api/products';
import { useAsync } from '../hooks/useAsync';
import CategoryCard from '../components/CategoryCard/CategoryCard';
import ProductCard from '../components/ProductCard/ProductCard';
import StateMessage from '../components/StateMessage/StateMessage';
import styles from './Home.module.css';

function Home() {
  const { data: categories } = useAsync(() => listCategories({ is_active: true, limit: 6 }), []);
  const { data: products, loading: productsLoading } = useAsync(
    () => listProducts({ is_available: true, limit: 8 }),
    [],
  );

  return (
    <>
      <section className={styles.hero}>
        <div className={styles.badge}>🌿 Fresh from the farm, daily</div>
        <h1>Shri Ram Sabji Delivery</h1>
        <p className={styles.hindi}>ताज़ी सब्ज़ियाँ, अब आपके घर तक</p>
        <p className={styles.tagline}>Fresh vegetables delivered to your doorstep.</p>
        <Link to="/products" className={styles.cta}>
          Order Fresh Vegetables
        </Link>
      </section>

      {categories && categories.length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionHeader}>
            <h2>Shop by Category</h2>
            <Link to="/categories" className={styles.viewAll}>
              View all
            </Link>
          </div>
          <div className={styles.categoryStrip}>
            {categories.map((category) => (
              <CategoryCard key={category.id} category={category} />
            ))}
          </div>
        </section>
      )}

      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2>Fresh Picks</h2>
          <Link to="/products" className={styles.viewAll}>
            View all
          </Link>
        </div>
        {productsLoading && <StateMessage icon="⏳" title="Loading fresh picks…" />}
        {!productsLoading && products && products.length > 0 && (
          <div className={styles.productGrid}>
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </section>
    </>
  );
}

export default Home;
