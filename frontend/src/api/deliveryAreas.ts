import { apiRequest } from './http';
import type { DeliveryArea } from '../types';

export interface DeliveryAreaInput {
  name: string;
  city: string;
  pincode: string;
  is_active?: boolean;
}

export function listDeliveryAreas(params?: { is_active?: boolean; limit?: number; offset?: number }) {
  return apiRequest<DeliveryArea[]>('/delivery-areas', { query: params });
}

export function getDeliveryArea(id: string) {
  return apiRequest<DeliveryArea>(`/delivery-areas/${id}`);
}

export function createDeliveryArea(input: DeliveryAreaInput) {
  return apiRequest<DeliveryArea>('/delivery-areas', { method: 'POST', body: input });
}

export function updateDeliveryArea(id: string, input: DeliveryAreaInput) {
  return apiRequest<DeliveryArea>(`/delivery-areas/${id}`, { method: 'PUT', body: input });
}

export function deleteDeliveryArea(id: string) {
  return apiRequest<void>(`/delivery-areas/${id}`, { method: 'DELETE' });
}
