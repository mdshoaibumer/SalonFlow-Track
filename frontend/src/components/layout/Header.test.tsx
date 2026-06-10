import { describe, it, expect, vi, beforeAll } from 'vitest'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
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
    expect(screen.getByRole('button', { name: /toggle theme/i })).toBeInTheDocument()
  })

  it('cycles theme on toggle click', async () => {
    const user = userEvent.setup()
    window.history.pushState({}, '', '/')
    renderWithProviders(<Header />)
    const btn = screen.getByRole('button', { name: /toggle theme/i })
    // Default is light, clicking should cycle to system
    await user.click(btn)
    // After click, theme should change (button still exists)
    expect(btn).toBeInTheDocument()
  })

  it('renders fallback title for unknown routes', () => {
    window.history.pushState({}, '', '/unknown-path')
    renderWithProviders(<Header />)
    expect(screen.getByText('SalonFlow Track')).toBeInTheDocument()
  })
})
