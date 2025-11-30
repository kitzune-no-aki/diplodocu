import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';

// We need to mock the modules before importing ProtectedRoute
vi.mock('../utils/AuthContext', () => ({
  useAuth: vi.fn(),
}));

// Mock Unauthorized component
vi.mock('./Unauthorized', () => ({
  default: () => <div data-testid="unauthorized">Unauthorized Access</div>,
}));

import ProtectedRoute from './protectedroute';
import { useAuth } from '../utils/AuthContext';

const mockUseAuth = useAuth as ReturnType<typeof vi.fn>;

const TestChild = () => <div data-testid="protected-content">Protected Content</div>;

const renderProtectedRoute = () => {
  return render(
    <BrowserRouter>
      <ProtectedRoute>
        <TestChild />
      </ProtectedRoute>
    </BrowserRouter>
  );
};

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders children when user is authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByTestId('protected-content')).toBeInTheDocument();
    expect(screen.getByText('Protected Content')).toBeInTheDocument();
    expect(screen.queryByTestId('unauthorized')).not.toBeInTheDocument();
  });

  it('renders Unauthorized component when user is not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByTestId('unauthorized')).toBeInTheDocument();
    expect(screen.getByText('Unauthorized Access')).toBeInTheDocument();
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
  });

  it('does not render children when not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });
});
