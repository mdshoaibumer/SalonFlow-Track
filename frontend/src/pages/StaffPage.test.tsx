import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { StaffPage } from './StaffPage'

describe('StaffPage', () => {
  it('renders loading state initially', () => {
    renderWithProviders(<StaffPage />)
    expect(screen.getByText('Staff')).toBeInTheDocument()
  })

  it('renders staff data after loading', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('Priya Sharma')).toBeInTheDocument()
    })
  })

  it('renders staff code in table', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('STF001')).toBeInTheDocument()
    })
  })

  it('renders KPI cards with stats', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('Total Staff')).toBeInTheDocument()
    })
    expect(screen.getByText('Active')).toBeInTheDocument()
    expect(screen.getByText('Inactive')).toBeInTheDocument()
  })

  it('renders Add Staff button', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add staff/i })).toBeInTheDocument()
    })
  })

  it('opens form dialog when Add Staff clicked', async () => {
    const user = userEvent.setup()
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add staff/i })).toBeInTheDocument()
    })
    await user.click(screen.getByRole('button', { name: /add staff/i }))
    await waitFor(() => {
      // Dialog opens with form fields
      expect(screen.getByLabelText(/full name/i)).toBeInTheDocument()
    })
  })

  it('renders status filter', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('All Status')).toBeInTheDocument()
    })
  })

  it('renders search input', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText('Search by name, phone, or code...')).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it('renders designation in table', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('Stylist')).toBeInTheDocument()
    })
  })

  it('renders staff status badge', async () => {
    renderWithProviders(<StaffPage />)
    await waitFor(() => {
      expect(screen.getByText('active')).toBeInTheDocument()
    })
  })
})
