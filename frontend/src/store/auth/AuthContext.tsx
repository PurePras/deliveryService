import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import type { User } from '../../types';
import {
  getMe,
  logout as apiLogout,
  register as apiRegister,
  requestLoginOtp,
  verifyLoginOtp,
  type RegisterInput,
} from '../../api/auth';

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  requestOtp: (phone: string) => Promise<void>;
  verifyOtp: (phone: string, code: string) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // The session lives in an httpOnly cookie the frontend can't read directly,
  // so /auth/me is the only way to discover "am I already logged in" on load.
  useEffect(() => {
    let cancelled = false;

    getMe()
      .then((me) => {
        if (!cancelled) setUser(me);
      })
      .catch(() => {
        if (!cancelled) setUser(null);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  async function requestOtp(phone: string) {
    await requestLoginOtp({ phone });
  }

  async function verifyOtp(phone: string, code: string) {
    setUser(await verifyLoginOtp({ phone, code }));
  }

  async function register(input: RegisterInput) {
    setUser(await apiRegister(input));
  }

  async function logout() {
    await apiLogout();
    setUser(null);
  }

  return (
    <AuthContext.Provider value={{ user, loading, requestOtp, verifyOtp, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within an AuthProvider');
  return ctx;
}
