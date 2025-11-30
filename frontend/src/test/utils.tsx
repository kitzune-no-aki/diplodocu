import { ReactElement, ReactNode } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { BrowserRouter, MemoryRouter } from 'react-router-dom';

interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  initialEntries?: string[];
}

// Wrapper with BrowserRouter for basic routing
const BrowserRouterWrapper = ({ children }: { children: ReactNode }) => {
  return <BrowserRouter>{children}</BrowserRouter>;
};

// Custom render with router support
export const renderWithRouter = (
  ui: ReactElement,
  options?: CustomRenderOptions
) => {
  const { initialEntries, ...renderOptions } = options || {};

  if (initialEntries) {
    return render(ui, {
      wrapper: ({ children }) => (
        <MemoryRouter initialEntries={initialEntries}>{children}</MemoryRouter>
      ),
      ...renderOptions,
    });
  }

  return render(ui, { wrapper: BrowserRouterWrapper, ...renderOptions });
};

export * from '@testing-library/react';
export { renderWithRouter as render };
