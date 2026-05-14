export type PlatformRole = 'participant' | 'organizer';

export interface User {
  id: string;
  email: string;
  full_name: string;
  platform_role: PlatformRole;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  user: User;
}

export interface RegisterRequest {
  email: string;
  password: string;
  full_name: string;
  platform_role?: PlatformRole;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface UpdateProfileRequest {
  full_name: string;
}

export interface ApiErrorBody {
  error: {
    code: string;
    message: string;
  };
}
