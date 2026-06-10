import { describe, it, expect, vi } from 'vitest'
import { screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { Sidebar } from './Sidebar'

describe('Sidebar', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/')
  })

  it('renders app title', () => {
    renderWithProviders(<Sidebar />)
    expect(screen.getByText('SalonFlow')).toBeInTheDocument()
  })

  it('renders all navigation groups', () => {
    renderWithProviders(<Sidebar />)
    expect(screen.getByText('Main')).toBeInTheDocument()
    expect(screen.getByText('Management')).toBeInTheDocument()
    expect(screen.getByText('Finance')).toBeInTheDocument()
    expect(screen.getByText('Inventory')).toBeInTheDocument()
    expect(screen.getByText('Reports')).toBeInTheDocument()
    expect(screen.getByText('System')).toBeInTheDocument()
  })

  it('expands active group by default (Main for /)', () => {
    renderWithProviders(<Sidebar />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Billing')).toBeInTheDocument()
  })

  it('shows Staff link when on /staff route', () => {
    window.history.pushState({}, '', '/staff')
    renderWithProviders(<Sidebar />)
    expect(screen.getByText('Staff')).toBeInTheDocument()
  })

  it('collapses group when clicking group label', async () => {
    const user = userEvent.setup()
    renderWithProviders(<Sidebar />)
    // Main group is expanded, click to collapse
    await user.click(screen.getByText('Main'))
    // After collapse, Dashboard link should be hidden
    expect(screen.queryByText('Dashboard')).not.toBeInTheDocument()
  })

  it('expands collapsed group on click', async () => {
    const user = userEvent.setup()
    renderWithProviders(<Sidebar />)
    // Finance might be collapsed initially (not active route), click to expand
    await user.click(screen.getByText('Finance'))
    expect(screen.getByText('Invoices')).toBeInTheDocument()
    expect(screen.getByText('Salary')).toBeInTheDocument()
    expect(screen.getByText('Advances')).toBeInTheDocument()
    expect(screen.getByText('Expenses')).toBeInTheDocument()
  })

  it('nav links have correct href', () => {
    renderWithProviders(<Sidebar />)
    const dashLink = screen.getByText('Dashboard').closest('a')
    expect(dashLink).toHaveAttribute('href', '/')
  })

  it('highlights active nav item', () => {
    window.history.pushState({}, '', '/')
    renderWithProviders(<Sidebar />)
    const dashLink = screen.getByText('Dashboard').closest('a')
    expect(dashLink).toHaveClass('bg-sidebar-accent')
  })
})
