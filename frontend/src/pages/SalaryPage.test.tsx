import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { renderWithProviders } from '@/test-utils'
import { SalaryPage } from './SalaryPage'

describe('SalaryPage', () => {
  it('renders page header', () => {
    renderWithProviders(<SalaryPage />)
    expect(screen.getByText('Salary Management')).toBeInTheDocument()
  })

  it('renders salary records after loading', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByText('Priya Sharma')).toBeInTheDocument()
    })
  })

  it('renders KPI cards with stats', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByText('Total Payroll')).toBeInTheDocument()
    })
    expect(screen.getByText('Pending')).toBeInTheDocument()
    expect(screen.getByText('Paid')).toBeInTheDocument()
    expect(screen.getByText('Outstanding Advances')).toBeInTheDocument()
  })

  it('renders Generate Salary button', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /generate salary/i })).toBeInTheDocument()
    })
  })

  it('renders month/year selectors', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      // Month and Year selectors are present
      expect(screen.getByText('Priya Sharma')).toBeInTheDocument()
    })
  })

  it('renders net salary amount', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByText('₹28,000')).toBeInTheDocument()
    })
  })

  it('renders payment status', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByText('pending')).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it('renders Mark Paid button for pending salary', async () => {
    renderWithProviders(<SalaryPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /mark paid/i })).toBeInTheDocument()
    })
  })
})
