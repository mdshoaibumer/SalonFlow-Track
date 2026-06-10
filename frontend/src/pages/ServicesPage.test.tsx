import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import { renderWithProviders } from '@/test-utils'
import { ServicesPage } from './ServicesPage'

describe('ServicesPage', () => {
  it('renders page header', () => {
    renderWithProviders(<ServicesPage />)
    expect(screen.getByText('Services')).toBeInTheDocument()
  })

  it('renders service data after loading', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByText('Haircut - Ladies')).toBeInTheDocument()
    })
  })

  it('renders KPI cards with stats', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByText('Total Services')).toBeInTheDocument()
    })
  })

  it('renders Add Service button', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /add service/i })).toBeInTheDocument()
    })
  })

  it('renders search input', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/search/i)).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export csv/i })).toBeInTheDocument()
    })
  })

  it('renders service price', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getAllByText('₹500').length).toBeGreaterThanOrEqual(1)
    })
  })

  it('renders category badge', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByText('Hair')).toBeInTheDocument()
    })
  })

  it('renders category filter', async () => {
    renderWithProviders(<ServicesPage />)
    await waitFor(() => {
      expect(screen.getByText('All Categories')).toBeInTheDocument()
    })
  })
})
