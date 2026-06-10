import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { ExpensesPage } from './ExpensesPage'

describe('ExpensesPage', () => {
  it('renders page header', () => {
    renderWithProviders(<ExpensesPage />)
    expect(screen.getByText('Expenses')).toBeInTheDocument()
  })

  it('renders expense data after loading', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByText('EXP-001')).toBeInTheDocument()
    })
  })

  it('renders KPI cards with stats', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByText("Today's Expenses")).toBeInTheDocument()
    })
    expect(screen.getByText('Monthly Expenses')).toBeInTheDocument()
    expect(screen.getByText('Monthly Profit')).toBeInTheDocument()
    expect(screen.getByText('Profit Margin')).toBeInTheDocument()
  })

  it('renders Add Expense button', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add expense/i })).toBeInTheDocument()
    })
  })

  it('opens add expense dialog', async () => {
    const user = userEvent.setup()
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add expense/i })).toBeInTheDocument()
    })
    await user.click(screen.getByRole('button', { name: /add expense/i }))
    await waitFor(() => {
      // Dialog opens - check for label that only appears in dialog
      expect(screen.getByText('Amount (₹) *')).toBeInTheDocument()
    })
  })

  it('renders search input', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/search/i)).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it('renders expense amount', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByText(/25,000/)).toBeInTheDocument()
    })
  })

  it('renders vendor name', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByText('Landlord')).toBeInTheDocument()
    })
  })

  it('renders category name', async () => {
    renderWithProviders(<ExpensesPage />)
    await waitFor(() => {
      expect(screen.getByText('Rent')).toBeInTheDocument()
    })
  })
})
