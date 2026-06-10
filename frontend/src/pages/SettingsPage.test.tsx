import { describe, it, expect, vi, beforeAll } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { SettingsPage } from './SettingsPage'

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

describe('SettingsPage', () => {
  it('renders page header', async () => {
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument()
    })
  })

  it('renders tab navigation', async () => {
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('General')).toBeInTheDocument()
    })
    expect(screen.getByText('Salon')).toBeInTheDocument()
    expect(screen.getByText('Notifications')).toBeInTheDocument()
    expect(screen.getByText('System')).toBeInTheDocument()
  })

  it('shows Appearance section in General tab by default', async () => {
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('Appearance')).toBeInTheDocument()
    })
  })

  it('shows theme selector', async () => {
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('Theme')).toBeInTheDocument()
    })
  })

  it('switches to Salon tab on click', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('Salon')).toBeInTheDocument()
    })
    await user.click(screen.getByText('Salon'))
    await waitFor(() => {
      expect(screen.getByText('Salon Details')).toBeInTheDocument()
    })
  })

  it('shows salon name input with value from API', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('Salon')).toBeInTheDocument()
    })
    await user.click(screen.getByText('Salon'))
    await waitFor(() => {
      expect(screen.getByDisplayValue('SalonFlow Studio')).toBeInTheDocument()
    })
  })

  it('switches to System tab on click', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('System')).toBeInTheDocument()
    })
    await user.click(screen.getByText('System'))
    await waitFor(() => {
      expect(screen.getByText('System Information')).toBeInTheDocument()
    })
  })

  it('renders health status in System tab', async () => {
    const user = userEvent.setup()
    renderWithProviders(<SettingsPage />)
    await waitFor(() => {
      expect(screen.getByText('System')).toBeInTheDocument()
    })
    await user.click(screen.getByText('System'))
    await waitFor(() => {
      expect(screen.getByText('healthy')).toBeInTheDocument()
    })
  })
})
