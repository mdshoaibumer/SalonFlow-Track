import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider } from '../providers/ThemeProvider'
import { AppRouter } from './AppRouter'

beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  })
})

function renderWithRouter(initialEntry = '/') {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <MemoryRouter initialEntries={[initialEntry]}>
          <AppRouter />
        </MemoryRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

describe('AppRouter', () => {
  it('renders dashboard on root route', () => {
    const { container } = renderWithRouter('/')
    expect(container.querySelector('main')).toBeInTheDocument()
  })

  it('renders staff page', () => {
    const { container } = renderWithRouter('/staff')
    expect(container.querySelector('main')).toBeInTheDocument()
  })

  it('renders services page', () => {
    const { container } = renderWithRouter('/services')
    expect(container.querySelector('main')).toBeInTheDocument()
  })
})
