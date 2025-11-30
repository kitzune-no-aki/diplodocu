import { vi } from 'vitest';

const keycloakMock = {
  init: vi.fn().mockResolvedValue(true),
  login: vi.fn(),
  logout: vi.fn(),
  register: vi.fn(),
  accountManagement: vi.fn(),
  createLoginUrl: vi.fn(),
  createLogoutUrl: vi.fn(),
  createRegisterUrl: vi.fn(),
  createAccountUrl: vi.fn(),
  isTokenExpired: vi.fn().mockReturnValue(false),
  updateToken: vi.fn().mockResolvedValue(true),
  clearToken: vi.fn(),
  hasRealmRole: vi.fn().mockReturnValue(false),
  hasResourceRole: vi.fn().mockReturnValue(false),
  loadUserProfile: vi.fn().mockResolvedValue({
    username: 'testuser',
    email: 'test@example.com',
    firstName: 'Test',
    lastName: 'User',
  }),
  loadUserInfo: vi.fn(),
  authenticated: false,
  token: 'mock-token',
  tokenParsed: {
    realm_access: { roles: [] },
  },
  refreshToken: 'mock-refresh-token',
  idToken: 'mock-id-token',
  onTokenExpired: undefined as (() => void) | undefined,
};

export default keycloakMock;
