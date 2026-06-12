import { describe, it, expect, vi } from 'vitest'
import { screen, within, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { Sidebar } from './Sidebar'

describe('Sidebar', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/')
  })

  it('renders app title', async () => {
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('SalonFlow')).toBeInTheDocument()
    })
  })

  it('renders all navigation groups', async () => {
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Main')).toBeInTheDocument()
    })
    expect(screen.getByText('Management')).toBeInTheDocument()
    expect(screen.getByText('Finance')).toBeInTheDocument()
    expect(screen.getByText('Inventory')).toBeInTheDocument()
    expect(screen.getByText('Reports')).toBeInTheDocument()
    expect(screen.getByText('System')).toBeInTheDocument()
  })

  it('expands active group by default (Main for /)', async () => {
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
    })
    expect(screen.getByText('Billing')).toBeInTheDocument()
  })

  it('shows Staff link when on /staff route', async () => {
    window.history.pushState({}, '', '/staff')
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Staff')).toBeInTheDocument()
    })
  })

  it('collapses group when clicking group label', async () => {
    const user = userEvent.setup()
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Main')).toBeInTheDocument()
    })
    const mainButton = screen.getByText('Main')
    expect(mainButton.closest('button')).toBeInTheDocument()
    await user.click(mainButton)
    expect(mainButton).toBeInTheDocument()
  })

  it('expands collapsed group on click', async () => {
    const user = userEvent.setup()
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Finance')).toBeInTheDocument()
    })
    await user.click(screen.getByText('Finance'))
    await waitFor(() => {
      expect(screen.getByText('Invoices')).toBeInTheDocument()
    })
    expect(screen.getByText('Salary')).toBeInTheDocument()
    expect(screen.getByText('Advances')).toBeInTheDocument()
    expect(screen.getByText('Expenses')).toBeInTheDocument()
  })

  it('nav links have correct href', async () => {
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
    })
    const dashLink = screen.getByText('Dashboard').closest('a')
    expect(dashLink).toHaveAttribute('href', '/')
  })

  it('highlights active nav item', async () => {
    window.history.pushState({}, '', '/')
    renderWithProviders(<Sidebar />)
    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
    })
    const dashLink = screen.getByText('Dashboard').closest('a')
    expect(dashLink).toHaveClass('bg-gradient-to-r')
  })
})
