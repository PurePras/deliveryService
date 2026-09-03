import { apiRequest } from './http';
import type { Order, OrderStatus, OrderWithItems } from '../types';

export interface CreateOrderItemInput {
  product_id: string;
  quantity: string;
}

export interface CreateOrderInput {
  delivery_area_id: string;
  delivery_slot_id: string;
  delivery_date: string;
  delivery_address: string;
  items: CreateOrderItemInput[];
}

export function createOrder(input: CreateOrderInput) {
  return apiRequest<OrderWithItems>('/orders', { method: 'POST', body: input });
}

export function getOrder(id: string) {
  return apiRequest<OrderWithItems>(`/orders/${id}`);
}

export function listOrders(params?: { status?: OrderStatus; limit?: number; offset?: number }) {
  return apiRequest<Order[]>('/orders', { query: params });
}

export function updateOrderStatus(id: string, status: OrderStatus) {
  return apiRequest<Order>(`/orders/${id}/status`, { method: 'PATCH', body: { status } });
}
