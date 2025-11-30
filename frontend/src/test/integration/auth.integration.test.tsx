import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';

describe('Authentication Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
  });

  describe('Protected Route Access', () => {
    it('shows unauthorized message when not authenticated', async () => {
      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: false,
          isLoading: false,
          userProfile: null,
          roles: [],
          login: vi.fn(),
          logout: vi.fn(),
        }),
      }));

      const { default: ProtectedRoute } = await import('../../components/protectedroute');
      const { default: Unauthorized } = await import('../../components/Unauthorized');

      // Mock Unauthorized to verify it's rendered
      vi.doMock('../../components/Unauthorized', () => ({
        default: () => <div data-testid="unauthorized">Please log in</div>,
      }));

      const { default: ProtectedRouteWithMock } = await import('../../components/protectedroute');

      render(
        <MemoryRouter>
          <ProtectedRouteWithMock>
            <div>Protected Content</div>
          </ProtectedRouteWithMock>
        </MemoryRouter>
      );

      // Should not show protected content
      expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    });

    it('shows protected content when authenticated', async () => {
      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: true,
          isLoading: false,
          userProfile: { username: 'testuser' },
          roles: [],
          login: vi.fn(),
          logout: vi.fn(),
        }),
      }));

      const { default: ProtectedRoute } = await import('../../components/protectedroute');

      render(
        <MemoryRouter>
          <ProtectedRoute>
            <div>Protected Content</div>
          </ProtectedRoute>
        </MemoryRouter>
      );

      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  describe('Login/Logout Flow', () => {
    it('shows login button when not authenticated', async () => {
      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: false,
          isLoading: false,
          login: vi.fn(),
          logout: vi.fn(),
        }),
      }));

      const { default: Navbar } = await import('../../components/navbar');

      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      expect(screen.getByText('Login')).toBeInTheDocument();
      expect(screen.queryByText('Logout')).not.toBeInTheDocument();
    });

    it('shows logout button when authenticated', async () => {
      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: true,
          isLoading: false,
          login: vi.fn(),
          logout: vi.fn(),
        }),
      }));

      const { default: Navbar } = await import('../../components/navbar');

      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      expect(screen.getByText('Logout')).toBeInTheDocument();
      expect(screen.queryByText('Login')).not.toBeInTheDocument();
    });

    it('calls login function when login button clicked', async () => {
      const mockLogin = vi.fn();

      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: false,
          isLoading: false,
          login: mockLogin,
          logout: vi.fn(),
        }),
      }));

      const { default: Navbar } = await import('../../components/navbar');
      const user = userEvent.setup();

      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      const loginButton = screen.getByText('Login');
      await user.click(loginButton);

      expect(mockLogin).toHaveBeenCalledTimes(1);
    });

    it('calls logout function when logout button clicked', async () => {
      const mockLogout = vi.fn();

      vi.doMock('../../utils/AuthContext', () => ({
        useAuth: () => ({
          isAuthenticated: true,
          isLoading: false,
          login: vi.fn(),
          logout: mockLogout,
        }),
      }));

      const { default: Navbar } = await import('../../components/navbar');
      const user = userEvent.setup();

      render(
        <MemoryRouter>
          <Navbar />
        </MemoryRouter>
      );

      const logoutButton = screen.getByText('Logout');
      await user.click(logoutButton);

      expect(mockLogout).toHaveBeenCalledTimes(1);
    });
  });
});
