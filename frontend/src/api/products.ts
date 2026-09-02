import { apiRequest } from './http';
import type { Product } from '../types';

export interface ProductInput {
  category_id: string;
  name: string;
  slug: string;
  description?: string;
  unit: string;
  price: string;
  stock_quantity: string;
  image_url?: string;
  is_available?: boolean;
}

export function listProducts(params?: {
  category_id?: string;
  is_available?: boolean;
  q?: string;
  limit?: number;
  offset?: number;
}) {
  return apiRequest<Product[]>('/products', { query: params });
}

export function getProduct(id: string) {
  return apiRequest<Product>(`/products/${id}`);
}

export function createProduct(input: ProductInput) {
  return apiRequest<Product>('/products', { method: 'POST', body: input });
}

export function updateProduct(id: string, input: ProductInput) {
  return apiRequest<Product>(`/products/${id}`, { method: 'PUT', body: input });
}

export function deleteProduct(id: string) {
  return apiRequest<void>(`/products/${id}`, { method: 'DELETE' });
}
