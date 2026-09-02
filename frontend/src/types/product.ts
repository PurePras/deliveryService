export interface Product {
  id: string;
  category_id: string;
  name: string;
  slug: string;
  description?: string;
  unit: string;
  price: string;
  stock_quantity: string;
  image_url?: string;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}
