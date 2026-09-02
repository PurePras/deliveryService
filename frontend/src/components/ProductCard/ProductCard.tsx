import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { Product } from '../../types';
import ProductImage from '../ProductImage/ProductImage';
import QuantityStepper from '../QuantityStepper/QuantityStepper';
import { useCart } from '../../store/cart/CartContext';
import { formatPricePerUnit } from '../../utils/format';
import styles from './ProductCard.module.css';

interface ProductCardProps {
  product: Product;
}

function ProductCard({ product }: ProductCardProps) {
  const { items, addItem, setQuantity } = useCart();
  const [justAdded, setJustAdded] = useState(false);

  const cartItem = items.find((item) => item.productId === product.id);

  function handleAdd() {
    addItem(product, 1);
    setJustAdded(true);
    setTimeout(() => setJustAdded(false), 1000);
  }

  return (
    <Link to={`/products/${product.id}`} className={styles.card}>
      <div className={styles.imageWrap}>
        <ProductImage src={product.image_url} alt={product.name} />
        {!product.is_available && <span className={styles.badge}>Out of stock</span>}
      </div>
      <div className={styles.body}>
        <h3 className={styles.name}>{product.name}</h3>
        <p className={styles.price}>{formatPricePerUnit(product.price, product.unit)}</p>

        {product.is_available && (
          // Stops the click from bubbling to the surrounding <Link> so cart controls
          // don't also trigger navigation to the product detail page.
          <div
            className={styles.action}
            onClick={(event) => {
              event.preventDefault();
              event.stopPropagation();
            }}
          >
            {cartItem ? (
              <QuantityStepper
                size="sm"
                quantity={cartItem.quantity}
                onIncrease={() => setQuantity(product.id, cartItem.quantity + 1)}
                onDecrease={() => setQuantity(product.id, cartItem.quantity - 1)}
              />
            ) : (
              <button type="button" className={styles.addBtn} onClick={handleAdd}>
                {justAdded ? 'Added ✓' : 'Add to Cart'}
              </button>
            )}
          </div>
        )}
      </div>
    </Link>
  );
}

export default ProductCard;
