import { apiRequest } from './http';
import type { Category } from '../types';

export interface CategoryInput {
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
  is_active?: boolean;
}

export function listCategories(params?: { is_active?: boolean; limit?: number; offset?: number }) {
  return apiRequest<Category[]>('/categories', { query: params });
}

export function getCategory(id: string) {
  return apiRequest<Category>(`/categories/${id}`);
}

export function createCategory(input: CategoryInput) {
  return apiRequest<Category>('/categories', { method: 'POST', body: input });
}

export function updateCategory(id: string, input: CategoryInput) {
  return apiRequest<Category>(`/categories/${id}`, { method: 'PUT', body: input });
}

export function deleteCategory(id: string) {
  return apiRequest<void>(`/categories/${id}`, { method: 'DELETE' });
}
