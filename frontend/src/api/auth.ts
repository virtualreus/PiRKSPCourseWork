import { api } from './client';
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
  User,
} from './types';

export function register(body: RegisterRequest): Promise<AuthResponse> {
  return api<AuthResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function login(body: LoginRequest): Promise<AuthResponse> {
  return api<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function getMe(): Promise<User> {
  return api<User>('/users/me');
}

export function updateMe(body: UpdateProfileRequest): Promise<User> {
  return api<User>('/users/me', {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}
