import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { CustomersPage } from './CustomersPage'

describe('CustomersPage', () => {
  it('renders page header', () => {
    renderWithProviders(<CustomersPage />)
    expect(screen.getByText('Customers')).toBeInTheDocument()
  })

  it('renders customer data after loading', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByText('Anjali Desai')).toBeInTheDocument()
    })
  })

  it('renders customer code', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByText('CUS001')).toBeInTheDocument()
    })
  })

  it('renders KPI cards', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByText('Total Customers')).toBeInTheDocument()
    })
    expect(screen.getByText('Active')).toBeInTheDocument()
    expect(screen.getByText('Total Revenue')).toBeInTheDocument()
    expect(screen.getByText('Avg Visits')).toBeInTheDocument()
  })

  it('renders Add Customer button', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add customer/i })).toBeInTheDocument()
    })
  })

  it('renders search input', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/search/i)).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it('renders customer phone number', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getByText('9876543211')).toBeInTheDocument()
    })
  })

  it('renders visit count', async () => {
    renderWithProviders(<CustomersPage />)
    await waitFor(() => {
      expect(screen.getAllByText('5').length).toBeGreaterThanOrEqual(1)
    })
  })
})
