import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Navbar from './navbar';

// Mock the useAuth hook
const mockLogin = vi.fn();
const mockLogout = vi.fn();

vi.mock('../utils/AuthContext', () => ({
  useAuth: () => ({
    isAuthenticated: false,
    login: mockLogin,
    logout: mockLogout,
  }),
}));

const renderNavbar = () => {
  return render(
    <BrowserRouter>
      <Navbar />
    </BrowserRouter>
  );
};

describe('Navbar', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders navigation links', () => {
    renderNavbar();

    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('My Table')).toBeInTheDocument();
  });

  it('renders Home link with correct path', () => {
    renderNavbar();

    const homeLink = screen.getByRole('link', { name: /home/i });
    expect(homeLink).toHaveAttribute('href', '/');
  });

  it('renders My Table link with correct path', () => {
    renderNavbar();

    const myTableLink = screen.getByRole('link', { name: /my table/i });
    expect(myTableLink).toHaveAttribute('href', '/mytable');
  });

  it('renders Login button when not authenticated', () => {
    renderNavbar();

    expect(screen.getByText('Login')).toBeInTheDocument();
    expect(screen.queryByText('Logout')).not.toBeInTheDocument();
  });

  it('calls login function when Login button is clicked', () => {
    renderNavbar();

    const loginButton = screen.getByText('Login');
    fireEvent.click(loginButton);

    expect(mockLogin).toHaveBeenCalledTimes(1);
  });
});

describe('Navbar - Authenticated', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Override mock for authenticated state
    vi.doMock('../utils/AuthContext', () => ({
      useAuth: () => ({
        isAuthenticated: true,
        login: mockLogin,
        logout: mockLogout,
      }),
    }));
  });

  it('renders Logout button when authenticated', async () => {
    // Re-import with new mock
    vi.resetModules();
    vi.doMock('../utils/AuthContext', () => ({
      useAuth: () => ({
        isAuthenticated: true,
        login: mockLogin,
        logout: mockLogout,
      }),
    }));

    const { default: NavbarAuth } = await import('./navbar');

    render(
      <BrowserRouter>
        <NavbarAuth />
      </BrowserRouter>
    );

    expect(screen.getByText('Logout')).toBeInTheDocument();
    expect(screen.queryByText('Login')).not.toBeInTheDocument();
  });
});
