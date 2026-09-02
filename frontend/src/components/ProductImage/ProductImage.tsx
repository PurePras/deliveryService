import styles from './ProductImage.module.css';

interface ProductImageProps {
  src?: string | null;
  alt: string;
  className?: string;
}

function ProductImage({ src, alt, className = '' }: ProductImageProps) {
  if (src) {
    return <img src={src} alt={alt} loading="lazy" className={`${styles.image} ${className}`} />;
  }

  return (
    <div className={`${styles.placeholder} ${className}`} role="img" aria-label={alt}>
      <span aria-hidden="true">🥬</span>
    </div>
  );
}

export default ProductImage;
