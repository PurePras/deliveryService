import { apiRequest } from './http';
import type { User } from '../types';

export function getUser(id: string) {
  return apiRequest<User>(`/users/${id}`);
}
