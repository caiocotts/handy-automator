import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { clearTokens, getAccessToken, login as apiLogin, LoginResponse, saveTokens } from '../lib/api';

interface AuthContextValue {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getAccessToken().then((token) => {
      setIsAuthenticated(token !== null);
      setIsLoading(false);
    });
  }, []);

  const login = useCallback(async (username: string, password: string) => {
    const data: LoginResponse = await apiLogin(username, password);
    await saveTokens(data.accessToken, data.refreshToken);
    setIsAuthenticated(true);
  }, []);

  const logout = useCallback(async () => {
    await clearTokens();
    setIsAuthenticated(false);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
