import { apiRequest } from './http';

export interface HealthStatus {
  status: string;
  service: string;
  database: string;
}

export function getHealth() {
  return apiRequest<HealthStatus>('/health');
}
