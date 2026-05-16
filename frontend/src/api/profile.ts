import { api } from './client';
import type { UserDashboard } from './profileTypes';

export function getDashboard(): Promise<UserDashboard> {
  return api<UserDashboard>('/users/me/dashboard');
}
