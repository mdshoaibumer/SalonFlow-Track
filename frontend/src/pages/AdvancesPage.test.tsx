import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { AdvancesPage } from './AdvancesPage'

describe('AdvancesPage', () => {
  it('renders page header', () => {
    renderWithProviders(<AdvancesPage />)
    expect(screen.getByText('Advance Management')).toBeInTheDocument()
  })

  it('renders advance data after loading', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('Rahul Kumar')).toBeInTheDocument()
    })
  })

  it('renders New Advance button', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /new advance/i })).toBeInTheDocument()
    })
  })

  it('opens add advance dialog', async () => {
    const user = userEvent.setup()
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /new advance/i })).toBeInTheDocument()
    })
    await user.click(screen.getByRole('button', { name: /new advance/i }))
    await waitFor(() => {
      expect(screen.getByText('New Advance Request')).toBeInTheDocument()
    })
  })

  it('renders advance amount', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('₹5,000')).toBeInTheDocument()
    })
  })

  it('renders advance status', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('pending')).toBeInTheDocument()
    })
  })

  it('renders approve/reject action buttons for pending', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('Rahul Kumar')).toBeInTheDocument()
    })
    // Action buttons are icon buttons in a flex container
    const actionButtons = screen.getAllByRole('button').filter(
      b => b.classList.contains('h-7')
    )
    expect(actionButtons.length).toBe(2) // Approve + Reject
  })

  it('renders search input', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/search/i)).toBeInTheDocument()
    })
  })

  it('renders reason for advance', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('Personal')).toBeInTheDocument()
    })
  })

  it('renders status filter', async () => {
    renderWithProviders(<AdvancesPage />)
    await waitFor(() => {
      expect(screen.getByText('All Status')).toBeInTheDocument()
    })
  })
})
