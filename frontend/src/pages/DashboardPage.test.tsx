import { describe, it, expect, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import { renderWithProviders } from '@/test-utils'
import { DashboardPage } from './DashboardPage'

// Mock recharts to avoid canvas rendering issues in jsdom
vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div data-testid="chart-container">{children}</div>,
  AreaChart: ({ children }: { children: React.ReactNode }) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div />,
  BarChart: ({ children }: { children: React.ReactNode }) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div />,
  XAxis: () => <div />,
  YAxis: () => <div />,
  CartesianGrid: () => <div />,
  Tooltip: () => <div />,
}))

describe('DashboardPage', () => {
  it('renders loading state initially', () => {
    renderWithProviders(<DashboardPage />)
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })

  it('renders KPI cards after loading', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByText('Revenue Today')).toBeInTheDocument()
    })
    expect(screen.getByText('Customers Today')).toBeInTheDocument()
    expect(screen.getByText('Average Bill')).toBeInTheDocument()
    expect(screen.getByText('Top Performer')).toBeInTheDocument()
  })

  it('renders revenue value from API', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getAllByText('₹12,500').length).toBeGreaterThanOrEqual(1)
    })
  })

  it('renders customer count from API', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getAllByText('8').length).toBeGreaterThanOrEqual(1)
    })
  })

  it('renders top performer name', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByText('Priya Sharma')).toBeInTheDocument()
    })
  })

  it('renders quick actions section', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByText('Quick Actions')).toBeInTheDocument()
    })
    expect(screen.getByText('New Bill')).toBeInTheDocument()
    expect(screen.getByText('Add Customer')).toBeInTheDocument()
    expect(screen.getByText('Appointment')).toBeInTheDocument()
    expect(screen.getByText('Services')).toBeInTheDocument()
  })

  it('renders charts section', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByText('Revenue Trend (14 Days)')).toBeInTheDocument()
    })
    expect(screen.getByText('Top Performers')).toBeInTheDocument()
  })

  it('renders staff summary card', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByText('Staff Summary')).toBeInTheDocument()
    })
  })

  it('renders export button', async () => {
    renderWithProviders(<DashboardPage />)
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /export/i })).toBeInTheDocument()
    })
  })
})
