export type UserRole = 'customer' | 'admin';

export interface User {
  id: string;
  name: string;
  phone: string;
  email?: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}
