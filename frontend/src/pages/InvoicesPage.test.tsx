import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import { renderWithProviders } from '@/test-utils'
import { InvoicesPage } from './InvoicesPage'

describe('InvoicesPage', () => {
  it('renders page header', () => {
    renderWithProviders(<InvoicesPage />)
    expect(screen.getByText('Invoices')).toBeInTheDocument()
  })

  it('renders invoice data after loading', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByText('INV-001')).toBeInTheDocument()
    })
  })

  it('renders KPI cards with stats', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByText("Today's Revenue")).toBeInTheDocument()
    })
    expect(screen.getByText("Today's Invoices")).toBeInTheDocument()
    expect(screen.getByText('Avg Bill Value')).toBeInTheDocument()
  })

  it('renders New Invoice button', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /new invoice/i })).toBeInTheDocument()
    })
  })

  it('renders search input', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/search/i)).toBeInTheDocument()
    })
  })

  it('renders revenue amount', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByText('₹12,500')).toBeInTheDocument()
    })
  })

  it('renders invoice count', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByText('8')).toBeInTheDocument()
    })
  })

  it('renders payment status badge', async () => {
    renderWithProviders(<InvoicesPage />)
    await waitFor(() => {
      expect(screen.getByText('paid')).toBeInTheDocument()
    })
  })
})
