import { apiRequest } from './http';
import type { User } from '../types';

export interface RegisterUserInput {
  name: string;
  phone: string;
  email?: string;
}

export function registerUser(input: RegisterUserInput) {
  return apiRequest<User>('/users', { method: 'POST', body: input });
}

export function getUser(id: string) {
  return apiRequest<User>(`/users/${id}`);
}
