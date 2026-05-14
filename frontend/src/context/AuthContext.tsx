import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';

import * as authApi from '../api/auth';
import { clearToken, getToken, setToken } from '../api/client';
import type { LoginRequest, RegisterRequest, User } from '../api/types';

interface AuthState {
  user: User | null;
  loading: boolean;
  login: (body: LoginRequest) => Promise<void>;
  register: (body: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    const token = getToken();
    if (!token) {
      setUser(null);
      return;
    }
    const me = await authApi.getMe();
    setUser(me);
  }, []);

  useEffect(() => {
    (async () => {
      try {
        if (getToken()) {
          await refreshUser();
        }
      } catch {
        clearToken();
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
  }, [refreshUser]);

  const login = useCallback(async (body: LoginRequest) => {
    const resp = await authApi.login(body);
    setToken(resp.access_token);
    setUser(resp.user);
  }, []);

  const register = useCallback(async (body: RegisterRequest) => {
    const resp = await authApi.register(body);
    setToken(resp.access_token);
    setUser(resp.user);
  }, []);

  const logout = useCallback(() => {
    clearToken();
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({ user, loading, login, register, logout, refreshUser }),
    [user, loading, login, register, logout, refreshUser],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth outside AuthProvider');
  }
  return ctx;
}
