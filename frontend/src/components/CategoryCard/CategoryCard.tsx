import { Link } from 'react-router-dom';
import type { Category } from '../../types';
import ProductImage from '../ProductImage/ProductImage';
import styles from './CategoryCard.module.css';

interface CategoryCardProps {
  category: Category;
}

function CategoryCard({ category }: CategoryCardProps) {
  return (
    <Link to={`/products?category=${category.id}`} className={styles.card}>
      <div className={styles.imageWrap}>
        <ProductImage src={category.image_url} alt={category.name} />
      </div>
      <span className={styles.name}>{category.name}</span>
    </Link>
  );
}

export default CategoryCard;
