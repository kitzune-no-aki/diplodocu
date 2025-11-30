import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import NotFound from './notfound';

const renderNotFound = () => {
  return render(
    <BrowserRouter>
      <NotFound />
    </BrowserRouter>
  );
};

describe('NotFound', () => {
  it('renders 404 error code', () => {
    renderNotFound();

    expect(screen.getByText('404')).toBeInTheDocument();
  });

  it('displays error message', () => {
    renderNotFound();

    expect(screen.getByText(/the page you're looking for has vanished/i)).toBeInTheDocument();
  });

  it('renders a link back to home', () => {
    renderNotFound();

    const homeLink = screen.getByRole('link', { name: /return to safety/i });
    expect(homeLink).toBeInTheDocument();
    expect(homeLink).toHaveAttribute('href', '/');
  });

  it('has correct styling classes', () => {
    renderNotFound();

    const heading = screen.getByText('404');
    expect(heading.tagName).toBe('H1');
  });
});
