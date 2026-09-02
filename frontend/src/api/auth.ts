import { apiRequest } from './http';
import type { User } from '../types';

export interface RegisterInput {
  name: string;
  phone: string;
  password: string;
  email?: string;
}

export interface LoginInput {
  phone: string;
  password: string;
}

export function register(input: RegisterInput) {
  return apiRequest<User>('/auth/register', { method: 'POST', body: input });
}

export function login(input: LoginInput) {
  return apiRequest<User>('/auth/login', { method: 'POST', body: input });
}

export function logout() {
  return apiRequest<{ status: string }>('/auth/logout', { method: 'POST' });
}

export function getMe() {
  return apiRequest<User>('/auth/me');
}
