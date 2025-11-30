import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Home from './home';

// Mock the useAuth hook for Navbar
vi.mock('../utils/AuthContext', () => ({
  useAuth: () => ({
    isAuthenticated: false,
    login: vi.fn(),
    logout: vi.fn(),
  }),
}));

const renderHome = () => {
  return render(
    <BrowserRouter>
      <Home />
    </BrowserRouter>
  );
};

describe('Home', () => {
  it('renders the Home component', () => {
    renderHome();

    // Check that the component renders without crashing
    expect(document.body).toBeInTheDocument();
  });

  it('displays greeting in multiple languages', () => {
    renderHome();

    // Check for various greetings
    expect(screen.getByText('Hallo')).toBeInTheDocument();
    expect(screen.getByText('Hello')).toBeInTheDocument();
    expect(screen.getByText('你好')).toBeInTheDocument();
    expect(screen.getByText('こんにちは')).toBeInTheDocument();
    expect(screen.getByText('Grüzi')).toBeInTheDocument();
    expect(screen.getByText('안녕하세요')).toBeInTheDocument();
    expect(screen.getByText('Salut')).toBeInTheDocument();
  });

  it('includes the Navbar component', () => {
    renderHome();

    // Navbar should be rendered (check for Home link in navbar)
    const homeLinks = screen.getAllByText('Home');
    expect(homeLinks.length).toBeGreaterThan(0);
  });
});
