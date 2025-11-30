import { ReactNode } from 'react';
import { vi } from 'vitest';

interface MockAuthContextValue {
  isAuthenticated: boolean;
  isLoading: boolean;
  userProfile: { username?: string; email?: string } | null;
  roles: string[];
  login: () => void;
  logout: () => void;
}

const defaultMockAuth: MockAuthContextValue = {
  isAuthenticated: false,
  isLoading: false,
  userProfile: null,
  roles: [],
  login: vi.fn(),
  logout: vi.fn(),
};

export const createMockAuthContext = (overrides: Partial<MockAuthContextValue> = {}) => ({
  ...defaultMockAuth,
  ...overrides,
});

export const MockAuthProvider = ({
  children,
  value = defaultMockAuth
}: {
  children: ReactNode;
  value?: MockAuthContextValue
}) => {
  return <>{children}</>;
};

export { defaultMockAuth };
