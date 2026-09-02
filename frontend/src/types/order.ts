export type OrderStatus = 'pending' | 'confirmed' | 'out_for_delivery' | 'delivered' | 'cancelled';

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  quantity: string;
  unit_price: string;
  subtotal: string;
  created_at: string;
}

export interface Order {
  id: string;
  user_id: string;
  delivery_area_id: string;
  delivery_slot_id: string;
  delivery_date: string;
  delivery_address: string;
  status: OrderStatus;
  total_amount: string;
  created_at: string;
  updated_at: string;
}

export interface OrderWithItems extends Order {
  items: OrderItem[];
}
