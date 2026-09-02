import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getProduct } from '../api/products';
import { useAsync } from '../hooks/useAsync';
import ProductImage from '../components/ProductImage/ProductImage';
import StateMessage from '../components/StateMessage/StateMessage';
import QuantityStepper from '../components/QuantityStepper/QuantityStepper';
import { useCart } from '../store/cart/CartContext';
import { formatPrice, formatQuantity } from '../utils/format';
import styles from './ProductDetail.module.css';

function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const { items, addItem, setQuantity } = useCart();
  const [selectedQty, setSelectedQty] = useState(1);

  const {
    data: product,
    loading,
    error,
  } = useAsync(() => {
    if (!id) return Promise.reject(new Error('Missing product id'));
    return getProduct(id);
  }, [id]);

  if (loading) {
    return <StateMessage icon="⏳" title="Loading product…" />;
  }

  if (error || !product) {
    return (
      <StateMessage
        icon="⚠️"
        title="Product not found"
        description="It may have been removed, or the link is incorrect."
      />
    );
  }

  const cartItem = items.find((item) => item.productId === product.id);

  return (
    <section className={styles.page}>
      <nav className={styles.breadcrumb} aria-label="Breadcrumb">
        <Link to="/">Home</Link>
        <span>/</span>
        <Link to="/products">Products</Link>
        <span>/</span>
        <span className={styles.current}>{product.name}</span>
      </nav>

      <div className={styles.layout}>
        <div className={styles.imageWrap}>
          <ProductImage src={product.image_url} alt={product.name} />
          {!product.is_available && <span className={styles.badge}>Out of stock</span>}
        </div>

        <div className={styles.info}>
          <h1>{product.name}</h1>
          <p className={styles.price}>
            {formatPrice(product.price)} <span className={styles.unit}>/ {product.unit}</span>
          </p>
          {product.description && <p className={styles.description}>{product.description}</p>}
          <p className={styles.stock}>
            {product.is_available
              ? `${formatQuantity(product.stock_quantity)} ${product.unit} in stock`
              : 'Currently unavailable'}
          </p>

          {product.is_available && (
            <div className={styles.cartAction}>
              {cartItem ? (
                <>
                  <QuantityStepper
                    quantity={cartItem.quantity}
                    onIncrease={() => setQuantity(product.id, cartItem.quantity + 1)}
                    onDecrease={() => setQuantity(product.id, cartItem.quantity - 1)}
                  />
                  <Link to="/cart" className={styles.viewCart}>
                    View cart →
                  </Link>
                </>
              ) : (
                <>
                  <QuantityStepper
                    quantity={selectedQty}
                    onIncrease={() => setSelectedQty((q) => q + 1)}
                    onDecrease={() => setSelectedQty((q) => Math.max(1, q - 1))}
                  />
                  <button type="button" className={styles.addBtn} onClick={() => addItem(product, selectedQty)}>
                    Add to Cart
                  </button>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </section>
  );
}

export default ProductDetail;
