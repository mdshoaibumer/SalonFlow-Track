import { describe, it, expect, vi, beforeAll } from 'vitest'
import { screen, fireEvent } from '@testing-library/react'
import { renderWithProviders } from '@/test-utils'
import { Header } from './Header'

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

describe('Header', () => {
  it('renders Dashboard title on root route', () => {
    window.history.pushState({}, '', '/')
    renderWithProviders(<Header />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })

  it('renders Staff Management title on /staff route', () => {
    window.history.pushState({}, '', '/staff')
    renderWithProviders(<Header />)
    expect(screen.getByText('Staff Management')).toBeInTheDocument()
  })

  it('renders Settings title on /settings route', () => {
    window.history.pushState({}, '', '/settings')
    renderWithProviders(<Header />)
    expect(screen.getByText('Settings')).toBeInTheDocument()
  })

  it('renders theme toggle button', () => {
    window.history.pushState({}, '', '/')
    renderWithProviders(<Header />)
    expect(screen.getByRole('button', { name: /theme/i })).toBeInTheDocument()
  })

  it('renders fallback title for unknown routes', () => {
    window.history.pushState({}, '', '/unknown-path')
    renderWithProviders(<Header />)
    expect(screen.getByText('Unknown Path')).toBeInTheDocument()
  })

  it('cycles theme from light to dark on click', () => {
    window.history.pushState({}, '', '/')
    renderWithProviders(<Header />)
    const btn = screen.getByRole('button', { name: /theme/i })
    // light → dark
    fireEvent.click(btn)
    // dark → system
    fireEvent.click(btn)
    // system → light
    fireEvent.click(btn)
  })

  it('renders breadcrumbs for nested routes', () => {
    window.history.pushState({}, '', '/reports/revenue')
    renderWithProviders(<Header />)
    expect(screen.getByText('Revenue Reports')).toBeInTheDocument()
  })

  it('renders multi-segment breadcrumbs with separator', () => {
    window.history.pushState({}, '', '/reports/staff')
    renderWithProviders(<Header />)
    expect(screen.getByText('Staff Reports')).toBeInTheDocument()
    expect(screen.getByText('/')).toBeInTheDocument()
  })
})
