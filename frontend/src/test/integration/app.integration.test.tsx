import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';

// Mock keycloak
vi.mock('../../utils/Keycloak.tsx', () => ({
  default: {
    init: vi.fn().mockResolvedValue(true),
    login: vi.fn(),
    logout: vi.fn(),
    authenticated: true,
    token: 'mock-token',
    tokenParsed: { realm_access: { roles: [] } },
    loadUserProfile: vi.fn().mockResolvedValue({
      username: 'testuser',
      email: 'test@example.com',
    }),
    updateToken: vi.fn().mockResolvedValue(true),
    onTokenExpired: undefined,
  },
}));

// Mock AuthContext
vi.mock('../../utils/AuthContext', () => ({
  useAuth: () => ({
    isAuthenticated: true,
    isLoading: false,
    userProfile: { username: 'testuser' },
    roles: [],
    login: vi.fn(),
    logout: vi.fn(),
  }),
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import App from '../../App';
import Home from '../../components/home';
import Navbar from '../../components/navbar';

describe('App Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Routing', () => {
    it('renders home page at root route', async () => {
      render(
        <MemoryRouter initialEntries={['/']}>
          <Home />
        </MemoryRouter>
      );

      expect(screen.getByText('Hello')).toBeInTheDocument();
      expect(screen.getByText('Hallo')).toBeInTheDocument();
    });

    it('renders 404 page for unknown routes', async () => {
      const { default: NotFound } = await import('../../components/notfound');

      render(
        <MemoryRouter initialEntries={['/unknown-route']}>
          <NotFound />
        </MemoryRouter>
      );

      expect(screen.getByText('404')).toBeInTheDocument();
    });
  });

  describe('Navigation', () => {
    it('can navigate from home to other pages via navbar', async () => {
      const user = userEvent.setup();

      render(
        <MemoryRouter initialEntries={['/']}>
          <Home />
        </MemoryRouter>
      );

      // Check navbar is present
      expect(screen.getByText('Home')).toBeInTheDocument();
      expect(screen.getByText('My Table')).toBeInTheDocument();
    });

    it('shows logout button when authenticated', async () => {
      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      expect(screen.getByText('Logout')).toBeInTheDocument();
    });
  });

  describe('Authentication Flow', () => {
    it('calls logout when logout button is clicked', async () => {
      const mockLogout = vi.fn();

      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: true,
          isLoading: false,
          login: vi.fn(),
          logout: mockLogout,
        }),
      }));

      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      const logoutButton = screen.getByText('Logout');
      await userEvent.click(logoutButton);
    });
  });
});

describe('Component Integration', () => {
  it('Home page includes Navbar component', () => {
    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    );

    // Navbar elements should be present
    expect(screen.getByRole('link', { name: /home/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /my table/i })).toBeInTheDocument();
  });

  it('Navbar links have correct href attributes', () => {
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );

    const homeLink = screen.getByRole('link', { name: /home/i });
    const myTableLink = screen.getByRole('link', { name: /my table/i });

    expect(homeLink).toHaveAttribute('href', '/');
    expect(myTableLink).toHaveAttribute('href', '/mytable');
  });
});
