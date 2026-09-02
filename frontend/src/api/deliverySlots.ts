import { apiRequest } from './http';
import type { DeliverySlot } from '../types';

export interface DeliverySlotInput {
  label: string;
  start_time: string;
  end_time: string;
  is_active?: boolean;
}

export function listDeliverySlots(params?: { is_active?: boolean; limit?: number; offset?: number }) {
  return apiRequest<DeliverySlot[]>('/delivery-slots', { query: params });
}

export function getDeliverySlot(id: string) {
  return apiRequest<DeliverySlot>(`/delivery-slots/${id}`);
}

export function createDeliverySlot(input: DeliverySlotInput) {
  return apiRequest<DeliverySlot>('/delivery-slots', { method: 'POST', body: input });
}

export function updateDeliverySlot(id: string, input: DeliverySlotInput) {
  return apiRequest<DeliverySlot>(`/delivery-slots/${id}`, { method: 'PUT', body: input });
}

export function deleteDeliverySlot(id: string) {
  return apiRequest<void>(`/delivery-slots/${id}`, { method: 'DELETE' });
}
